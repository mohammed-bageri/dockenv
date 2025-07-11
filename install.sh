#!/bin/bash

# dockenv installer script
# Usage: curl -s https://raw.githubusercontent.com/mohammed-bageri/dockenv/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/mohammed-bageri/dockenv"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="dockenv"
VERSION="latest"

# Platform detection
OS=""
ARCH=""

detect_platform() {
    echo -e "${BLUE}🔍 Detecting platform...${NC}"
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)     OS="linux";;
        Darwin*)    OS="darwin";;
        CYGWIN*)    OS="windows";;
        MINGW*)     OS="windows";;
        MSYS*)      OS="windows";;
        *)          echo -e "${RED}❌ Unsupported operating system: $(uname -s)${NC}"; exit 1;;
    esac
    
    # Detect Architecture
    case "$(uname -m)" in
        x86_64)     ARCH="amd64";;
        amd64)      ARCH="amd64";;
        arm64)      ARCH="arm64";;
        aarch64)    ARCH="arm64";;
        armv7l)     ARCH="armv7";;
        *)          echo -e "${RED}❌ Unsupported architecture: $(uname -m)${NC}"; exit 1;;
    esac
    
    echo -e "${GREEN}   Platform: ${OS}/${ARCH}${NC}"
}

check_dependencies() {
    echo -e "${BLUE}🔍 Checking dependencies...${NC}"
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        echo -e "${YELLOW}⚠️  Docker not found. Installing Docker...${NC}"
        install_docker
    else
        echo -e "${GREEN}✅ Docker found: $(docker --version)${NC}"
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${YELLOW}⚠️  Docker Compose not found. Installing Docker Compose...${NC}"
        install_docker_compose
    else
        if command -v docker-compose &> /dev/null; then
            echo -e "${GREEN}✅ Docker Compose found: $(docker-compose --version)${NC}"
        else
            echo -e "${GREEN}✅ Docker Compose found: $(docker compose version)${NC}"
        fi
    fi
    
    # Check if curl is installed
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}❌ curl is required but not installed.${NC}"
        echo -e "${YELLOW}   Please install curl and try again.${NC}"
        exit 1
    fi
}

install_docker() {
    case "$OS" in
        linux)
            echo -e "${BLUE}📦 Installing Docker on Linux...${NC}"
            if command -v apt-get &> /dev/null; then
                # Ubuntu/Debian
                curl -fsSL https://get.docker.com -o get-docker.sh
                sudo sh get-docker.sh
                sudo usermod -aG docker $USER
                rm get-docker.sh
            elif command -v yum &> /dev/null; then
                # CentOS/RHEL
                curl -fsSL https://get.docker.com -o get-docker.sh
                sudo sh get-docker.sh
                sudo usermod -aG docker $USER
                rm get-docker.sh
            else
                echo -e "${YELLOW}⚠️  Automatic Docker installation not supported for this Linux distribution.${NC}"
                echo -e "${YELLOW}   Please install Docker manually: https://docs.docker.com/get-docker/${NC}"
            fi
            ;;
        darwin)
            echo -e "${YELLOW}⚠️  Please install Docker Desktop for Mac: https://docs.docker.com/docker-for-mac/install/${NC}"
            ;;
        windows)
            echo -e "${YELLOW}⚠️  Please install Docker Desktop for Windows: https://docs.docker.com/docker-for-windows/install/${NC}"
            ;;
    esac
}

install_docker_compose() {
    echo -e "${BLUE}📦 Installing Docker Compose...${NC}"
    
    # Try to use Docker Compose plugin (newer method)
    if docker compose version &> /dev/null; then
        echo -e "${GREEN}✅ Docker Compose plugin already available${NC}"
        return
    fi
    
    # Install standalone docker-compose
    case "$OS" in
        linux)
            COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep 'tag_name' | cut -d\" -f4)
            sudo curl -L "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
            sudo chmod +x /usr/local/bin/docker-compose
            ;;
        darwin)
            echo -e "${YELLOW}⚠️  Docker Compose should be included with Docker Desktop for Mac${NC}"
            ;;
        windows)
            echo -e "${YELLOW}⚠️  Docker Compose should be included with Docker Desktop for Windows${NC}"
            ;;
    esac
}

