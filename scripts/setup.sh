#!/usr/bin/env bash

# Nixopus Development Environment Setup Script
# 
# This script sets up the development environment for the project
# Supported platforms: Linux (Ubuntu, CentOS, Fedora, Arch) and macOS

# Prerequisites:
# - Linux: Run with sudo privileges
# - macOS: Homebrew should be installed (https://brew.sh)
#         Docker Desktop for Mac should be installed and running

# Usage:
# - Linux: sudo ./setup.sh
# - macOS: ./setup.sh (no sudo required)

set -euo pipefail


BRANCH="feat/dev_environment"
OS="$(uname)"

function detect_package_manager() {
    if [[ "$OS" == "Darwin" ]]; then
        if command -v brew &>/dev/null; then
            echo "brew"
        else
            echo "Error: Homebrew not found. Please install Homebrew first: https://brew.sh" >&2
            exit 1
        fi
    elif command -v apt-get &>/dev/null; then
        echo "apt"
    elif command -v dnf &>/dev/null; then
        echo "dnf"
    elif command -v yum &>/dev/null; then
        echo "yum"
    elif command -v pacman &>/dev/null; then
        echo "pacman"
    else
        echo "Error: Unsupported package manager" >&2
        exit 1
    fi
}

function install_package() {
    local pkg_manager
    pkg_manager=$(detect_package_manager)
    
    case $pkg_manager in
        "brew")
            brew install "$1"
            ;;
        "apt")
            apt-get update
            apt-get install -y "$1"
            ;;
        "dnf")
            dnf install -y "$1"
            ;;
        "yum")
            yum install -y "$1"
            ;;
        "pacman")
            pacman -Sy --noconfirm "$1"
            ;;
    esac
}

# check if the os is linux or macOS
function check_os() {
    if [[ "$OS" != "Linux" && "$OS" != "Darwin" ]]; then
        echo "Error: This script is only supported on Linux and macOS." >&2
        exit 1
    fi
}

# check for Docker availability globally
function check_docker() {
    if ! command -v docker &>/dev/null; then
        if [[ "$OS" == "Darwin" ]]; then
            echo "Error: Docker not found. Please ensure Docker Desktop for Mac is installed and running."
            echo "Download from: https://www.docker.com/products/docker-desktop"
        else
            echo "Error: Docker not found. Please install Docker first."
            echo "You can install it using your package manager or from: https://docs.docker.com/engine/install/"
        fi
        exit 1
    fi
    
    # Check if Docker daemon is running
    if ! docker info &>/dev/null; then
        echo "Error: Docker daemon is not running. Please start Docker service."
        exit 1
    fi
    
    echo "Docker check completed."
}

# check for prerequisites on macOS
function check_macos_prerequisites() {
    if [[ "$OS" == "Darwin" ]]; then
        echo "Checking macOS prerequisites"
        
        # Check for Homebrew
        if ! command -v brew &>/dev/null; then
            echo "Error: Homebrew is required on macOS but not found." >&2
            echo "Please install Homebrew first by running:" >&2
            echo '/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"' >&2
            exit 1
        fi
        
        echo "macOS prerequisites check completed."
    fi
}

# check if the script is running as root (only required for Linux)
function check_root() {
    if [[ "$OS" == "Linux" && "$EUID" -ne 0 ]]; then
        echo "Error: Please run as root (sudo) on Linux systems" >&2
        exit 1
    elif [[ "$OS" == "Darwin" && "$EUID" -eq 0 ]]; then
        echo "Warning: consider running without sudo on macos as it is not equired." >&2
    fi
}

# check if the required commands are installed
function check_command() {
    local cmd="$1"
    if ! command -v "$cmd" &>/dev/null; then
        echo "Command '$cmd' not found. Attempting to install"
        case "$cmd" in
            "git")
                install_package "git"
                ;;
            "yarn")
                if [[ "$OS" == "Darwin" ]]; then
                    install_package "node"
                    install_package "yarn"
                else
                    install_package "nodejs"
                    install_package "npm"
                    npm install -g yarn
                fi
                ;;
            *)
                echo "Error: Automatic installation not supported for '$cmd'" >&2
                exit 1
                ;;
        esac
    fi
}

function check_required_commands() {
    local commands=("git" "yarn" "go")  
    for cmd in "${commands[@]}"; do
        check_command "$cmd"
    done
}

# check if the go version is installed 
function check_go_version() {
    local go_version=$(go version | awk '{print $3}' | sed 's/^go//')
    local required_version="1.23.4"
    
    local ver_num=$(echo "$go_version" | sed 's/\.//g')
    local req_num=$(echo "$required_version" | sed 's/\.//g')
    
    if [ "$ver_num" -lt "$req_num" ]; then
        echo "Error: Go version $required_version or higher is required. Current version: $go_version" >&2
        exit 1
    fi
}

