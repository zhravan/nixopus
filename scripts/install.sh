#!/bin/bash

set -e

# Function to print colored text
print_color() {
    local color=$1
    local text=$2
    case $color in
        red) echo -e "\033[0;31m$text\033[0m" ;;
        green) echo -e "\033[0;32m$text\033[0m" ;;
        yellow) echo -e "\033[0;33m$text\033[0m" ;;
        blue) echo -e "\033[0;34m$text\033[0m" ;;
        purple) echo -e "\033[0;35m$text\033[0m" ;;
        cyan) echo -e "\033[0;36m$text\033[0m" ;;
        *) echo "$text" ;;
    esac
}

print_welcome() {
    echo
    print_color "cyan" "  _   _ _ _                           "
    print_color "cyan" " | \ | (_)                          "
    print_color "cyan" " |  \| |___  _____  _ __  _   _ ___ "
    print_color "cyan" " | . \` | \ \/ / _ \| '_ \| | | / __|"
    print_color "cyan" " | |\  | |>  < (_) | |_) | |_| \__ "
    print_color "cyan" " |_| \_|_/_/\_\___/| .__/ \__,_|___/"
    print_color "cyan" "                   | |              "
    print_color "cyan" "                   |_|              "
    echo
    print_color "yellow" "====================================================="
    print_color "yellow" "           Welcome to Nixopus Installation           "
    print_color "yellow" "====================================================="
    echo
    print_color "green" "This script will help you install Nixopus on your server."
    print_color "green" "Please make sure you have root access and a valid domain name."
    echo
}

# Call welcome message at the start
print_welcome

# detect Linux distribution and version
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo "$ID:$VERSION_ID"
    elif [ -f /etc/redhat-release ]; then
        echo "centos:$(grep -o '[0-9]\+' /etc/redhat-release | head -1)"
    elif [ -f /etc/alpine-release ]; then
        echo "alpine:$(cat /etc/alpine-release)"
    elif [ -f /etc/amazon-linux-release ]; then
        echo "amzn:$(grep -o '[0-9]\+' /etc/amazon-linux-release | head -1)"
    else
        echo "unknown:unknown"
    fi
}

# Function to install Docker based on distribution
install_docker() {
    local distro=$1
    local version=$2
    local arch=$(uname -m)
    
    case $distro in
        ubuntu|debian)
            apt-get update
            apt-get install -y apt-transport-https ca-certificates curl software-properties-common
            curl -fsSL https://download.docker.com/linux/$distro/gpg | apt-key add -
            add-apt-repository "deb [arch=$arch] https://download.docker.com/linux/$distro $(lsb_release -cs) stable"
            apt-get update
            apt-get install -y docker-ce docker-ce-cli containerd.io
            ;;
        centos|rhel|amzn)
            if [ "$distro" = "amzn" ]; then
                amazon-linux-extras install -y docker
            else
                yum install -y yum-utils
                yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
                yum install -y docker-ce docker-ce-cli containerd.io
            fi
            ;;
        fedora)
            dnf -y install dnf-plugins-core
            dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
            dnf install -y docker-ce docker-ce-cli containerd.io
            ;;
        alpine)
            apk add --no-cache docker
            ;;
        *)
            echo "Unsupported distribution: $distro"
            exit 1
            ;;
    esac
    if command -v systemctl >/dev/null 2>&1; then
        systemctl start docker
        systemctl enable docker
    elif command -v rc-update >/dev/null 2>&1; then
        rc-update add docker default
        rc-service docker start
    else
        echo "Warning: Could not determine init system. Docker may need to be started manually."
    fi
}

# Function to install Docker Compose
install_docker_compose() {
    local arch=$(uname -m)
    local compose_version="v2.24.5"
    
    case $arch in
        x86_64) arch="x86_64" ;;
        aarch64|arm64) arch="aarch64" ;;
        armv7l) arch="armv7" ;;
        *) echo "Unsupported architecture: $arch"; exit 1 ;;
    esac
    
    curl -L "https://github.com/docker/compose/releases/download/${compose_version}/docker-compose-$(uname -s)-${arch}" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
}

