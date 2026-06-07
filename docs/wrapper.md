---
title: The Wrapper Script
description: A self-contained way to run doth without installing Go or doth globally.
author: 5000K
---

# The Wrapper Script

The wrapper is a small shell script that ships with a doth project. It manages the project's own `doth` installation. Go and the `doth` binary live inside the project's `.doth/` directory. Nothing is installed system wide.

## What it solves

Setting up a new machine needs `doth` to deploy the project. Installing it by hand means installing Go first, then `doth`, then keeping both up to date. The wrapper skips all of that.

## Who it is for

The wrapper is for anyone who wants to set up a new computer without installing Go or `doth` globally. It is for anyone who wants the same `doth` version on every machine. It is for anyone who wants a fully self contained setup.

## The workflow

Setting up a new computer is two steps.

1. Clone the project.
2. Run `./doth.sh deploy` or `./doth.sh install --pacman`.

```sh
git clone https://github.com/you/dotfiles.git ~/dotfiles
cd ~/dotfiles
./doth.sh deploy
```

```sh
./doth.sh install --pacman
```

The wrapper downloads Go and `doth` on the first run. Both stay inside the project folder. The wrapper updates both in place when new versions appear. A `doth.lock` file pins the version. The wrapper does not touch anything outside the project.

See [[Setup]] for the lock workflow and how to create the wrapper.
