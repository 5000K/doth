<p align="center">
  <img src="./badge.png" alt="doth" width="360">
</p>

---

<br/>

[![CI](https://github.com/5000K/doth/actions/workflows/ci.yml/badge.svg)](https://github.com/5000K/doth/actions/workflows/ci.yml) [![Release](https://github.com/5000K/doth/actions/workflows/release.yml/badge.svg)](https://github.com/5000K/doth/actions/workflows/release.yml)

**doth** is a go-based .config manager for linux.

> doth is a pre-release version. we actively use it at 5000K and it proved stable. Since two of us actively based their arch setup on doth, we are very careful not to introduce breaking changes at this point in time. But it still is pre-release software and should be treated like it.

## The idea

**doth** can simply copy or symlink your files from your git repo. Or it can render them as templates.

Show the battery state in the top bar on your laptop specifically? Deploy with different themes by simply passing another one into doth? Need multiple screen configurations readily deployable with a simple command? Doth is the tool for you. It's a bit more complex, but you are in control in return.

## Modules

A module is a subfolder with a `module.yaml` and the files it deploys. Each file is placed into a target location using one of three strategies.

- **copy** writes a copy of the file.
- **link** creates a symbolic link to the file in the repository.
- **render** runs the file as a Go template and writes the result.

See the [modules guide](https://doth.5000K.org/modules) for the full reference.

## Templates

The render strategy uses the Go [text/template](https://pkg.go.dev/text/template) system. The values come from configuration files passed to `doth deploy`. See the [modules guide](https://doth.5000K.org/modules) for examples.

## Dependencies

Dependencies are declared in `doth.yaml` and `module.yaml`. `doth install` reads the declarations and installs the packages. Built-in flags cover common package managers such as `apt`, `dnf`, and `pacman`. Custom sources can be defined in configuration files. See the [dependencies guide](https://doth.5000K.org/dependencies) for the full reference.

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
