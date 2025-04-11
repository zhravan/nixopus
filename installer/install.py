#!/usr/bin/env python3

import os
import sys
import subprocess
import platform
from pathlib import Path
import re
from validation import Validation
from environment import EnvironmentSetup
import shutil
import json
import secrets
import time

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
    
    def generate_strong_password(self):
        return secrets.token_urlsafe(16)

    def ask_admin_credentials(self):
        email = input("Please enter the email for the admin: ")
        validation = Validation()
        validation.validate_email(email)
        password = input("Please enter the password for the admin(generates a strong password if left blank): ")
        if not password:
            password = self.generate_strong_password()
        validation.validate_password(password)
        return email, password
    
    def check_docker_version(self):
        try:
            subprocess.run(["docker", "--version"], check=True, capture_output=True)
        except subprocess.CalledProcessError as e:
            print(f"Error: Docker version {self.required_docker_version} or higher is required")
            print(e.stderr.decode())
            sys.exit(1)
            
    def check_docker_compose_version(self):
        try:
            subprocess.run(["docker-compose", "--version"], check=True, capture_output=True)
        except subprocess.CalledProcessError as e:
            print(f"Error: Docker Compose version {self.required_compose_version} or higher is required")
            print(e.stderr.decode())
            sys.exit(1)
            
    def check_curl_installed(self):
        if not shutil.which("curl"):
            print("Error: Curl is not installed")
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
        self.check_curl_installed()

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
            os.environ["DOCKER_HOST"] = "tcp://localhost:2376"
            os.environ["DOCKER_TLS_VERIFY"] = "1"
            os.environ["DOCKER_CERT_PATH"] = "/etc/nixopus/docker-certs"
            
            compose_cmd = ["docker", "compose"] if shutil.which("docker") else ["docker-compose"]
            result = subprocess.run(compose_cmd + ["up", "--build", "-d"], capture_output=True, text=True, cwd=self.project_root)
            
            if result.returncode != 0:
                print("Error starting services:")
                print(result.stderr)
                sys.exit(1)
        except Exception as e:
            print(f"Error starting services: {str(e)}")
            sys.exit(1)

    def verify_installation(self):
        print("\nVerifying installation...")
        try:
            result = subprocess.run(["docker", "ps"], capture_output=True, text=True)
            if result.returncode != 0:
                print("Error verifying installation:")
                print(result.stderr)
                sys.exit(1)
                
            containers = result.stdout
            required_containers = ["nixopus-api-container", "nixopus-db-container", "nixopus-view-container"]
            
            for container in required_containers:
                if container not in containers:
                    print(f"Error: {container} is not running")
                    sys.exit(1)

            print("Installation verified successfully!")
        except Exception as e:
            print(f"Error verifying installation: {str(e)}")
            sys.exit(1)
    
    def setup_caddy(self):
        print("\nSetting up Proxy...")
        try:
            with open('api/helpers/caddy.json', 'r') as f:
                config = json.dumps(json.load(f))
            
            result = subprocess.run(
                ['curl', '-X', 'POST', 'http://localhost:2019/load',
                 '-H', 'Content-Type: application/json',
                 '-d', config],
                capture_output=True,
                text=True
            )
            
            if result.returncode == 0:
                print("✓ Caddy configuration loaded successfully")
            else:
                print("✗ Failed to load Caddy configuration:")
                print(result.stderr)
        except Exception as e:
            print(f"✗ Error setting up Caddy: {str(e)}")
    
    def setup_admin(self, email, password, domain):
        print("\nSetting up admin...")
        username = email.split('@')[0]
        
        raw_command = f'curl -X POST https://api.{domain}/api/v1/auth/register -H "Content-Type: application/json" --connect-timeout 10 --max-time 30 --insecure -d \'{{"email":"{email}","password":"{password}","type":"admin","username":"{username}","organization":""}}\''
        print("Raw command:", raw_command)
        
        curl_command = [
            "curl", "-X", "POST", f"https://api.{domain}/api/v1/auth/register",
            "-H", "Content-Type: application/json",
            "--connect-timeout", "10",
            "--max-time", "30",
            "--insecure",
            "-d", f'{{"email":"{email}","password":"{password}","type":"admin","username":"{username}","organization":""}}'
        ]
        
        result = subprocess.run(curl_command, capture_output=True, text=True, check=False)
        print("Response:", result.stdout)
        
        try:
            response = json.loads(result.stdout)
            if response.get("status") == "success":
                print("✓ Admin setup completed successfully")
                return
            raise Exception(f"API Error: {response.get('message', 'Unknown error')}")
        except json.JSONDecodeError:
            raise Exception("Invalid response from API")

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
    print("\n")
    print("\033[1mWelcome to Nixopus Installation Wizard\033[0m")
    print("This wizard will guide you through the installation process of Nixopus.")
    print("Please follow the prompts carefully to complete the setup.\n")
    
    installer.ask_for_sudo()
    domain = installer.ask_domain()
    email, password = installer.ask_admin_credentials()
    installer.check_system_requirements()
    env_vars = installer.setup_environment(domain)
    installer.start_services()
    installer.verify_installation()
    installer.setup_caddy()
    
    try:
        installer.setup_admin(email, password, domain)
    except Exception as e:
        print(f"✗ {str(e)}")
        sys.exit(1)
    
    print("\n\033[1mInstallation Complete!\033[0m")
    print("\n\033[1mAccess Information:\033[0m")
    print(f"• Nixopus is now installed and running on: {domain}")
    
    print("\n\033[1mAdmin Credentials:\033[0m")
    print(f"• Email: {email}")
    print(f"• Password: {password}")
    print("\n\033[1mImportant:\033[0m Please save these credentials securely. You will need them to log in.")
    print("\n\033[1mThank you for installing Nixopus!\033[0m")
    print("\n\033[1mPlease visit the documentation at https://docs.nixopus.com for more information.\033[0m")
    print("\n\033[1mIf you have any questions, please visit the community forum at https://community.nixopus.com\033[0m")
    print("\n\033[1mSee you in the community!\033[0m")

if __name__ == "__main__":
    main() 