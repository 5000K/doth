package model

const DothFileLocation = "./doth.yaml"
const DothFileTemplate = `# This is the doth.yaml file. It defines the configuration for your doth project.

# The directory where your modules are located. Relative to the current working directory.
moduleDir: "./modules"

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

const GitignoreFileLocation = "./.gitignore"
const GitignoreFileTemplate = `# This is the .gitignore file. It defines which files and directories should be ignored by git.


# Ignore the .doth directory, it contains local state and should not be committed to version control.
.doth/
`

const LocalStateDir = "./.doth"