# clone the nixopus repository
function clone_nixopus() {
    if [ -d "nixopus" ]; then
        echo "nixopus directory already exists, please remove it manually and run the script again"
        exit 1
    fi
    if ! git clone --branch "$BRANCH" https://github.com/raghavyuva/nixopus.git; then
        echo "Error: Failed to clone nixopus repository" >&2
        exit 1
    fi
}   

# checkout to the branch
function checkout_branch() {
    local branch="$1"
    if ! git checkout "$branch"; then
        echo "Error: Failed to checkout to $branch" >&2
        exit 1
    fi
}

# move to the folder
function move_to_folder() {
    local folder="$1"
    if ! cd "$folder"; then
        echo "Error: Failed to change directory to $folder" >&2
        exit 1
    fi
}

# get architecture name in format expected by most download URLs
function get_arch_name() {
    local arch
    arch=$(uname -m)
    case "$arch" in
        x86_64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
    esac
}

# get OS name in lowercase format (for downloads, packages, etc.)
function get_os_name() {
    case "$OS" in
        "Linux") echo "linux" ;;
        "Darwin") echo "darwin" ;;
        *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
    esac
}

function install_go() {
    local version="1.23.4"
    local arch
    arch=$(get_arch_name)
    
    local os
    os=$(get_os_name)
    
    local temp_dir
    temp_dir=$(mktemp -d)
    
    echo "Downloading Go ${version}"
    if ! curl -L "https://go.dev/dl/go${version}.${os}-${arch}.tar.gz" -o "${temp_dir}/go.tar.gz"; then
        echo "Error: Failed to download Go" >&2
        rm -rf "$temp_dir"
        exit 1
    fi

    echo "Verifying checksum"
    local checksum_url="https://go.dev/dl/go${version}.${os}-${arch}.tar.gz.sha256"
    local expected_sum
    expected_sum=$(curl -sL "$checksum_url" | awk '{print $1}')
    local actual_sum
    if [[ "$OS" == "Darwin" ]]; then
        actual_sum=$(shasum -a 256 "${temp_dir}/go.tar.gz" | awk '{print $1}')
    else
        actual_sum=$(sha256sum "${temp_dir}/go.tar.gz" | awk '{print $1}')
    fi
    
    if [[ $expected_sum != $actual_sum ]]; then
        echo "Error: Checksum mismatch for Go archive" >&2
        rm -rf "$temp_dir"
        exit 1
    fi
    
    echo "Installing Go ${version}"
    local go_install_path
    if [[ "$OS" == "Darwin" ]]; then #incase of macos
        go_install_path="/usr/local"
        if [[ "$EUID" -ne 0 ]]; then
            echo "Installing Go to user directory"
            go_install_path="$HOME"
        fi
    else
        go_install_path="/usr/local"
    fi
    
    if ! rm -rf "${go_install_path}/go" && tar -C "${go_install_path}" -xzf "${temp_dir}/go.tar.gz"; then
        echo "Error: Failed to install Go" >&2
        rm -rf "$temp_dir"
        exit 1
    fi
    
    rm -rf "$temp_dir"
    
    # Set up PATH
    if [[ "$OS" == "Darwin" ]]; then
        local shell_profile
        if [[ "$SHELL" == *"zsh"* ]]; then
            shell_profile="$HOME/.zshrc"
        else
            shell_profile="$HOME/.bash_profile"
        fi
        
        if ! grep -q "${go_install_path}/go/bin" "$shell_profile" 2>/dev/null; then
            echo "export PATH=\$PATH:${go_install_path}/go/bin" >> "$shell_profile"
            echo "Added Go to PATH in $shell_profile"
        fi
        export PATH="$PATH:${go_install_path}/go/bin"
    else
        if ! grep -q "/usr/local/go/bin" /etc/profile.d/go.sh 2>/dev/null; then
            echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile.d/go.sh
            chmod +x /etc/profile.d/go.sh
            source /etc/profile.d/go.sh
        fi
    fi
}

# check if the go version is installed else install it 1.23.4
function check_go_version() {
    if ! command -v go &>/dev/null; then
        echo "Go is not installed. Installing"
        install_go
    fi
}

# install air hot reload for golang
function install_air_hot_reload(){
    local user_home
    if [[ "$OS" == "Darwin" ]]; then
        user_home="$HOME"
        if [[ "$EUID" -eq 0 && -n "${SUDO_USER:-}" ]]; then
            user_home=$(eval echo "~$SUDO_USER")
        fi
    else
        user_home=$(eval echo ~${SUDO_USER:-$USER})
    fi
    
    sudo -u "${SUDO_USER:-$USER}" env GOPATH="$user_home/go" go install github.com/air-verse/air@latest
    export PATH="$PATH:$user_home/go/bin"
}

