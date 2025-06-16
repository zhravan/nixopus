#! /bin/bash

# This script is used to setup the environment for the project

set -euo pipefail

function detect_package_manager() {
    if command -v apt-get &>/dev/null; then
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

# check if the os is linux
function check_os() {
    if [ "$(uname)" != "Linux" ]; then
        echo "Error: This script is only supported on Linux." >&2
        exit 1
    fi
}

# check if the script is running as root
function check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo "Error: Please run as root (sudo)" >&2
        exit 1
    fi
}

# check if the required commands are installed
function check_command() {
    local cmd="$1"
    if ! command -v "$cmd" &>/dev/null; then
        echo "Error: '$cmd' is not installed. Please install '$cmd' before running this script." >&2
        exit 1
    fi
}  

function check_required_commands() {
    local commands=("git" "docker" "yarn")
    for cmd in "${commands[@]}"; do
        check_command "$cmd"
    done
}

# clone the nixopus repository
function clone_nixopus() {
    if [ -d "nixopus" ]; then
        echo "nixopus directory already exists. Removing..."
        rm -rf nixopus
    fi
    if ! git clone https://github.com/raghavyuva/nixopus.git; then
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

function install_go() {
    local version="1.23.4"
    local arch
    arch=$(uname -m)
    case "$arch" in
        x86_64) arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
    esac
    local os="linux"
    local temp_dir
    temp_dir=$(mktemp -d)
    
    echo "Downloading Go ${version}..."
    if ! curl -L "https://go.dev/dl/go${version}.${os}-${arch}.tar.gz" -o "${temp_dir}/go.tar.gz"; then
        echo "Error: Failed to download Go" >&2
        rm -rf "$temp_dir"
        exit 1
    fi
    
    echo "Installing Go ${version}..."
    if ! rm -rf /usr/local/go && tar -C /usr/local -xzf "${temp_dir}/go.tar.gz"; then
        echo "Error: Failed to install Go" >&2
        rm -rf "$temp_dir"
        exit 1
    fi
    
    rm -rf "$temp_dir"
    
    if ! grep -q "/usr/local/go/bin" /etc/profile.d/go.sh 2>/dev/null; then
        echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
        chmod +x /etc/profile.d/go.sh
        source /etc/profile.d/go.sh
    fi
}

# check if the go version is installed else install it 1.23.4
function check_go_version() {
    if ! command -v go &>/dev/null; then
        echo "Go is not installed. Installing..."
        install_go
    fi
}

# install air hot reload for golang
function install_air_hot_reload(){
    go install github.com/air-verse/air@latest
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
    docker run -d --name nixopus-db \
        -e POSTGRES_USER="${USERNAME:-nixopus}" \
        -e POSTGRES_PASSWORD="${PASSWORD:-nixopus}" \
        -e POSTGRES_DB="${DB_NAME:-nixopus}" \
        -p "${DB_PORT:-5432}:5432" \
        postgres
    echo "Postgres setup completed successfully"
}

# setup ssh will create a ssh key and add it to the authorized_keys file
function setup_ssh(){
    # TODO: generate SSH key and add to authorized_keys
    return 0
}

# setup environment variables
function setup_environment_variables(){
    move_to_folder "api"
    cp .env.sample .env
    move_to_folder ".."
    move_to_folder "view"
    cp .env.sample .env
    move_to_folder ".."
    echo "Environment variables setup completed successfully"
}

# start the api server
function start_api(){
    move_to_folder "api"
    go mod tidy
    go mod download
    air 
}

# start the view server
function start_view(){
    move_to_folder "view"
    yarn install --frozen-lockfile
    yarn run dev
}

# main function
function main() {
    echo "Starting Nixopus development environment setup..."
    check_root
    check_os
    check_required_commands
    check_go_version
    clone_nixopus
    move_to_folder "nixopus"
    echo "Nixopus repository cloned successfully"
    setup_postgres_with_docker
    echo "Postgres setup completed successfully"
    setup_environment_variables
    echo "Environment variables setup completed successfully"
    install_air_hot_reload
    echo "Air hot reload installed successfully"
    setup_ssh
    echo "SSH setup completed successfully"
    start_api
    echo "API server started successfully"
    start_view
    echo "View server started successfully"
    echo "Nixopus development environment setup completed successfully"
}

main