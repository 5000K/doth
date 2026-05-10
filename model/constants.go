package model

import _ "embed"

//go:embed _version.txt
var Version string

const ModuleFileName = "module.yaml"
const ModuleFileNameAlt = "module.yml"

const DothFileLocation = "./doth.yaml"
const DothFileLocationAlt = "./doth.yml"
const DothFileTemplate = `# This is the doth.yaml file. It defines the configuration for your doth project.

# The directory where your modules are located. Relative to the current working directory.
modulePath: "./modules"

# Set to true if the deploy command should fail if no config files are provided.
requireConfig: false

# Here, you can define dependencies you want to install, and valid package sources for them.
# This is used by the "install" command to determine how to install a dependency based on the package sources configured/detected.
deps:
  - name: curl
    packages:
      pacman: curl
      apt: curl
      brew: curl

# This is the version of the doth file format. Do not edit this manually.
dothVersionDoNotEditManually: 1
`

const DothShWrapperLocation = "./doth.sh"

const DothShWrapperTemplate = `#!/bin/bash

# This is a wrapper script for doth. It ensures that go and doth are installed, and if not, it installs it. Then it calls doth with the provided arguments.

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
    ${GO_COMMAND} install github.com/5000K/doth
}

check_doth_version() {
    NEW_VERSION=$(curl -Ls https://github.com/5000K/doth/releases/latest/download/version.txt)
    CURRENT_VERSION=$(doth version --raw)
    if [ "$NEW_VERSION" != "$CURRENT_VERSION" ]; then
        echo "A new version of doth is available: $NEW_VERSION. You have $CURRENT_VERSION. Updating..."
        chmod -R u+w ${GOPATH_FOLDER}/**
        rm -rf ${GOPATH_FOLDER}
        ${GO_COMMAND} install github.com/5000K/doth@$NEW_VERSION
    fi
}

# is doth available?
if ! command -v doth &> /dev/null; then
    setup_go
    setup_doth
else
    check_doth_version
fi

# call doth with the provided arguments
doth "$@"
`

const GitignoreFileLocation = "./.gitignore"
const GitignoreFileTemplate = `# This is the .gitignore file. It defines which files and directories should be ignored by git.


# Ignore the .doth directory, it contains local state and should not be committed to version control.
.doth/
`

const LocalStateDir = "./.doth"