download_binary() {
    echo -e "${BLUE}⬇️  Downloading dockenv...${NC}"
    
    # GitHub releases URL
    DOWNLOAD_URL="${REPO_URL}/releases/latest/download/dockenv-${OS}-${ARCH}"
    
    if [ "$VERSION" != "latest" ]; then
        DOWNLOAD_URL="${REPO_URL}/releases/download/v${VERSION}/dockenv-${OS}-${ARCH}"
    fi
    
    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    # Download binary
    echo -e "${BLUE}   Downloading from: ${DOWNLOAD_URL}${NC}"
    
    # Try to download the actual binary
    if curl -L "$DOWNLOAD_URL" -o "$BINARY_NAME" --fail --silent --show-error; then
        chmod +x "$BINARY_NAME"
        echo -e "${GREEN}✅ Downloaded successfully${NC}"
    else
        echo -e "${YELLOW}⚠️  Release binary not available, building from source...${NC}"
        
        # Fallback: try to build from source if Go is available
        if command -v go &> /dev/null; then
            echo -e "${BLUE}   Cloning repository and building...${NC}"
            git clone "$REPO_URL.git" dockenv-src
            cd dockenv-src
            go build -o "../$BINARY_NAME" .
            cd ..
            chmod +x "$BINARY_NAME"
            echo -e "${GREEN}✅ Built from source successfully${NC}"
        else
            echo -e "${RED}❌ Failed to download binary and Go is not available for building from source${NC}"
            echo -e "${YELLOW}   Please install Go or download a release manually from: ${REPO_URL}/releases${NC}"
            exit 1
        fi
    fi
}

install_binary() {
    echo -e "${BLUE}📦 Installing dockenv...${NC}"
    
    # Check if we have write permission to install directory
    if [ ! -w "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}   Using sudo for installation (requires root privileges)${NC}"
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        mv "$BINARY_NAME" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Verify installation
    if command -v dockenv &> /dev/null; then
        echo -e "${GREEN}✅ dockenv installed successfully!${NC}"
        echo -e "${GREEN}   Version: $(dockenv --version 2>/dev/null || echo 'dockenv development build')${NC}"
    else
        echo -e "${RED}❌ Installation failed${NC}"
        exit 1
    fi
}

cleanup() {
    echo -e "${BLUE}🧹 Cleaning up...${NC}"
    cd /
    rm -rf "$TMP_DIR"
}

show_next_steps() {
    echo
    echo -e "${GREEN}🎉 Installation completed successfully!${NC}"
    echo
    echo -e "${BLUE}📋 Next steps:${NC}"
    echo -e "${YELLOW}   1. Start a new terminal session or run: source ~/.bashrc${NC}"
    echo -e "${YELLOW}   2. Navigate to your project directory${NC}"
    echo -e "${YELLOW}   3. Initialize dockenv: ${GREEN}dockenv init${NC}"
    echo -e "${YELLOW}   4. Start services: ${GREEN}dockenv up${NC}"
    echo
    echo -e "${BLUE}📚 Common commands:${NC}"
    echo -e "${YELLOW}   dockenv init               # Interactive setup${NC}"
    echo -e "${YELLOW}   dockenv init --profile laravel  # Quick Laravel setup${NC}"
    echo -e "${YELLOW}   dockenv up                 # Start services${NC}"
    echo -e "${YELLOW}   dockenv down               # Stop services${NC}"
    echo -e "${YELLOW}   dockenv status             # Check service status${NC}"
    echo -e "${YELLOW}   dockenv add mysql          # Add a service${NC}"
    echo -e "${YELLOW}   dockenv autostart enable   # Auto-start on boot${NC}"
    echo
    echo -e "${BLUE}📖 Documentation: ${REPO_URL}#readme${NC}"
    echo -e "${BLUE}🐛 Issues: ${REPO_URL}/issues${NC}"
}

# Main installation flow
main() {
    echo -e "${BLUE}🐳 dockenv installer${NC}"
    echo -e "${BLUE}   Local development environments made easy${NC}"
    echo
    
    detect_platform
    check_dependencies
    download_binary
    install_binary
    cleanup
    show_next_steps
}

# Handle script interruption
trap cleanup EXIT

# Check if running with sudo (not recommended)
if [ "$EUID" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  Running as root is not recommended.${NC}"
    echo -e "${YELLOW}   Please run as a regular user.${NC}"
    echo -e "${YELLOW}   The installer will request sudo when needed.${NC}"
    exit 1
fi

# Run main function
main "$@"
