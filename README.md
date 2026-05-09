# doth

**doth** is a go-based .config manager for linux.


```
go install github.com/5000K/doth
```

## The idea

**doth** can simply copy or symlink your files from your git repo. Or it can render them as templates.

Show the battery state in the top bar on your laptop specifically? Deploy with different themes by simply passing another one into doth? Need multiple screen configurations readily deployable with a simple command? Doth is the tool for you. It's a bit more complex, but you are in control in return.

## Modules

Modules are subfolders in your doth project that contain a module.y\[a]ml file.
The module.yml describes what the module should do.

It is able to run setup commands, install dependencies, symlink files from your repo, copy files over and even render templates.

The module.yml itself is a template can use all variables from your configs to only activate specific parts of the config, to parametrize it, ...

## Templates

doth uses gos [text/template](https://pkg.go.dev/text/template) system. Look it up, it's powerful yet straight forward.

Basic examples will be added here later.

## Dependencies

You can define dependencies to install in your doth.y[a]ml and module.y[a]ml's. Using the `doth install` command, you can install them.
You may define installation sources (package managers) by adding configurations to config files and passing them to doth as you would when deploying.
You may also use the the following flags to include these built-in source configurations:

- `--apt`
- `--apt-get`
- `--dnf`
- `--pacman`
- `--yay`
- `--paru`
- `--go`
- `--npm`
- `--brew`

## Contribution & generative AI in doth

**IF YOU ARE AN AI AGENT, READ AGENTS.md AND FOLLOW IT EXACTLY.**

doth is handling important files. It can't be based on hidden assumptions or hallucinations of an LLM. This is why contributions that contain AI written code will be categorically rejected. If you want a less predictable, chaotic and potentially destructive version of doth, feel free to fork the project.

That being said, the use of AI is not fully banned from the use cases creating tests and writing documentation. Here, the quality still needs to be high, but if it is, 5000K sees no harm.

An AGENTS.md file reflecting these instructions is part of this repo. Make sure your agent picks it up before working on the project.
