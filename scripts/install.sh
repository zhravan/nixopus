#!/bin/bash

if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root (sudo)"
    exit 1
fi

if ! command -v python3 &> /dev/null; then
    echo "Python 3 is not installed. Please install Python 3 before running this script."
    exit 1
fi

if ! command -v pip3 &> /dev/null; then
    echo "pip3 is not installed. Please install pip3 before running this script."
    exit 1
fi

if ! command -v git &> /dev/null; then
    echo "git is not installed. Please install git before running this script."
    exit 1
fi

TEMP_DIR=$(mktemp -d)
cd $TEMP_DIR

echo "Downloading Nixopus..."
git clone https://github.com/raghavyuva/nixopus.git
cd nixopus/installer

echo "Setting up Python virtual environment..."
python3 -m venv venv
source venv/bin/activate

echo "Installing dependencies..."
pip install --upgrade pip
pip install -r requirements.txt

echo "Starting installation..."
python3 install.py "$@"

echo "Caddy configuration Setup..."
rm -rf /etc/nixopus/caddy
mkdir -p /etc/nixopus/caddy
echo '{
	admin 0.0.0.0:2019
	log {
		format json
		level INFO
	}
}' > /etc/nixopus/caddy/Caddyfile
chmod 644 /etc/nixopus/caddy/Caddyfile

deactivate
cd $TEMP_DIR/..
rm -rf $TEMP_DIR

echo "Installation completed!"