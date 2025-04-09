import os
import secrets
import string
from pathlib import Path
import subprocess
import socket
import json

class EnvironmentSetup:
    def __init__(self, domain):
        self.domain = domain
        self.project_root = Path(__file__).parent.parent
        self.env_file = self.project_root / ".env"
        self.env_sample = self.project_root / ".env.sample"
        self.config_dir = Path("/etc/nixopus")
        self.docker_certs_dir = self.config_dir / "docker-certs"
        self.ssh_dir = self.config_dir / "ssh"

    def generate_random_string(self, length=12):
        alphabet = string.ascii_letters + string.digits
        return ''.join(secrets.choice(alphabet) for _ in range(length))

    def generate_ssh_key(self):
        self.ssh_dir.mkdir(parents=True, exist_ok=True)
        private_key_path = self.ssh_dir / "id_rsa"
        public_key_path = self.ssh_dir / "id_rsa.pub"
        
        if not private_key_path.exists():
            subprocess.run(["ssh-keygen", "-t", "rsa", "-b", "4096", "-f", str(private_key_path), "-N", ""], check=True)
        
        return private_key_path, public_key_path

    def setup_docker_certs(self):
        self.docker_certs_dir.mkdir(parents=True, exist_ok=True)
        subprocess.run([
            "openssl", "req", "-newkey", "rsa:4096", "-nodes", "-sha256",
            "-keyout", str(self.docker_certs_dir / "key.pem"),
            "-x509", "-days", "365",
            "-out", str(self.docker_certs_dir / "cert.pem"),
            "-subj", "/CN=nixopus"
        ], check=True)

    def get_local_ip(self):
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            s.connect(("8.8.8.8", 80))
            local_ip = s.getsockname()[0]
            s.close()
            return local_ip
        except Exception:
            return "localhost"

    def setup_docker_tls(self):
        docker_config_dir = Path("/etc/docker")
        docker_config_dir.mkdir(parents=True, exist_ok=True)

        daemon_config = {
            "tls": True,
            "tlscacert": str(self.docker_certs_dir / "cert.pem"),
            "tlscert": str(self.docker_certs_dir / "cert.pem"),
            "tlskey": str(self.docker_certs_dir / "key.pem"),
            "hosts": [
                "unix:///var/run/docker.sock",
                f"tcp://0.0.0.0:2376"
            ]
        }

        with open(docker_config_dir / "daemon.json", "w") as f:
            json.dump(daemon_config, f, indent=2)

        try:
            status = subprocess.run(["systemctl", "status", "docker"], capture_output=True, text=True)
            print("\nDocker service status before restart:")
            print(status.stdout)

            subprocess.run(["systemctl", "restart", "docker"], check=True)
            
            status = subprocess.run(["systemctl", "status", "docker"], capture_output=True, text=True)
            print("\nDocker service status after restart:")
            print(status.stdout)

        except subprocess.CalledProcessError as e:
            journal = subprocess.run(["journalctl", "-u", "docker", "-n", "50"], capture_output=True, text=True)
            print("\nDocker service error logs:")
            print(journal.stdout)
            raise Exception(f"Failed to restart Docker service. Error: {e.stderr}\nJournal logs: {journal.stdout}")

    def setup_environment(self):
        db_name = f"nixopus_{self.generate_random_string(8)}"
        username = f"nixopus_{self.generate_random_string(8)}"
        password = self.generate_random_string(16)
        
        api_port = 8443
        next_public_port = 7443 
        db_port = 5432 

        private_key_path, public_key_path = self.generate_ssh_key()
        local_ip = self.get_local_ip()

        self.setup_docker_certs()
        self.setup_docker_tls()

        env_vars = {
            "DB_NAME": db_name,
            "USERNAME": username,
            "PASSWORD": password,
            "HOST_NAME": "nixopus-db",
            "DB_PORT": str(db_port),
            "SSL_MODE": "disable",
            "API_PORT": str(api_port),
            "NEXT_PUBLIC_BASE_URL": f"https://api.{self.domain}/api",
            "NEXT_PUBLIC_WEBSOCKET_URL": f"wss://api.{self.domain}/ws",
            "NEXT_PUBLIC_PORT": str(next_public_port),
            "MOUNT_PATH": "/etc/nixopus/configs",
            "PORT": str(api_port),
            "SSH_HOST": local_ip,
            "SSH_PORT": "22",
            "SSH_USER": "root",
            "SSH_PRIVATE_KEY": str(private_key_path),
            "DOCKER_HOST": f"tcp://{local_ip}:2376",
            "DOCKER_TLS_VERIFY": "1",
            "DOCKER_CERT_PATH": str(self.docker_certs_dir),
            "CADDY_ENDPOINT": "http://localhost:2019"
        }

        with open(self.env_file, 'w') as f:
            for key, value in env_vars.items():
                f.write(f"{key}={value}\n")

        # copy to api/.env
        api_env_file = self.project_root / "api" / ".env"
        with open(api_env_file, 'w') as f:
            for key, value in env_vars.items():
                f.write(f"{key}={value}\n")

        self.env_file.chmod(0o600)
        private_key_path.chmod(0o600)
        public_key_path.chmod(0o644)
        self.docker_certs_dir.chmod(0o700)
        for cert_file in self.docker_certs_dir.glob("*"):
            cert_file.chmod(0o600)

        return env_vars 