# load the env variables from the api/.env.sample file
function load_api_env_variables(){
    move_to_folder "api"
    if [ -f .env.sample ]; then
        while IFS='=' read -r key value; do
            # Skip empty lines and comments
            [[ -z "$key" || "$key" =~ ^# ]] && continue
            # Remove any quotes from the value
            value=$(echo "$value" | tr -d '"'"'")
            # Export the variable
            export "$key=$value"
        done < .env.sample
    else
        echo "Error: .env.sample file not found in api directory" >&2
        exit 1
    fi
    move_to_folder ".."
}

# setup postgres with docker
function setup_postgres_with_docker(){
    load_api_env_variables

    # Check if container already exists
    if docker ps -a --format 'table {{.Names}}' | grep -q "^nixopus-db$"; then
        echo "Already nixopus-db container exists"
        return 0
    fi
    
    # Start PostgreSQL container with credentials matching .env.sample
    docker run -d --name nixopus-db \
        -e POSTGRES_USER="${USERNAME:-postgres}" \
        -e POSTGRES_PASSWORD="${PASSWORD:-12344}" \
        -e POSTGRES_DB="${DB_NAME:-postgres}" \
        -e POSTGRES_HOST_AUTH_METHOD=trust \
        -p "${DB_PORT:-5432}:5432" \
        --health-cmd="pg_isready -U ${USERNAME:-postgres} -d ${DB_NAME:-postgres}" \
        postgres:14-alpine
    
    echo "Waiting for PostgreSQL to be ready"
    sleep 5
    
    # Wait for PostgreSQL to be ready
    local max_attempts=30
    local attempt=1
    while [ $attempt -le $max_attempts ]; do
        if docker exec nixopus-db pg_isready -U "${USERNAME:-postgres}" -d "${DB_NAME:-postgres}" >/dev/null 2>&1; then
            echo "PostgreSQL is ready!"
            break
        fi
        echo "Waiting for PostgreSQL (attempt $attempt/$max_attempts)"
        sleep 2
        attempt=$((attempt + 1))
    done
    
    if [ $attempt -gt $max_attempts ]; then
        echo "Error: PostgreSQL failed to start within expected time" >&2
        exit 1
    fi
    
    echo "Postgres setup completed successfully"
}

# verify database connection
function verify_database_connection(){
    echo "Verifying database connection"
    load_api_env_variables
    
    # Test connection using docker exec
    if docker exec nixopus-db psql -U "${USERNAME:-postgres}" -d "${DB_NAME:-postgres}" -c "SELECT 1;" >/dev/null 2>&1; then
        echo "Database connection verified successfully"
    else
        echo "Error: Failed to connect to database" >&2
        echo "Details: user=${USERNAME:-postgres}, db=${DB_NAME:-postgres}, host=${HOST_NAME:-localhost}, port=${DB_PORT:-5432}" >&2
        exit 1
    fi
}

# setup ssh will create a ssh key and add it to the authorized_keys file
function setup_ssh(){
    local ssh_dir="$HOME/.ssh"
    local private_key="$ssh_dir/id_ed25519_nixopus"
    local public_key="$ssh_dir/id_ed25519_nixopus.pub"
    
    if [[ -f "$private_key" && -f "$public_key" ]]; then
        echo "SSH key for Nixopus already exists, skipping ssh setup"
        return 0
    fi
    
    echo "setting up SSH config"
    
    # Check if ssh-keygen is available and install if needed
    if ! command -v ssh-keygen &>/dev/null; then
        echo "Installing openssh"
        if [[ "$OS" == "Darwin" ]]; then
            install_package "openssh"
        else
            case $(detect_package_manager) in
                "apt") install_package "openssh-client" ;;
                "dnf"|"yum") install_package "openssh-clients" ;;
                "pacman") install_package "openssh" ;;
            esac
        fi
    fi
    
    local authorized_keys="$ssh_dir/authorized_keys"

    mkdir -p "$ssh_dir" && chmod 700 "$ssh_dir"
    
    echo "Generating Nixopus SSH key"
    ssh-keygen -t ed25519 -f "$private_key" -N "" -C "nixopus-$(whoami)@$(hostname)-$(date +%Y%m%d)"
    chmod 600 "$private_key" && chmod 644 "$public_key"
    
    if [[ ! -f "$authorized_keys" || ! $(grep -Fq "$(cat "$public_key")" "$authorized_keys" 2>/dev/null) ]]; then
        cat "$public_key" >> "$authorized_keys" && chmod 600 "$authorized_keys"
        echo "Nixopus public key added to authorized_keys"
    fi
    
    echo "Nixopus SSH setup is done"
}

