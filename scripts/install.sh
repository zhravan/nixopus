#!/bin/bash

if [ -z "$1" ]; then
    echo "Please specify the domain name or IP address"
    exit 1
else
    DOMAIN="$1"
fi

echo "Installing the application"

# We should install the dependencies here such as docker

echo "Setting up environment variables"

cp .env.sample .env

perl -pi -e "s|DOMAIN|${DOMAIN}|g" .env

if docker compose up --build -d; then
    echo "Started successfully"
else
    echo "Failed to start"
    exit 1
fi

echo "Application installed successfully"