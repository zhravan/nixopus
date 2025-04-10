#!/usr/bin/env python3

import os
import sys
import subprocess
import platform
from pathlib import Path
import re
from validation import Validation
from environment import EnvironmentSetup

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
    
    def ask_domain(self):
        domain = input("Please enter the domain in which you want to run the application (e.g. nixopus.com): ")
        validation = Validation()
        validation.validate_domain(domain)
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

    def setup_environment(self, domain):
        print("\nSetting up environment...")
        env_setup = EnvironmentSetup(domain)
        env_vars = env_setup.setup_environment()
        print("Environment setup completed!")
        return env_vars

    def start_services(self):
        print("\nStarting services...")
        try:
            subprocess.run(["docker-compose", "up", "--build", "-d"], check=True, cwd=self.project_root)
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
    
    print("\033[36m  _   _ _ _                           \033[0m")
    print("\033[36m | \ | (_)                          \033[0m")
    print("\033[36m |  \| |___  _____  _ __  _   _ ___ \033[0m")
    print("\033[36m | . \` | \ \/ / _ \| '_ \| | | / __|\033[0m")
    print("\033[36m | |\  | |>  < (_) | |_) | |_| \__ \033[0m")
    print("\033[36m |_| \_|_/_/\_\___/| .__/ \__,_|___/\033[0m")
    print("\033[36m                   | |              \033[0m")
    print("\033[36m                   |_|              \033[0m")
    print("")
    print("Welcome to the Nixopus installer!")
    print("This script will install Nixopus on your system. Hold on tight!")
    installer.ask_for_sudo()
    domain = installer.ask_domain()
    
    installer.check_system_requirements()
    env_vars = installer.setup_environment(domain)
    installer.start_services()
    installer.verify_installation()
    
    print("\nInstallation completed successfully!")
    print("You can access the application at:")
    print(f"- API: https://api.{domain}/api")
    print(f"- View: https://{domain}")
    
    print("\nGenerated credentials:")
    print(f"Database:")
    print(f"  Name: {env_vars['DB_NAME']}")
    print(f"  Username: {env_vars['USERNAME']}")
    print(f"  Password: {env_vars['PASSWORD']}")
    print(f"SSH:")
    print(f"  Username: {env_vars['SSH_USER']}")
    print(f"  Private Key: {env_vars['SSH_PRIVATE_KEY']}")

if __name__ == "__main__":
    main() 