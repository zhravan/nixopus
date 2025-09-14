import os
import platform
import shutil
import subprocess
import sys
import time
import traceback
from dataclasses import dataclass
from pathlib import Path
from typing import Optional

from app.utils.config import Config
from app.utils.logger import Logger
from app.utils.lib import HostInformation, DirectoryManager

from app.commands.preflight.run import PreflightRunner
from app.commands.conflict.conflict import ConflictConfig, ConflictService
from app.commands.install.deps import install_all_deps
from app.commands.install.ssh import SSH, SSHConfig
from app.commands.clone.clone import Clone, CloneConfig

# Import installer infrastructure

# Add the project root to Python path for installer imports
current_file = Path(__file__).resolve()
project_root = current_file.parent.parent.parent.parent  # cli/app/commands/setup -> project root
installer_path = project_root / "installer"

if installer_path.exists() and str(project_root) not in sys.path:
    sys.path.insert(0, str(project_root))

try:
    from installer.service_manager import ServiceManager

    INSTALLER_AVAILABLE = True
except ImportError:
    # Fallback when installer module is not available
    INSTALLER_AVAILABLE = False

    class ServiceManager:
        def __init__(self, *args, **kwargs):
            pass

        def check_system_requirements(self):
            pass

        def start_services(self, env):
            pass

        def verify_installation(self, env):
            pass

        def check_api_up_status(self, port):
            return True

        def setup_admin(self, email, password, port):
            pass


@dataclass
class DevSetupConfig:
    """Configuration for development environment setup"""

    api_port: int = 8080
    view_port: int = 7443
    db_port: int = 5432
    redis_port: int = 6379
    branch: str = "feat/develop"
    repo: Optional[str] = None
    workspace: str = "./nixopus-dev"
    ssh_key_path: str = "~/.ssh/id_ed25519_nixopus"
    ssh_key_type: str = "ed25519"
    skip_preflight: bool = False
    skip_conflict: bool = False
    skip_deps: bool = False
    skip_docker: bool = False
    skip_ssh: bool = False
    skip_admin: bool = False
    admin_email: Optional[str] = None
    admin_password: str = "Nixopus123!"
    force: bool = False
    dry_run: bool = False
    verbose: bool = False
    timeout: int = 300
    output: str = "text"
    config_file: str = "../helpers/config.dev.yaml"


