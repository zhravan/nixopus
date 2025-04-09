import os
import sys
import subprocess
import platform
from pathlib import Path
import re

class Installer:
    def __init__(self):
        self.required_docker_version = "20.10.0"
        self.required_compose_version = "2.0.0"
        self.project_root = Path(__file__).parent.parent
        self.env_file = self.project_root / ".env"
        self.env_sample = self.project_root / ".env.sample"
    
    # this script will only work with root privileges
    def ask_for_sudo(self):
        if os.geteuid() != 0:
            print("Please run the script with sudo privileges")
            sys.exit(1)
    
    def ask_email(self):
        email = input("Please enter your email: ")
        if not email:
            print("Error: Email is required")
            sys.exit(1)
        return email
    
    def ask_password(self):
        password = input("Please enter your password: ")
        if not password:
            print("Error: Password is required")
            sys.exit(1)
        return password
    
    def ask_domain(self):
        domain = input("Please enter the domain in which you want to run the application (e.g. nixopus.com): ")
        if not domain:
            print("Error: Domain is required")
            sys.exit(1)
        return domain
    
    def check_docker_version(self):
        docker_version = subprocess.check_output(["docker", "--version"]).decode()
        if not self._version_check(docker_version, self.required_docker_version):
            print(f"Error: Docker version {self.required_docker_version} or higher is required")
            sys.exit(1)
            
    def check_docker_compose_version(self):
        compose_version = subprocess.check_output(["docker-compose", "--version"]).decode()
        if not self._version_check(compose_version, self.required_compose_version):
            print(f"Error: Docker Compose version {self.required_compose_version} or higher is required")
            sys.exit(1)

    def check_system_requirements(self):
        print("Checking system requirements...")
        
        # we will only support linux for now
        system = platform.system()
        if system not in ["Linux"]:
            print(f"Error: Unsupported operating system: {system}")
            sys.exit(1)

        self.check_docker_version()
        self.check_docker_compose_version()
        
        print("System requirements check passed!")

    def _version_check(self, version_string, required_version):
        version = re.search(r'\d+\.\d+\.\d+', version_string)
        if not version:
            return False
        return tuple(map(int, version.group().split('.'))) >= tuple(map(int, required_version.split('.')))

    def setup_environment(self):
        print("\nSetting up environment variables...")
        
        if not self.env_sample.exists():
            print("Error: .env.sample file not found")
            sys.exit(1)

        if self.env_file.exists():
            print("Warning: .env file already exists")
            if input("Do you want to overwrite it? (y/n): ").lower() != 'y':
                print("Using existing .env file")
                return

        env_vars = {}
        with open(self.env_sample) as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith('#'):
                    key, default = line.split('=', 1)
                    env_vars[key] = default

        print("\nPlease provide the following configuration values:")
        for key, default in env_vars.items():
            value = input(f"{key} [{default}]: ").strip()
            env_vars[key] = value if value else default

        with open(self.env_file, 'w') as f:
            for key, value in env_vars.items():
                f.write(f"{key}={value}\n")

        print("Environment setup completed!")

    def start_services(self):
        print("\nStarting services...")
        try:
            subprocess.run(["docker-compose", "up", "-d"], check=True, cwd=self.project_root)
            print("Services started successfully!")
        except subprocess.CalledProcessError as e:
            print(f"Error starting services: {e}")
            sys.exit(1)

    def verify_installation(self):
        print("\nVerifying installation...")
        try:
            containers = subprocess.check_output(["docker", "ps"]).decode()
            required_containers = ["nixopus-api-container", "nixopus-db-container", "nixopus-view-container"]
            
            for container in required_containers:
                if container not in containers:
                    print(f"Error: {container} is not running")
                    sys.exit(1)

            print("Installation verified successfully!")
        except subprocess.CalledProcessError as e:
            print(f"Error verifying installation: {e}")
            sys.exit(1)

def main():
    installer = Installer()
    
    print("Nixopus Installation Script")
    print("===========================")
    
    installer.ask_for_sudo()
    email = installer.ask_email()
    password = installer.ask_password()
    domain = installer.ask_domain()
    
    installer.check_system_requirements()
    installer.setup_environment()
    installer.start_services()
    installer.verify_installation()
    
    print("\nInstallation completed successfully!")
    print("You can access the application at:")
    print(f"- API: https://{domain}/api")
    print(f"- View: https://{domain}")
    
    print("\nYou can now login with the following credentials:")
    print(f"Email: {email}")
    print(f"Password: {password}")

if __name__ == "__main__":
    main() 