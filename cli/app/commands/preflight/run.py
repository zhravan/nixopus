import platform
import shutil
import subprocess
from typing import List, Optional, Sequence

from app.utils.config import Config
from app.utils.protocols import LoggerProtocol
from app.utils.logger import Logger

from .messages import ports_unavailable
from .port import PortConfig, PortService, PortCheckResult


class PreflightRunner:
    """Centralized preflight check runner for port availability"""

    def __init__(self, logger: Optional[LoggerProtocol] = None, verbose: bool = False):
        self.logger = logger
        self.verbose = verbose
        self.config = Config()

    def _have(self, cmd: str) -> bool:
        return shutil.which(cmd) is not None

    def check_windows_environment(self) -> None:
        """On Windows hosts, verify Docker Desktop, WSL2 readiness.

        1.  Ensures docker CLI exists and Docker daemon is reachable.
        2. Checks WSL presence and recommends WSL2 if not detected.
        """
        if platform.system().lower() != "windows":
            return

        if self.logger:
            self.logger.info("Running Windows preflight checks...")

        # check Docker CLI 
        if not self._have("docker"):
            raise Exception("Docker CLI not found on Windows. Please install Docker Desktop and ensure 'docker' is in PATH.")

        # pin g Docker daemon
        try:
            # quick daemon ping via 'docker info'
            subprocess.run(["docker", "info"], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
            if self.logger:
                self.logger.success("Docker daemon is running.")
        except subprocess.CalledProcessError:
            raise Exception("Docker daemon is not running. Start Docker Desktop and retry.")

        # WSL presence and version check (recommendation only)
        if self._have("wsl"):
            try:
                result = subprocess.run(
                    ["wsl", "-l", "-v"],
                    check=False,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                    text=True,
                )
                output = result.stdout or ""
                # Look for any distro line containing version '2'
                has_wsl2 = any(" 2" in line or "\t2" in line for line in output.splitlines())
                if has_wsl2:
                    if self.logger:
                        self.logger.success("WSL2 detected.")
                else:
                    if self.logger:
                        self.logger.warning("WSL detected but no WSL2 distro found. Docker Desktop works best with WSL2.")
            except Exception:
                if self.logger:
                    self.logger.warning("Unable to verify WSL version. Proceeding.")
        else:
            if self.logger:
                self.logger.warning("WSL not found. Install WSL2 for the best Docker Desktop compatibility (optional).")

    def run_port_checks(self, ports: List[int], host: str = "localhost") -> List[PortCheckResult]:
        """Run port availability checks and return results"""
        port_config = PortConfig(ports=ports, host=host, verbose=self.verbose)
        # Ensure a concrete logger instance is provided
        effective_logger = self.logger or Logger(verbose=self.verbose)
        port_service = PortService(port_config, logger=effective_logger)
        return port_service.check_ports()

    def check_required_ports(self, ports: List[int], host: str = "localhost") -> None:
        """Check required ports and raise exception if any are unavailable"""
        port_results = self.run_port_checks(ports, host)
        unavailable_ports = [result for result in port_results if not result.get("is_available", True)]

        if unavailable_ports:
            error_msg = f"{ports_unavailable}: {[p['port'] for p in unavailable_ports]}"
            raise Exception(error_msg)

    def check_ports_from_config(
        self, config_key: str = "required_ports", user_config: Optional[dict] = None, defaults: Optional[dict] = None
    ) -> None:
        """Check ports using configuration values"""
        if user_config is not None and defaults is not None:
            ports = self.config.get_config_value(config_key, user_config, defaults)
        else:
            ports = self.config.get_yaml_value("ports")

        if not isinstance(ports, list):
            raise Exception("Configured 'ports' must be a list of integers")
        self.check_required_ports([int(p) for p in ports])
