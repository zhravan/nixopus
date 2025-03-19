#!/bin/bash

if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root" >&2
    echo "Please run with sudo or as root user" >&2
    exit 1
fi

KEY_TYPE="rsa"
KEY_NAME="nixopus"
APP_DIR="/opt/nixopus"
SSH_DIR="$APP_DIR/.ssh"

while [ $# -gt 0 ]; do
  case "$1" in
    --app-dir=*)
      APP_DIR="${1#*=}"
      SSH_DIR="$APP_DIR/.ssh"
      shift
      ;;
    --key-name=*)
      KEY_NAME="${1#*=}"
      shift
      ;;
    --key-type=*)
      KEY_TYPE="${1#*=}"
      shift
      ;;
    *)
      echo "Unknown parameter: $1" >&2
      exit 1
      ;;
  esac
done

PRIVATE_KEY_FILE="$SSH_DIR/${KEY_NAME}_private_key"
PUBLIC_KEY_FILE="$SSH_DIR/${KEY_NAME}_private_key.pub"

if [ ! -d "$APP_DIR" ]; then
  echo "Creating application directory $APP_DIR..."
  mkdir -p "$APP_DIR"
  chmod 755 "$APP_DIR"
fi

if [ ! -d "$SSH_DIR" ]; then
  echo "Creating SSH directory $SSH_DIR..."
  mkdir -p "$SSH_DIR"
  chmod 700 "$SSH_DIR"
fi


echo "Generating $KEY_TYPE SSH key pair for application in $APP_DIR..."
ssh-keygen -t "$KEY_TYPE" -f "$PRIVATE_KEY_FILE" -N ""


chmod 600 "$PRIVATE_KEY_FILE"
chmod 644 "$PUBLIC_KEY_FILE"

echo -e "\nKeys have been saved to:"
echo "Private key: $PRIVATE_KEY_FILE"
echo "Public key: $PUBLIC_KEY_FILE"


if [ -n "$SUDO_USER" ]; then
  echo "Setting ownership to $SUDO_USER..."
  chown -R "$SUDO_USER" "$APP_DIR"
else
  echo "Running as root user, ownership remains with root"
fi