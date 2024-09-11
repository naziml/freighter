#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Update package list and install prerequisites
echo "Updating package list..."
sudo apt update -y

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

   # Install Go from the official source (replace the version as needed)
    GO_VERSION="1.21.1"
    wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz

    # Remove any previous Go installation
    sudo rm -rf /usr/local/go

    # Install Go
    sudo tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz

    # Add Go to the PATH
    if ! grep -q 'export PATH=$PATH:/usr/local/go/bin' ~/.bashrc; then
        echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    fi

    # Source the updated bashrc to apply changes immediately
    source ~/.bashrc

    # Verify Go installation
    go version

    # Clean up the downloaded archive
    rm go$GO_VERSION.linux-amd64.tar.gz
fi

# Verify installations
echo "Verifying installations..."
git --version
go version

# Install build tools
sudo apt install gcc make
sudo apt install libsqlite3-dev
export CGO_ENABLED=1

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
