package model

const dothFileLocation = "./doth.yaml"
const dothFileTemplate = `# This is the doth.yaml file. It defines the configuration for your doth project.

# The directory where your modules are located. Relative to the current working directory.
moduleDir: "./modules"

# This is the version of the doth file format. Do not edit this manually.
dothVersionDoNotEditManually: 1
`

const gitignoreFileLocation = "./.gitignore"
const gitignoreFileTemplate = `# This is the .gitignore file. It defines which files and directories should be ignored by git.

# Ignore the .doth directory, it contains local state and should not be committed to version control.
.doth/
`

const localStateDir = "./.doth"
