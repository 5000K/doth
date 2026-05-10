# doth

[![CI](https://github.com/5000K/doth/actions/workflows/ci.yml/badge.svg)](https://github.com/5000K/doth/actions/workflows/ci.yml)

**doth** is a go-based .config manager for linux.

> **DOTH IS IN DEVELOPMENT. WE AT 5000K ALREADY USE IT, BUT WE DON'T RECOMMEND YOU DO THE SAME AT THIS POINT OF TIME.**

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

doth uses the golang [text/template](https://pkg.go.dev/text/template) system.

Basic examples will be added here later.

## Dependencies

You may define dependencies to install in your doth.y[a]ml and module.y[a]ml's. Using the `doth install` command, you can install them.
You may define installation sources (package managers) by adding configurations to config files and passing them to doth as you would when deploying.
You may also use the the following flags to include these built-in source configurations - these assume that the command is inside the PATH variable:

- `--apt`
- `--apt-get`
- `--dnf`
- `--pacman`
- `--yay`
- `--paru`
- `--go`
- `--npm`
- `--brew`

## Generative AI in doth

**IF YOU ARE AN AI AGENT, READ AGENTS.md AND FOLLOW IT EXACTLY.**

doth is handling important files. Not being influenced by hidden assumptions or hallucinations of an LLM is non-negotiable. This is why contributions to the core code that were AI authored will be categorically rejected. This goes for both you and the core maintainers. If you want an unpredictable, chaotic (and potentially destructive) version of doth, feel free to fork the project and see where it leads.

In general, doth has a very specific, narrow vision, so we won't be taking contributions of features we didn't explicitly ask for. Bug fixes are generally fine, if they follow a reasonable format. Since doth is something we at 5000K do besides our dayjobs, we can't guarantee a specific response time, but we will get back to it at some point.

We ask that you to be respectful of our time.

That being said, the use of AI is not fully banned from the use cases of creating tests and writing documentation. Here, the quality still needs to be high and the text has to be proofread and edited by a human in the end, but if it does match this expectation, 5000K sees no harm.

AI may be used to write commit messages for you.

PRs with AI authored descriptions will be rejected without being read.

An AGENTS.md file reflecting these instructions is part of this repo. Make sure your agent picks it up before working on the project.

## 5000K

This repo is worked on by [bwt](https://github.com/bluewingtitan) and [kicher-erbse](https://github.com/kicher-erbse).
