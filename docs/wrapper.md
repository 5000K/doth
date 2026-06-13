---
title: The Wrapper Script
description: A self-contained way to run doth without installing Go or doth globally.
author: 5000K
---

The wrapper is a small shell script that ships with a doth project. It manages the project's own `doth` installation. The `doth` binary and it's dependencies live inside the project's `.doth/` directory. Nothing is installed system wide.

## What it solves

Setting up a new machine needs `doth` to deploy the project. Installing it by hand means installing Go first, then `doth`, then keeping both up to date. The wrapper skips all of that.

## Who it is for

The wrapper is for anyone who wants to set up a new computer without installing Go or `doth` globally before doing the setup. This way, doth is self-contained and your working environment is not influenced by you using doth in any way.

## The workflow

Setting up a new computer has two steps.

1. Clone the project.
2. Run `./doth.sh deploy` and/or `./doth.sh install {package manager flags/configs}`.

```sh
git clone https://github.com/you/dotfiles.git ~/dotfiles
cd ~/dotfiles
./doth.sh deploy # deploy config files
./doth.sh install --pacman # install deps
```

The wrapper downloads Go and `doth` on the first run. Both stay inside the project folder. The wrapper updates both in place when new versions appear. You can lock the used version using the `doth lock` command, and unlock it (=> "use newest version") using the `doth unlock` command. The wrapper itself does not touch anything (files or environment) outside the project.

See [Setup](./setup.md) for the lock workflow and how to create the wrapper.
