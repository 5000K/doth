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

const DothShWrapperTemplate = `
#!/bin/bash

# This is a wrapper script for doth. It ensures that go and doth are installed, and if not, it installs it. Then it calls doth with the provided arguments.

ROOT_DIR=$(dirname "$(realpath "$0")")

GO_VERSION=1.26.2
TARGET_FOLDER=${ROOT_DIR}/.doth/go


# set paths for this script
export PATH=${TARGET_FOLDER}/bin:$PATH
export GOPATH=${TARGET_FOLDER}
export GOBIN=${TARGET_FOLDER}/bin

setup_go() {
    # check if Go is already installed
    if command -v go &> /dev/null; then
        echo "Go is already installed. Skipping installation."
        return
    fi
    
    # Install Go
    curl -L -o go.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    mkdir -p ${TARGET_FOLDER}
    tar -C ${TARGET_FOLDER} --strip-components=1 -xzf go.tar.gz
    rm go.tar.gz
}

setup_doth() {
    go install github.com/5000K/doth
}

check_doth_version() {
    NEW_VERSION=$(curl -L https://github.com/5000K/doth/releases/latest/download/version.txt)
    CURRENT_VERSION=$(doth --version --raw)
    if [ "$NEW_VERSION" != "$CURRENT_VERSION" ]; then
        echo "A new version of doth is available: $NEW_VERSION. You have $CURRENT_VERSION. Updating..."
        go install github.com/5000K/doth@$NEW_VERSION
    else
        echo "You have the latest version of doth: $CURRENT_VERSION."
    fi
}

# is doth available?
if ! command -v doth &> /dev/null; then
    echo "Doth is not installed. Installing Go and Doth..."
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
