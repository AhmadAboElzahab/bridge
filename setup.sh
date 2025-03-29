#!/bin/bash
set -e

echo "Setting up Go development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    exit 1
fi

# Install pre-commit
echo "Installing pre-commit..."
if ! command -v pip3 &> /dev/null; then
    echo "Installing pip first..."
    if command -v apt-get &> /dev/null; then
        sudo apt-get update && sudo apt-get install -y python3-pip
    elif command -v brew &> /dev/null; then
        brew install python3
    else
        echo "Error: Could not install pip. Please install pip3 manually."
        exit 1
    fi
fi

pip3 install pre-commit

# Install Go development tools
echo "Installing Go development tools..."
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install pre-commit hooks
echo "Installing pre-commit hooks..."
pre-commit install

echo "✅ Setup complete! Your development environment is ready."
echo "Every commit will now automatically format your Go code."