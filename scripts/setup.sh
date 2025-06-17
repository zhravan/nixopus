#!/usr/bin/env bash

# This script is used to setup the environment for the project

set -euo pipefail

BRANCH="feat/dev_environment"

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
        echo "Error: Command '$cmd' not found. Please install it manually" >&2
        exit 1
    fi
}

function check_required_commands() {
    local commands=("git" "docker" "yarn" "go")  
    for cmd in "${commands[@]}"; do
        check_command "$cmd"
    done
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

# install air hot reload for golang
function install_air_hot_reload(){
    local user_home
    user_home=$(eval echo ~${SUDO_USER:-$USER})
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
    user_home=$(eval echo ~${SUDO_USER:-$USER})
    
    nohup "$user_home/go/bin/air" > api.log 2>&1 &
}

# start the view server
function start_view(){
    move_to_folder "view"
    yarn install --frozen-lockfile
    nohup yarn run dev > view.log 2>&1 &
}

# main function
function main() {
    echo "Starting Nixopus development environment setup..."
    check_root
    check_os
    check_required_commands
    clone_nixopus
    move_to_folder "nixopus"
    checkout_branch "$BRANCH"
    install_air_hot_reload
    echo "Nixopus repository cloned and configured successfully"
    
    setup_postgres_with_docker
    setup_environment_variables
    setup_ssh
    echo "SSH setup completed successfully"
    
    start_api
    echo "API server started successfully"
    move_to_folder ".."
    start_view
    echo "View server started successfully"
    echo "Nixopus development environment setup completed successfully"
}

main
