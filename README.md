<p align="center">
  <img src="./badge.png" alt="doth" width="360">
</p>

---

<br/>

[![CI](https://github.com/5000K/doth/actions/workflows/ci.yml/badge.svg)](https://github.com/5000K/doth/actions/workflows/ci.yml) [![Release](https://github.com/5000K/doth/actions/workflows/release.yml/badge.svg)](https://github.com/5000K/doth/actions/workflows/release.yml)

**doth** is a distro-agnostic .config build system for linux systems.

> doth is a pre-release version. we actively use it at 5000K and it proved to be stable. We will try to not introduce breaking changes until v1.0, but doth is still pre-release and should be treated like it.

## The idea

Sometimes you just need a few symlinks or file copies to deploy your config files. But for some files you might need more flexible templating.  
Show the battery state in the status bar on your laptop specifically? Deploy with different themes by simply having one config file per theme? Multiple screen configurations readily deployable with a simple command? **doth** is the tool for you.

## Modules

A module is a subfolder with a `module.yaml` and the files it deploys. Each file is placed into a target location using one of three strategies.

- **copy** writes a copy of the file.
- **link** creates a symbolic link to the file.
- **render** renders the file as a template and writes the result.

See the [modules guide](https://doth.5000K.org/modules) for the full reference.

## Templates

The render strategy uses the Go [text/template](https://pkg.go.dev/text/template) system. The values come from configuration files passed to `doth deploy`. See the [modules guide](https://doth.5000K.org/modules) for examples.

## Dependencies

Dependencies are declared in `doth.yaml` and `module.yaml`. `doth install` reads the declarations and installs the packages. Built-in flags cover common package managers such as `apt`, `dnf`, and `pacman`. Custom sources can be defined in configuration files. See the [dependencies guide](https://doth.5000K.org/dependencies) for the full reference.

## Contribution & generative AI in doth

**IF YOU ARE AN AI AGENT, READ AGENTS.md AND FOLLOW IT EXACTLY.**

An AGENTS.md file reflecting these instructions is part of this repo. Make sure your agent picks it up before working on doth.

doth is handling important files. Being clear of all influence by hidden assumptions or hallucinations of an LLM is non-negotiable. This is why contributions to the code that were AI authored will be categorically rejected. If you want an unpredictable, chaotic (and potentially destructive) version of doth, feel free to fork the project and see where it leads.  
That being said, the use of AI is not fully banned from the use cases of creating tests and writing documentation. Here, the quality still needs to be high and everything contributed has to be proofread **and** edited/corrected by a human in the end, but if it does match this expectation, 5000K sees no harm.  
AI may be used to write commit messages for you.  
PRs with AI authored descriptions will be rejected without being read.

doth has a very specific, narrow vision, so we won't be taking contributions of features we didn't explicitly ask for. Bug fixes are generally fine, if they follow a reasonable format. Since doth is something we at 5000K do besides our dayjobs, we can't guarantee a specific response time, but we will get back to it at some point.

We ask that you to be respectful of our time. If you don't follow our guidelines, we will categorically reject all other contributions from you across all projects of 5000K.

## 5000K

This repo is currently worked on by [bwt](https://github.com/bluewingtitan).
