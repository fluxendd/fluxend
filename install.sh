#!/bin/bash

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command_exists apt-get; then
            echo "ubuntu"
        elif command_exists yum; then
            echo "centos"
        elif command_exists dnf; then
            echo "fedora"
        else
            echo "linux"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    else
        echo "unknown"
    fi
}

# Function to install Docker
install_docker() {
    local os=$(detect_os)
    print_status "Installing Docker..."

    case $os in
        "ubuntu")
            sudo apt-get update
            sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release
            curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
            echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
            sudo apt-get update
            sudo apt-get install -y docker-ce docker-ce-cli containerd.io
            sudo systemctl start docker
            sudo systemctl enable docker
            sudo usermod -aG docker $USER
            ;;
        "centos"|"fedora")
            sudo yum install -y yum-utils
            sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
            sudo yum install -y docker-ce docker-ce-cli containerd.io
            sudo systemctl start docker
            sudo systemctl enable docker
            sudo usermod -aG docker $USER
            ;;
        "macos")
            if command_exists brew; then
                brew install --cask docker
                print_warning "Please start Docker Desktop manually after installation"
            else
                print_error "Homebrew not found. Please install Docker Desktop manually from https://docker.com/products/docker-desktop"
                exit 1
            fi
            ;;
        *)
            print_error "Unsupported OS for automatic Docker installation"
            exit 1
            ;;
    esac
}

# Function to install Docker Compose
install_docker_compose() {
    print_status "Installing Docker Compose..."

    # Try to install via package manager first
    local os=$(detect_os)
    case $os in
        "ubuntu")
            sudo apt-get install -y docker-compose-plugin
            ;;
        "centos"|"fedora")
            sudo yum install -y docker-compose-plugin
            ;;
        "macos")
            # Docker Compose comes with Docker Desktop on macOS
            print_status "Docker Compose comes with Docker Desktop on macOS"
            ;;
        *)
            # Fallback to manual installation
            sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
            sudo chmod +x /usr/local/bin/docker-compose
            ;;
    esac
}

# Function to install make
install_make() {
    local os=$(detect_os)
    print_status "Installing make..."

    case $os in
        "ubuntu")
            sudo apt-get install -y build-essential
            ;;
        "centos"|"fedora")
            sudo yum groupinstall -y "Development Tools"
            ;;
        "macos")
            if command_exists brew; then
                brew install make
            else
                xcode-select --install
            fi
            ;;
        *)
            print_error "Unsupported OS for automatic make installation"
            exit 1
            ;;
    esac
}

# Function to check if all required tools are installed
check_requirements() {
    local missing=()

    if ! command_exists docker; then
        missing+=("docker")
    fi

    if ! command_exists docker-compose && ! docker compose version >/dev/null 2>&1; then
        missing+=("docker-compose")
    fi

    if ! command_exists make; then
        missing+=("make")
    fi

    echo "${missing[@]}"
}

# Main script starts here
print_status "Starting Fluxend setup..."

# Step 1: Download the archive
print_status "Downloading fluxend.zip..."
curl -L -o fluxend.zip https://github.com/fluxendd/fluxend/archive/refs/heads/main.zip

# Step 2: Extract the archive
print_status "Extracting archive..."
if command_exists unzip; then
    unzip -q fluxend.zip
elif command_exists tar; then
    # Some systems might have tar but not unzip
    print_warning "unzip not found, trying with tar..."
    tar -xf fluxend.zip
else
    print_error "Neither unzip nor tar found. Cannot extract archive."
    exit 1
fi

# Step 3: Change to directory
print_status "Changing to fluxend directory..."
cd fluxend-main

# Step 4: Check and install missing requirements
print_status "Checking requirements..."
missing=($(check_requirements))

if [ ${#missing[@]} -gt 0 ]; then
    print_warning "Missing requirements: ${missing[*]}"

    for tool in "${missing[@]}"; do
        case $tool in
            "docker")
                install_docker
                ;;
            "docker-compose")
                install_docker_compose
                ;;
            "make")
                install_make
                ;;
        esac
    done

    print_status "Waiting for services to start..."
    sleep 5
fi

# Step 5: Final verification
print_status "Verifying installation..."
missing=($(check_requirements))

if [ ${#missing[@]} -gt 0 ]; then
    print_error "Installation failed. Still missing: ${missing[*]}"
    print_error "Please install the missing tools manually and run this script again."
    exit 1
fi

print_status "All requirements satisfied!"

# Special handling for Docker group membership
if groups $USER | grep -q docker; then
    print_status "User is in docker group"
else
    print_warning "User added to docker group. You may need to log out and back in, or run 'newgrp docker'"
fi

# Step 6: Run make setup
print_status "Running make setup..."
make setup

print_status "Setup complete!"
print_status "Fluxend has been successfully set up in $(pwd)"