# Check system requirements
check_system_requirements() {
    local kernel_version=$(uname -r)
    local kernel_major=$(echo $kernel_version | cut -d. -f1)
    local kernel_minor=$(echo $kernel_version | cut -d. -f2)
    
    # Check kernel version (Docker requires 3.10 or higher)
    if [ $kernel_major -lt 3 ] || ([ $kernel_major -eq 3 ] && [ $kernel_minor -lt 10 ]); then
        echo "Error: Kernel version must be 3.10 or higher"
        exit 1
    fi
    
    # Check for required kernel modules
    if ! lsmod | grep -q overlay; then
        echo "Error: Overlay filesystem module not loaded"
        exit 1
    fi
}

# Input validation
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: $0 <domain> <admin_email>"
    exit 1
fi

DOMAIN="$1"
ADMIN_EMAIL="$2"

# Validate email format
if ! [[ "$ADMIN_EMAIL" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
    echo "Invalid email format"
    exit 1
fi

# Check if running as root
if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root"
    exit 1
fi

echo "Starting installation process..."
echo "Domain: $DOMAIN"
echo "Admin Email: $ADMIN_EMAIL"

check_system_requirements

DISTRO_INFO=$(detect_distro)
DISTRO=$(echo $DISTRO_INFO | cut -d: -f1)
VERSION=$(echo $DISTRO_INFO | cut -d: -f2)

echo "Detected distribution: $DISTRO $VERSION"

if ! command -v docker &> /dev/null; then
    echo "Docker not found. Installing Docker..."
    install_docker "$DISTRO" "$VERSION"
fi

if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose not found. Installing Docker Compose..."
    install_docker_compose
fi

if ! docker info &> /dev/null; then
    echo "Docker is not running properly"
    exit 1
fi

echo "Docker and Docker Compose are installed and running"

echo "Installing the application"

echo "Setting up environment variables"

cp .env.sample .env

perl -pi -e "s|DOMAIN|${DOMAIN}|g" .env

if docker compose up --build -d; then
    echo "Started successfully"
else
    echo "Failed to start"
    exit 1
fi

generate_password() {
    tr -dc 'A-Za-z0-9!@#$%^&*()_+{}|:<>?=' < /dev/urandom | head -c 16
}

setup_environment() {
    local domain=$1
    local email=$2
    
    local admin_password=$(generate_password)
    
    if [ ! -f .env.sample ]; then
        echo "Error: .env.sample file not found"
        exit 1
    fi
    
    cp .env.sample .env
    
    perl -pi -e "s|DOMAIN=.*|DOMAIN=${domain}|g" .env
    perl -pi -e "s|ADMIN_EMAIL=.*|ADMIN_EMAIL=${email}|g" .env
    perl -pi -e "s|ADMIN_PASSWORD=.*|ADMIN_PASSWORD=${admin_password}|g" .env
    
    if [ ! -f ssh/id_rsa ]; then
        mkdir -p ssh
        ssh-keygen -t rsa -b 4096 -f ssh/id_rsa -N ""
    fi
    
    perl -pi -e "s|SSH_PRIVATE_KEY=.*|SSH_PRIVATE_KEY=ssh/id_rsa|g" .env
    perl -pi -e "s|SSH_PUBLIC_KEY=.*|SSH_PUBLIC_KEY=ssh/id_rsa.pub|g" .env
    
    echo "Environment setup completed"
    echo "Admin password: ${admin_password}"
}

wait_for_service() {
    local service=$1
    local max_attempts=30
    local attempt=1
    
    echo "Waiting for ${service} to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if docker compose ps | grep -q "${service}.*Up"; then
            echo "${service} is ready"
            return 0
        fi
        echo "Attempt ${attempt}/${max_attempts}: ${service} not ready yet..."
        sleep 10
        attempt=$((attempt + 1))
    done
    
    echo "Timeout waiting for ${service} to be ready"
    return 1
}

create_admin_user() {
    local email=$1
    local password=$2
    
    echo "Creating admin user..."
    
    wait_for_service "app"
    
    echo "Admin user created successfully"
}

setup_environment "$DOMAIN" "$ADMIN_EMAIL"

echo "Starting application containers..."
if ! docker compose up --build -d; then
    echo "Failed to start application containers"
    exit 1
fi

create_admin_user "$ADMIN_EMAIL" "$(grep ADMIN_PASSWORD .env | cut -d= -f2)"

echo "============================================="
echo "Installation completed successfully!"
echo "Application URL: https://${DOMAIN}"
echo "Admin Email: ${ADMIN_EMAIL}"
echo "Admin Password: $(grep ADMIN_PASSWORD .env | cut -d= -f2)"
echo "============================================="
echo "Please save these credentials securely"
echo "============================================="