# setup environment variables
function setup_environment_variables(){
    move_to_folder "api"
    if [ -f .env.sample ]; then
        cp .env.sample .env || { echo "Error: Failed to copy api/.env.sample to .env" >&2; exit 1; }
    else
        echo "Error: api/.env.sample file not found" >&2
        exit 1
    fi
    move_to_folder ".."
    
    move_to_folder "view"
    if [ -f .env.sample ]; then
        cp .env.sample .env || { echo "Error: Failed to copy view/.env.sample to .env" >&2; exit 1; }
    else
        echo "Error: view/.env.sample file not found" >&2
        exit 1
    fi
    move_to_folder ".."
    echo "Environment variables setup completed successfully"
}

# start the api server
function start_api(){
    move_to_folder "api"
    go mod tidy
    go mod download
    
    local user_home
    if [[ "$OS" == "Darwin" ]]; then
        user_home="$HOME"
        if [[ "$EUID" -eq 0 && -n "${SUDO_USER:-}" ]]; then
            user_home=$(eval echo "~$SUDO_USER")
        fi
    else
        user_home=$(eval echo ~${SUDO_USER:-$USER})
    fi
    
    echo "API server started with air hot reload"
    echo "Logs can be found in api.log"
    echo "You can stop the server using 'pkill -f air' command"
    nohup "$user_home/go/bin/air" > api.log 2>&1 &
    
}

open_discord_gh_link() {
  local url="https://discord.com/invite/skdcq39Wpv"
  local gh_url="https://github.com/raghavyuva/nixopus/"

  case "$OS" in
    Darwin)
      open "$url"
      open "$gh_url"
      ;;
    Linux)
      if command -v xdg-open &>/dev/null; then
        xdg-open "$url"
        xdg-open "$gh_url"
      else
        echo "Warning: Could not auto-open browser." >&2
      fi
      ;;
    *)
      echo "Warning: Unsupported OS for browser launch." >&2
      ;;
  esac
}


# start the view server
function start_view(){
    move_to_folder "view"
    yarn install --frozen-lockfile
    echo "View server started"
    echo "Logs can be found in view.log"
    echo "You can stop the server using 'pkill -f yarn' command"
    nohup yarn run dev > view.log 2>&1 &

}

# check if required ports are available
# This function checks if the essential ports needed by Nixopus services are free:
# - 5432: PostgreSQL database (configurable via DB_PORT in api/.env.sample)
# - 8080: API server (configurable via PORT in api/.env.sample)  
# - 3000: Frontend development server (Next.js default for 'yarn dev')
check_port_availability() {
  for p in 5432 8080 3000; do
    if pid=$(lsof -ti TCP:"$p"); then
      echo "Port $p is in use (PID $pid)" >&2
      exit 1
    fi
  done
  echo "All ports (5432, 8080, 3000) are free."
}

# main function
function main() {
    echo "Starting Nixopus development environment setup"
    check_os
    check_macos_prerequisites
    check_root
    check_docker
    check_required_commands
    check_port_availability
    check_go_version
    clone_nixopus
    move_to_folder "nixopus"
    checkout_branch "$BRANCH"
    install_air_hot_reload
    echo "Nixopus repository cloned and configured successfully"
    
    setup_postgres_with_docker
    verify_database_connection
    setup_environment_variables
    setup_ssh
    echo "SSH setup completed successfully"
    
    start_api
    echo "API server started successfully"
    move_to_folder ".."
    start_view
    echo "View server started successfully"
    echo "Nixopus development environment setup completed successfully"
    echo "-------------------------------------------------------------"    
    echo ""
    echo "=== Application Access ==="
    echo "Frontend: http://localhost:3000"
    echo "API: http://localhost:8080"
    echo "Database: localhost:5432"
    echo ""
    echo "=== Troubleshooting ==="
    echo "If you encounter database connection issues:"
    echo "1. Check if Docker container is running: docker ps | grep nixopus-db"
    echo "2. Check database logs: docker logs nixopus-db"
    echo "3. Verify connection: docker exec nixopus-db psql -U postgres -d postgres -c 'SELECT 1;'"
    echo "4. Restart the database: docker restart nixopus-db"
    echo ""
    echo "Log files:"
    echo "- API logs: nixopus/api/api.log"
    echo "- View logs: nixopus/view/view.log"
    echo "----------------------------------------------------------------------------"
    
    echo ""
    echo "Need help or have questions?"
    echo ">>>> Join our Discord :: https://discord.com/invite/skdcq39Wpv"
    echo ">>>> Star us on GitHub: https://github.com/raghavyuva/nixopus/"
    echo ">>>> Raise issues on GitHub Issues: https://github.com/raghavyuva/nixopus/issues"
    open_discord_gh_link

}

main
