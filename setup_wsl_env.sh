#!/bin/bash
# AAR Development Setup for WSL (Windows Subsystem for Linux)

echo "🐧 Setting up AAR Dev Environment in WSL..."

# 1. Update and Install GCC (needed for compiling some Go networking libs)
sudo apt-get update
sudo apt-get install -y gcc libc6-dev git

# 2. Check for Go
if ! command -v go &> /dev/null
then
    echo "⬇️ Go not found. Installing Go 1.21..."
    wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    echo "✅ Go installed."
else
    echo "✅ Go is already installed."
fi

# 3. Navigate to Backend
cd backend

# 4. Tidy Modules
/usr/local/go/bin/go mod tidy

echo "-----------------------------------------------------"
echo "🚀 READY! To start the AAR Brain in Linux Mode:"
echo "   cd backend"
echo "   sudo /usr/local/go/bin/go run main.go"
echo "-----------------------------------------------------"
