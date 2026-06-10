---
title: Overview
description: An introduction to doth, what it solves, and who it is for.
author: 5000K
---

# Overview

doth is a tool that manages configuration files on Linux. It treats your dotfiles as a modular repository. The repository is version controlled. The same repository can be deployed onto many machines with different settings.

## The problem

Configuration files are scattered across the home directory. They are tedious to back up. They are tedious to set up on a new machine. The usual approach is a single git repository and a few symlinks. This works for one machine. It breaks down when you have several machines with different needs.

## What doth does

doth reads a project that contains a `doth.yaml` and a folder of modules. Each module declares where its files go. `doth deploy` then places the files.

Three placement strategies are available. The strategy is chosen per file.

- **copy** writes a copy of the file into the target location.
- **link** creates a symbolic link from the target location into the repository.
- **render** runs the file through a template engine and writes the result.

The template engine is the Go standard library `text/template`. The values come from one or more configuration files. The user passes those files to `doth deploy`. This lets one repository serve many machines.

doth can also install the packages your configuration depends on. See [Dependencies](./dependencies.md) for details.

## Who it is for

doth is for people who want to manage their dotfiles flexibly. It is for people who run more than one machine. It is for people who want different flavours of the same configuration files ready to deploy - on different machines, or on a single one. It is for people who want to switch themes or layouts by passing a different configuration file.

doth is not for people with a handful of static dotfiles on a single machine. A plain git repository with symlinks is enough for that.

## The shape of a project

A doth project is a directory. It contains a `doth.yaml` and a `modules` directory. The `doth.yaml` holds top level settings and top level dependencies. Each subfolder of `modules/` is a module. The module folder contains a `module.yaml` and the files the module deploys.

`doth.yaml` and the module files are the definition files of a doth project. They define how to deploy what.

```
my-dotfiles/
├── doth.yaml
├── modules/
│   ├── bash/
│   │   ├── module.yaml
│   │   └── .bashrc
│   └── nvim/
│       ├── module.yaml
│       └── init.lua
```

Configuration files are kept seperately to the definition files. They hold the values that vary between unique configurations. They are passed to `doth deploy` with the `--config` flag.

```
my-dotfiles/
├── configs/
│   ├── laptop.yaml
│   └── server.yaml
```

## Where to go next

- [Setup](./setup.md) covers installing `doth` and creating a new project.
- [The Wrapper Script](./wrapper.md) explains the self contained workflow for deploying a project on a new machine.
- [Modules](./modules.md) covers writing modules.
- [Dependencies](./dependencies.md) covers the package install system.
- [YML Reference](./yml.md) covers the full format of the YAML files.