class DevSetup:
    """Development environment setup orchestrator"""

    def __init__(self, config: DevSetupConfig, logger: Logger):
        self.config = config
        self.logger = logger
        self.nixopus_config = Config()

        # Set config file to dev configuration
        self.nixopus_config._yaml_path = os.path.abspath(self.config.config_file)

        # Workspace setup
        self.workspace_path = Path(self.config.workspace).resolve()
        self.project_root = self.workspace_path

        # Set default admin email if not provided
        if not self.config.admin_email:
            self.config.admin_email = f"{os.environ.get('USER', 'admin')}@example.com"

    def run(self):
        """Execute complete development setup"""
        self.logger.info("Starting Nixopus development environment setup...")

        try:
            # Phase 1: System Validation
            if not self.config.skip_preflight or not self.config.skip_conflict:
                self._validate_system()

            # Phase 2: Dependencies
            if not self.config.skip_deps:
                self._install_dependencies()

            # Phase 3: Repository Setup
            self._setup_repository()

            # Phase 4: Project Dependencies
            self._setup_project_dependencies()

            # Phase 5: Infrastructure (Database)
            if not self.config.skip_docker:
                self._setup_infrastructure()

            # Phase 6: SSH Configuration
            if not self.config.skip_ssh:
                self._setup_ssh()

            # Phase 7: Environment Configuration
            self._configure_environment()

            # Phase 8: Service Startup & Validation
            self._start_services()

            # Phase 9: Admin Setup
            if not self.config.skip_admin:
                self._create_admin()

            # Phase 10: Completion
            self._display_completion_info()

        except Exception as e:
            self.logger.error(f"Setup failed at phase: {e}")
            if self.config.verbose:
                self.logger.error(traceback.format_exc())
            raise

    def _validate_system(self):
        """Phase 1: System validation using existing preflight and conflict checks"""
        self.logger.info("Phase 1: Validating system requirements...")

        if not self.config.skip_preflight:
            self.logger.info("Running preflight checks...")
            try:
                preflight_runner = PreflightRunner(logger=self.logger, verbose=self.config.verbose)
                # Run OS-specific preflight checks
                preflight_runner.check_windows_environment()
                preflight_runner.check_ports_from_config()
                self.logger.success("Preflight checks passed")
            except Exception as e:
                self.logger.error(f"Preflight checks failed: {e}")
                raise

        if not self.config.skip_conflict:
            self.logger.info("Checking for tool conflicts...")
            try:
                conflict_config = ConflictConfig(
                    config_file=self.config.config_file,
                    verbose=self.config.verbose,
                    output=self.config.output,
                )
                conflict_service = ConflictService(conflict_config, logger=self.logger)
                results = conflict_service.check_conflicts()
                conflicts = [r for r in results if r.conflict]

                missing_tools = []
                version_mismatches = []
                if conflicts:
                    for conflict in conflicts:
                        if conflict.current is None:
                            missing_tools.append(conflict)
                        else:
                            version_mismatches.append(conflict)

                if version_mismatches:
                    self.logger.warning(
                        f"Tool version mismatches found: {len(version_mismatches)}. This is not a blocker for dev setup."
                    )
                    for conflict in version_mismatches:
                        self.logger.warning(f"  - {conflict.tool}: Expected {conflict.expected}, Found {conflict.current}")

                if missing_tools:
                    self.logger.error(f"Missing required tools: {len(missing_tools)}")
                    for conflict in missing_tools:
                        self.logger.error(f"  - {conflict.tool}: {conflict.status}")
                    raise Exception("Missing required tools. Please install them and try again.")

                if not conflicts:
                    self.logger.success("No tool conflicts detected.")
                elif not missing_tools:
                    self.logger.success("Tool conflict check passed (version mismatches are non-blocking).")
            except Exception as e:
                self.logger.error(f"Conflict detection failed: {e}")
                raise

    def _install_dependencies(self):
        """Phase 2: Install system dependencies using existing installer"""
        self.logger.info("Phase 2: Installing system dependencies...")

        if self.config.dry_run:
            self.logger.info("[DRY RUN] Would install dependencies from config.dev.yaml")
            return

        try:
            install_all_deps(verbose=self.config.verbose, output=self.config.output, dry_run=self.config.dry_run)
            self.logger.success("System dependencies installed")
        except Exception as e:
            self.logger.error(f"Dependency installation failed: {e}")
            raise

    def _setup_repository(self):
        """Phase 3: Clone repository using existing clone functionality"""
        self.logger.info("Phase 3: Setting up repository...")

        # Get repository configuration from config.dev.yaml
        repo_url = self.config.repo or self.nixopus_config.get_yaml_value("clone.repo")
        branch = self.config.branch

        clone_path = str(self.workspace_path)

        if self.config.dry_run:
            self.logger.info(f"[DRY RUN] Would clone {repo_url}#{branch} to {clone_path}")
            return

        try:
            # Check if directory exists
            if self.workspace_path.exists() and not self.config.force:
                if any(self.workspace_path.iterdir()):
                    raise Exception(f"Directory {clone_path} already exists and is not empty. Use --force to overwrite.")

            clone_config = CloneConfig(
                repo=repo_url,
                branch=branch,
                path=clone_path,
                force=self.config.force,
                verbose=self.config.verbose,
                output=self.config.output,
                dry_run=self.config.dry_run,
            )

            clone_operation = Clone(logger=self.logger)
            result = clone_operation.clone(clone_config)

            if not result.success:
                raise Exception(f"Repository clone failed: {result.output}")

            self.logger.success(f"Repository cloned to {clone_path}")

        except Exception as e:
            self.logger.error(f"Repository setup failed: {e}")
            raise

    def _setup_project_dependencies(self):
        """Phase 4: Setup project-specific dependencies (Go, Yarn, Poetry)"""
        self.logger.info("Phase 4: Setting up project dependencies...")

        if self.config.dry_run:
            self.logger.info("[DRY RUN] Would setup Go modules, Yarn packages, and Poetry dependencies")
            return

        # Change to project directory
        original_cwd = os.getcwd()
        try:
            os.chdir(self.workspace_path)

            # Setup Go dependencies for API
            api_path = self.workspace_path / "api"
            if api_path.exists():
                self.logger.info("Setting up Go dependencies...")
                os.chdir(api_path)
                subprocess.run(["go", "mod", "tidy"], check=True)
                subprocess.run(["go", "mod", "download"], check=True)
                self.logger.success("Go dependencies installed")
                os.chdir(self.workspace_path)

            # Setup Yarn dependencies for view
            view_path = self.workspace_path / "view"
            if view_path.exists():
                self.logger.info("Setting up Yarn dependencies...")
                os.chdir(view_path)
                subprocess.run(["yarn", "install", "--frozen-lockfile"], check=True)
                self.logger.success("Yarn dependencies installed")
                os.chdir(self.workspace_path)

            # Setup Poetry dependencies for CLI
            cli_path = self.workspace_path / "cli"
            if cli_path.exists() and (cli_path / "pyproject.toml").exists():
                self.logger.info("Setting up Poetry dependencies...")
                os.chdir(cli_path)
                subprocess.run(["poetry", "install"], check=True)
                self.logger.success("Poetry dependencies installed")
                os.chdir(self.workspace_path)

        except subprocess.CalledProcessError as e:
            self.logger.error(f"❌ Project dependency setup failed: {e}")
            raise
        except Exception as e:
            self.logger.error(f"❌ Project dependency setup failed: {e}")
            raise
        finally:
            os.chdir(original_cwd)

    def _setup_infrastructure(self):
        """Phase 5: Setup PostgreSQL and Redis using existing service manager"""
        self.logger.info("Phase 5: Setting up database infrastructure...")

        if self.config.dry_run:
            self.logger.info("[DRY RUN] Would setup PostgreSQL and Redis containers")
            return

        if not INSTALLER_AVAILABLE:
            self.logger.warning("Installer module not available - using Docker directly for database + Redis setup")
            self._setup_database_docker_direct()
            return

        try:
            # Use existing service manager from installer
            service_manager = ServiceManager(
                project_root=self.workspace_path, env="development", debug=self.config.verbose  # Use development mode
            )

            service_manager.check_system_requirements()

            # Start database services
            service_manager.start_services("development")
            service_manager.verify_installation("development")

            self.logger.success("Database infrastructure ready")

        except Exception as e:
            self.logger.error(f"Infrastructure setup failed: {e}")
            raise

    def _setup_database_docker_direct(self):
        """Fallback method to setup Postgres and Redis using direct Docker commands"""
        self.logger.info("Setting up PostgreSQL and Redis with Docker directly...")

        try:
            # Ensure Docker CLI is available
            if not shutil.which("docker"):
                raise Exception(
                    "Docker is not installed or not in PATH. Please install Docker Desktop (macOS) or Docker Engine (Linux) and ensure it is running."
                )

            try:
                db_env = self.nixopus_config.get_service_env_values("services.db.env")
            except Exception:
                db_env = {}

            try:
                redis_env = self.nixopus_config.get_service_env_values("services.redis.env")
            except Exception:
                redis_env = {}

            db_container_name = db_env.get("DB_CONTAINER_NAME", "nixopus-db")
            db_image = db_env.get("DB_IMAGE", "postgres:14-alpine")
            db_port = str(db_env.get("DB_PORT", self.config.db_port))
            pg_user = db_env.get("POSTGRES_USER", "postgres")
            pg_password = db_env.get("POSTGRES_PASSWORD", "postgres")
            pg_db = db_env.get("POSTGRES_DB", "postgres")
            pg_auth = db_env.get("POSTGRES_HOST_AUTH_METHOD", "trust")

            redis_container_name = redis_env.get("REDIS_CONTAINER_NAME", "nixopus-redis")
            redis_image = redis_env.get("REDIS_IMAGE", "redis:7-alpine")
            redis_port = str(redis_env.get("REDIS_PORT", self.config.redis_port))

            # PostgreSQL
            check_db_cmd = [
                "docker",
                "ps",
                "-a",
                "--format",
                "{{.Names}}",
                "--filter",
                f"name={db_container_name}",
            ]
            db_result = subprocess.run(check_db_cmd, capture_output=True, text=True, check=True)
            if db_container_name in db_result.stdout:
                self.logger.info("Database container already exists")
            else:
                db_run_cmd = [
                    "docker",
                    "run",
                    "-d",
                    "--name",
                    db_container_name,
                    "-e",
                    f"POSTGRES_USER={pg_user}",
                    "-e",
                    f"POSTGRES_PASSWORD={pg_password}",
                    "-e",
                    f"POSTGRES_DB={pg_db}",
                    "-e",
                    f"POSTGRES_HOST_AUTH_METHOD={pg_auth}",
                    "-p",
                    f"{db_port}:5432",
                    "--health-cmd",
                    f"pg_isready -U {pg_user} -d {pg_db}",
                    db_image,
                ]
                subprocess.run(db_run_cmd, check=True)
                self.logger.success("Database container started")

            # Redis
            check_redis_cmd = [
                "docker",
                "ps",
                "-a",
                "--format",
                "{{.Names}}",
                "--filter",
                f"name={redis_container_name}",
            ]
            redis_result = subprocess.run(check_redis_cmd, capture_output=True, text=True, check=True)
            if redis_container_name in redis_result.stdout:
                self.logger.info("Redis container already exists")
            else:
                redis_run_cmd = [
                    "docker",
                    "run",
                    "-d",
                    "--name",
                    redis_container_name,
                    "-p",
                    f"{redis_port}:6379",
                    "--health-cmd",
                    "redis-cli ping || exit 1",
                    redis_image,
                ]
                subprocess.run(redis_run_cmd, check=True)
                self.logger.success("Redis container started")

        except subprocess.CalledProcessError as e:
            self.logger.error(f"Container setup failed: {e}")
            raise
        except Exception as e:
            self.logger.error(f"Infrastructure setup error: {e}")
            raise

    def _setup_ssh(self):
        """Phase 6: Setup SSH using existing SSH installer"""
        self.logger.info("Phase 6: Setting up SSH configuration...")

        if self.config.dry_run:
            self.logger.info(f"[DRY RUN] Would generate SSH key at {self.config.ssh_key_path}")
            return

        try:
            # Ensure SSH tooling is present (client and, where applicable, server)
            self._ensure_ssh_tools()
            ssh_config = SSHConfig(
                path=self.config.ssh_key_path,
                key_type=self.config.ssh_key_type,
                key_size=4096,
                passphrase="",
                verbose=self.config.verbose,
                output=self.config.output,
                dry_run=self.config.dry_run,
                force=self.config.force,
                set_permissions=True,
                add_to_authorized_keys=True,
                create_ssh_directory=True,
            )

            ssh_operation = SSH(logger=self.logger)
            result = ssh_operation.generate(ssh_config)

            self.logger.success("SSH configuration completed")

        except Exception as e:
            self.logger.error(f"SSH setup failed: {e}")
            raise

    def _ensure_ssh_tools(self):
        """Install or guide installing SSH client/server based on OS/distro."""
        os_name = platform.system().lower()

        def have(cmd: str) -> bool:
            return shutil.which(cmd) is not None

        if os_name == "darwin":
            # macOS ships with ssh; server is via Remote Login. Install brew openssh if missing.
            if not have("ssh"):
                if have("brew"):
                    self.logger.info("Installing OpenSSH via Homebrew (ssh client)...")
                    subprocess.run(["brew", "install", "openssh"], check=True)
                else:
                    raise Exception(
                        "ssh client not found. Install Command Line Tools (xcode-select --install) or Homebrew and 'brew install openssh'."
                    )
            # Provide guidance for enabling the server on macOS
            self.logger.info("For SSH server on macOS, enable 'Remote Login' in System Settings > Sharing.")
            return

        if os_name == "windows":
            # Windows 10/11: Prefer built-in OpenSSH via Windows Capabilities, fallback to winget/choco.
            def run_powershell(ps_cmd: str, check: bool = True):
                return subprocess.run(
                    ["powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", ps_cmd],
                    check=check,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                    text=True,
                )

            if not have("ssh"):
                installed_client = False
                # Try Windows capability for OpenSSH Client
                try:
                    self.logger.info("Checking OpenSSH Client capability on Windows...")
                    state_res = run_powershell(
                        "(Get-WindowsCapability -Online | Where-Object { $_.Name -like 'OpenSSH.Client*' }).State",
                        check=False,
                    )
                    if "Installed" not in state_res.stdout:
                        self.logger.info("Installing OpenSSH Client capability (requires admin)...")
                        run_powershell("Add-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0", check=True)
                    installed_client = True
                except Exception as e:
                    self.logger.warning(f"Windows capability install for OpenSSH Client failed: {e}")

                # Fallback: winget
                if not installed_client and have("winget"):
                    try:
                        self.logger.info("Installing OpenSSH Client via winget (requires admin consent)...")
                        subprocess.run(["winget", "install", "--id", "Microsoft.OpenSSH.Beta", "--source", "winget"], check=True)
                        installed_client = True
                    except Exception as e:
                        self.logger.warning(f"winget install of OpenSSH failed: {e}")

                # Fallback: choco
                if not installed_client and have("choco"):
                    try:
                        self.logger.info("Installing OpenSSH via Chocolatey (requires admin)...")
                        subprocess.run(["choco", "install", "openssh", "-y"], check=True)
                        installed_client = True
                    except Exception as e:
                        self.logger.warning(f"Chocolatey install of OpenSSH failed: {e}")

                if not have("ssh"):
                    self.logger.error(
                        "ssh client not found on Windows. Please install 'OpenSSH Client' via Optional Features or run PowerShell as Administrator: Add-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0"
                    )
                    return

            # Try to install and start the OpenSSH Server optionally (best-effort)
            try:
                state_res = run_powershell(
                    "(Get-WindowsCapability -Online | Where-Object { $_.Name -like 'OpenSSH.Server*' }).State",
                    check=False,
                )
                if "Installed" not in state_res.stdout:
                    self.logger.info("Installing OpenSSH Server capability (optional, requires admin)...")
                    run_powershell("Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0", check=False)

                # Enable and start sshd service (ignore failures if not elevated)
                self.logger.info("Ensuring sshd service is enabled and started (optional)...")
                run_powershell("Set-Service -Name sshd -StartupType 'Automatic'", check=False)
                run_powershell("Start-Service sshd", check=False)
            except Exception as e:
                self.logger.warning(f"OpenSSH Server setup skipped/failed: {e}")

            self.logger.info("Windows SSH tooling ready.")
            return

        # Linux: ensure ssh client and server packages
        try:
            package_manager = HostInformation.get_package_manager()
        except Exception as e:
            self.logger.warning(f"Could not determine package manager for SSH install: {e}")
            return

        install_cmds = []
        if package_manager == "apt":
            install_cmds = [
                ["sudo", "apt-get", "update"],
                ["sudo", "apt-get", "install", "-y", "openssh-client", "openssh-server"],
            ]
        elif package_manager in ("yum", "dnf"):
            pkg = "dnf" if package_manager == "dnf" else "yum"
            install_cmds = [["sudo", pkg, "install", "-y", "openssh-clients", "openssh-server"]]
        elif package_manager == "apk":
            install_cmds = [["sudo", "apk", "add", "openssh-client", "openssh-server"]]
        elif package_manager == "pacman":
            install_cmds = [["sudo", "pacman", "-S", "--noconfirm", "openssh"]]

        for cmd in install_cmds:
            try:
                subprocess.run(cmd, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
            except subprocess.CalledProcessError as e:
                self.logger.warning(f"SSH install step failed: {' '.join(cmd)} ({e})")

        # Try to enable and start sshd if available (ignore failures)
        if have("systemctl"):
            for svc in ["ssh", "sshd", "ssh.service", "sshd.service"]:
                try:
                    subprocess.run(
                        ["sudo", "systemctl", "enable", "--now", svc],
                        check=False,
                        stdout=subprocess.DEVNULL,
                        stderr=subprocess.DEVNULL,
                    )
                except Exception:
                    pass

    def _configure_environment(self):
        """Phase 7: Configure environment variables and service configuration"""
        self.logger.info("Phase 7: Configuring environment...")

        if self.config.dry_run:
            self.logger.info("[DRY RUN] Would configure .env files and service settings")
            return

        try:
            # Update ports in environment files using config.dev.yaml values
            self._update_env_files()

            self.logger.success("Environment configured")

        except Exception as e:
            self.logger.error(f"Environment configuration failed: {e}")
            raise

    def _update_env_files(self):
        """Update .env files with custom port configurations"""
        # Copy .env.sample to .env for API
        api_env_sample = self.workspace_path / "api" / ".env.sample"
        api_env = self.workspace_path / "api" / ".env"

        if api_env_sample.exists():
            shutil.copy2(api_env_sample, api_env)

            # Update API environment variables
            self._update_env_file(
                api_env,
                {
                    "PORT": str(self.config.api_port),
                    "DB_PORT": str(self.config.db_port),
                    "REDIS_URL": f"redis://localhost:{self.config.redis_port}",
                    "ALLOWED_ORIGIN": f"http://localhost:{self.config.view_port}",
                    "SSH_PRIVATE_KEY": os.path.expanduser(self.config.ssh_key_path),
                    "SSH_HOST": "localhost",
                    "SSH_PORT": "22",
                    "SSH_USER": os.environ.get("USER", "admin"),
                    "ENV": "development",
                    "DB_NAME": "postgres",
                    "USERNAME": "postgres",
                    "PASSWORD": "changeme",
                    "HOST_NAME": "localhost",
                    "SSL_MODE": "disable",
                },
            )

        # Copy .env.sample to .env for view
        view_env_sample = self.workspace_path / "view" / ".env.sample"
        view_env = self.workspace_path / "view" / ".env"

        if view_env_sample.exists():
            shutil.copy2(view_env_sample, view_env)

            # Update view environment variables
            self._update_env_file(
                view_env,
                {
                    "PORT": str(self.config.view_port),
                    "NEXT_PUBLIC_PORT": str(self.config.view_port),
                    "API_URL": f"http://localhost:{self.config.api_port}",
                    "WEBSOCKET_URL": f"ws://localhost:{self.config.api_port}",
                    "WEBHOOK_URL": f"http://localhost:{self.config.api_port}",
                },
            )

    def _update_env_file(self, env_file: Path, updates: dict):
        """Update specific environment file with new values"""
        if not env_file.exists():
            return

        # Read current content
        lines = env_file.read_text().splitlines()

        # Update lines
        updated_lines = []
        updated_keys = set()

        for line in lines:
            if "=" in line and not line.strip().startswith("#"):
                key = line.split("=")[0].strip()
                if key in updates:
                    updated_lines.append(f"{key}={updates[key]}")
                    updated_keys.add(key)
                else:
                    updated_lines.append(line)
            else:
                updated_lines.append(line)

        # Add any missing keys
        for key, value in updates.items():
            if key not in updated_keys:
                updated_lines.append(f"{key}={value}")

        # Write back
        env_file.write_text("\n".join(updated_lines) + "\n")

    def _start_services(self):
        """Phase 8: Start API and frontend services"""
        self.logger.info("Phase 8: Starting services...")

        if self.config.dry_run:
            self.logger.info("[DRY RUN] Would start API and frontend services")
            return

        try:
            # Start API service with Air hot reload
            self._start_api_service()

            # Start frontend service
            self._start_frontend_service()

            self.logger.success("Services started successfully")

        except Exception as e:
            self.logger.error(f"Service startup failed: {e}")
            raise

    def _start_api_service(self):
        """Start API service with Air hot reload"""
        api_path = self.workspace_path / "api"
        if not api_path.exists():
            raise Exception("API directory not found")

        self.logger.info("Starting API service with Air hot reload...")

        # Check if Air is installed
        if not shutil.which("air"):
            # Install Air if not available
            air_install_cmd = ["go", "install", "github.com/air-verse/air@latest"]
            subprocess.run(air_install_cmd, check=True, cwd=api_path)

        # Start Air in background
        log_file = api_path / "api.log"
        with open(log_file, "w") as f:
            subprocess.Popen(["air"], cwd=api_path, stdout=f, stderr=subprocess.STDOUT, start_new_session=True)

        self.logger.info(f"API service started (logs: {log_file})")

    def _start_frontend_service(self):
        """Start frontend service with Yarn"""
        view_path = self.workspace_path / "view"
        if not view_path.exists():
            raise Exception("View directory not found")

        self.logger.info("Starting frontend service...")

        # Start Yarn dev server in background
        log_file = view_path / "view.log"
        with open(log_file, "w") as f:
            subprocess.Popen(
                ["yarn", "run", "dev", "--", "-p", str(self.config.view_port)],
                cwd=view_path,
                stdout=f,
                stderr=subprocess.STDOUT,
                start_new_session=True,
            )

        self.logger.info(f"Frontend service started (logs: {log_file})")

    def _create_admin(self):
        """Phase 9: Create admin account using existing service manager"""
        self.logger.info("Phase 9: Creating admin account...")

        if self.config.dry_run:
            self.logger.info(f"[DRY RUN] Would create admin: {self.config.admin_email}")
            return

        if not INSTALLER_AVAILABLE:
            self.logger.warning("Installer module not available - skipping admin account creation")
            self.logger.info("You can create admin manually via API after setup")
            return

        try:
            service_manager = ServiceManager(project_root=self.workspace_path, env="development", debug=self.config.verbose)

            # Wait for API to be ready
            max_retries = 10
            for attempt in range(max_retries):
                if service_manager.check_api_up_status(self.config.api_port):
                    break
                self.logger.info(f"Waiting for API... (attempt {attempt + 1}/{max_retries})")
                time.sleep(3)
            else:
                raise Exception("API service not responding after maximum retries")

            # Create admin account
            service_manager.setup_admin(
                email=self.config.admin_email, password=self.config.admin_password, port=self.config.api_port
            )

            self.logger.success("Admin account created")

        except Exception as e:
            self.logger.error(f"Admin account creation failed: {e}")
            # Don't raise - this is not critical for development setup

    def _display_completion_info(self):
        """Phase 10: Display completion information and guidance"""
        self.logger.info("Phase 10: Setup completed!")

        print("\n" + "=" * 70)
        print("NIXOPUS DEVELOPMENT ENVIRONMENT READY!")
        print("=" * 70)
        print()
        print("Application Access:")
        print(f"   - Frontend: http://localhost:{self.config.view_port}")
        print(f"   - API:      http://localhost:{self.config.api_port}")
        print(f"   - Database: localhost:{self.config.db_port}")
        print()
        print("Default Login Credentials:")
        print(f"   - Email:    {self.config.admin_email}")
        print(f"   - Password: {self.config.admin_password}")
        print("   WARNING: Change these credentials after first login!")
        print()
        print("Project Location:")
        print(f"   - Workspace: {self.workspace_path}")
        print()
        print("Useful Commands:")
        print("   - Stop API:      pkill -f air")
        print("   - Stop Frontend: pkill -f yarn")
        print("   - View API logs: tail -f api/api.log")
        print("   - View Frontend logs: tail -f view/view.log")
        print()
        print("Troubleshooting:")
        print("   - Check containers: docker ps")
        print("   - Database logs:    docker logs nixopus-db")
        print("   - Redis logs:       docker logs nixopus-redis")
        print("   - Restart database: docker restart nixopus-db")
        print("   - Restart Redis:    docker restart nixopus-redis")
        print()
        print("Community & Support:")
        print("   - Discord: https://discord.com/invite/skdcq39Wpv")
        print("   - GitHub:  https://github.com/raghavyuva/nixopus")
        print("   - Issues:  https://github.com/raghavyuva/nixopus/issues")
        print()
        print("=" * 70)
        print("Happy coding!")
        print("=" * 70)
