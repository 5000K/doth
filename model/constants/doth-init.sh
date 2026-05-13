#!/bin/bash

# This is a one-time bootstrapping script to initialize a fully self-contained doth project, including the current wrapper.

ROOT_DIR=$(dirname "$(realpath "$0")")

GO_VERSION=1.26.2
TARGET_FOLDER=${ROOT_DIR}/.doth/go
GOPATH_FOLDER=${ROOT_DIR}/.doth/gopath

GO_COMMAND=${TARGET_FOLDER}/bin/go


# set paths for this script
export GOPATH=${GOPATH_FOLDER}
export GOBIN=${GOPATH_FOLDER}/bin
export PATH=${TARGET_FOLDER}/bin:${GOBIN}:$PATH

detect_platform() {
    local os arch
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)
    
    case "$os" in
        linux)  GO_OS="linux" ;;
        darwin) GO_OS="darwin" ;;
        *) echo "Unsupported OS: $os"; exit 1 ;;
    esac
    
    case "$arch" in
        x86_64)         GO_ARCH="amd64" ;;
        aarch64|arm64)  GO_ARCH="arm64" ;;
        armv6l|armv7l)  GO_ARCH="armv6l" ;;
        *) echo "Unsupported architecture: $arch"; exit 1 ;;
    esac
}

setup_go() {
    # check if Go is already installed
    if command -v ${GO_COMMAND} &> /dev/null; then
        return
    fi
    
    detect_platform
    echo "Installing Go ${GO_VERSION} (${GO_OS}-${GO_ARCH})..."
    
    # Install Go
    mkdir -p ${TARGET_FOLDER}
    mkdir -p ${GOBIN}
    curl -Ls -o ${ROOT_DIR}/.doth/go.tar.gz https://go.dev/dl/go${GO_VERSION}.${GO_OS}-${GO_ARCH}.tar.gz
    tar -C ${TARGET_FOLDER} --strip-components=1 -xzf ${ROOT_DIR}/.doth/go.tar.gz
    rm ${ROOT_DIR}/.doth/go.tar.gz
}

setup_doth() {
    echo "Installing the current version of doth..."
    NEW_VERSION=$(curl -Ls https://github.com/5000K/doth/releases/latest/download/version.txt)
    ${GO_COMMAND} install github.com/5000K/doth@$NEW_VERSION
}
# is doth available?
if ! command -v doth &> /dev/null; then
    setup_go
    setup_doth
fi

# initialize a fresh doth project
doth init --modules ./modules --verbose

echo "doth is initialized and ready to use!"
echo "instead of installing doth globally, use ./doth.sh to run the current version of doth."
