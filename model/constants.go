package model

import _ "embed"

//go:embed _version.txt
var Version string

const ModuleFileName = "module.yaml"
const ModuleFileNameAlt = "module.yml"

const DothFileLocation = "./doth.yaml"
const DothFileLocationAlt = "./doth.yml"

//go:embed constants/doth.yaml
var DothFileTemplate string

const DothShWrapperLocation = "./doth.sh"

//go:embed constants/doth.sh
var DothShWrapper string

const GitignoreFileLocation = "./.gitignore"

//go:embed constants/.gitignore
var GitignoreFileTemplate string

const LocalStateDir = "./.doth"

const DothLockFileLocation = "./doth.lock"
