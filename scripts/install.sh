#!/bin/bash

# check if the script is running as root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root (sudo)"
    exit 1
fi

# check if the required commands are installed
function check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "$1 is not installed. Please install $1 before running this script."
        exit 1
    fi
}

# check if the required python version is installed
function check_python_version() {
    if ! python3 --version | grep -q "Python 3.10"; then
        echo "Python 3.10 is not installed. Please install Python 3.10 before running this script."
        # we will have to install python 3.10 here based on the OS
        exit 1
    fi
}

# check if the required commands are installed
check_command "python3"
check_command "pip3"
check_command "git"

check_python_version

# default to production environment if no environment is specified through the command line arguments
ENV="production"

# parse the command line arguments
for arg in "$@"; do
    if [[ $arg == "--env"* ]]; then
        ENV=$(echo $arg | cut -d'=' -f2)
        break
    fi
done

# set the source and destination directories based on the environment
if [ "$ENV" == "staging" ]; then
    NIXOPUS_DIR="/etc/nixopus-staging"
    SOURCE_DIR="$NIXOPUS_DIR/source"
    BRANCH="feat/develop"
else
    NIXOPUS_DIR="/etc/nixopus"
    SOURCE_DIR="$NIXOPUS_DIR/source"
    BRANCH="master"
fi

# create the directories if they don't exist
mkdir -p $NIXOPUS_DIR
mkdir -p $SOURCE_DIR

# clone the repository if it doesn't exist
if [ -d "$SOURCE_DIR/.git" ]; then
    cd $SOURCE_DIR
    git fetch --all
    git reset --hard origin/$BRANCH
    git checkout $BRANCH
    git pull
else
    rm -rf $SOURCE_DIR/* $SOURCE_DIR/.[!.]*
    git clone https://github.com/raghavyuva/nixopus.git $SOURCE_DIR
    cd $SOURCE_DIR
    git checkout $BRANCH
fi

# set up the caddy configuration
echo "Setting up Caddy configuration..."
rm -rf $NIXOPUS_DIR/caddy > /dev/null 2>&1
mkdir -p $NIXOPUS_DIR/caddy > /dev/null 2>&1
echo '{
	admin 0.0.0.0:2019
	log {
		format json
		level INFO
	}
}' > $NIXOPUS_DIR/caddy/Caddyfile

cd $SOURCE_DIR/installer
echo "Setting up Nixopus Installation Environment..."
python3 -m venv venv > /dev/null 2>&1
source venv/bin/activate > /dev/null 2>&1

echo "Installing dependencies..."
pip install --upgrade pip > /dev/null 2>&1
pip install -r requirements.txt > /dev/null 2>&1

echo "Starting Nixopus Installation..."
# run the installer
PYTHONPATH=$SOURCE_DIR/installer python3 install.py "$@"

# deactivate the virtual environment
deactivate > /dev/null 2>&1