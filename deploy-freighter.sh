#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Update package list and install prerequisites
echo "Updating package list..."
sudo apt update

# Install Git if it's not already installed
if command_exists git; then
    echo "Git is already installed."
else
    echo "Installing Git..."
    sudo apt install -y git
fi

# Install Go if it's not already installed
if command_exists go; then
    echo "Go is already installed."
else
    echo "Installing Go..."

    # Download the latest version of Go
    GO_VERSION=$(curl -sL https://golang.org/dl/ | grep -oP 'go[0-9\.]+' | head -n 1)
    GO_TAR_FILE="go${GO_VERSION}.linux-amd64.tar.gz"
    GO_DOWNLOAD_URL="https://golang.org/dl/${GO_TAR_FILE}"

    echo "Downloading Go ${GO_VERSION}..."
    curl -LO "${GO_DOWNLOAD_URL}"

    echo "Extracting Go..."
    sudo tar -C /usr/local -xzf "${GO_TAR_FILE}"

    # Clean up
    rm "${GO_TAR_FILE}"

    # Set up Go environment variables
    echo "Setting up Go environment variables..."
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
fi

# Verify installations
echo "Verifying installations..."
git --version
go version

# Create repo directory and clone code
echo "Getting and building code..."
mkdir ~/code
cd ~/code

# Clone Go Container Registry
echo "Get container registry code"
git clone https://github.com/naziml/go-containerregistry.git
cd go-containerregistry
git checkout manifest-store

echo "Get Freighter server"
cd ~/code
git clone https://github.com/naziml/freighter.git
cd freighter/cmd/freighter_server
echo "Build Freighter server"
go build
echo "Running Freighter server"
./freigther_server &

echo "Installation script completed."
