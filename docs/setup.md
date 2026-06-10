---
title: Setup
description: Installing the doth binary and initializing a new project.
author: 5000K
---

# Setup

A new doth project needs the `doth` binary and a project directory. Each release of doth publishes two ways to get a working installation. The pre-built binaries cover the common platforms. The init script handles the full bootstrap on any supported system.

## Pre-built binaries

Each release attaches binaries for five platforms. The supported targets are `linux/amd64`, `linux/arm64`, `linux/arm`, `darwin/amd64`, and `darwin/arm64`. The file name is `doth-<version>-<os>-<arch>`. Place the binary on your `PATH` and make it executable.

```sh
curl -Ls -O https://github.com/5000K/doth/releases/latest/download/doth-v0.0.13-linux-amd64
chmod +x doth-v0.0.13-linux-amd64
sudo mv doth-v0.0.13-linux-amd64 /usr/local/bin/doth
```

The release also contains four extra files. They are used by the init script and the wrapper script.

| File           | Purpose                                                  |
|----------------|----------------------------------------------------------|
| `version.txt`  | The latest version number. Used to detect updates.       |
| `doth.sh`      | The wrapper script. See below.                           |
| `doth-init.sh` | The bootstrap script. See below.                         |
| `LICENSE`      | The license of the project.                              |

## The init script

The `doth-init.sh` script is published with every release. It bootstraps a self-contained doth installation in the current directory. It downloads Go into `./.doth/go/`. It installs the latest `doth` binary into `./.doth/gopath/bin/`. It then runs `doth init --modules ./modules --verbose` to create the project.

Run the script in an empty directory.

```sh
curl -Ls -O https://github.com/5000K/doth/releases/latest/download/doth-init.sh
chmod +x doth-init.sh
./doth-init.sh
```

The script puts the Go toolchain and the `doth` binary inside the project's `.doth/` directory. Both are out of the way of your files. The `doth` binary is on the `PATH` only for the duration of the script.

The init script also generates a [doth wrapper](./wrapper.md).

## The wrapper script

The wrapper script installs and updates `doth` on demand. It keeps the Go toolchain and the `doth` binary inside the project's `.doth/` directory. It is a portable way to run `doth` without installing it system wide.

There are two ways to create a wrapper script.

The `doth init` command can write one for you. Pass `--wrapper` on init.

```sh
doth init --wrapper
```

The `doth wrapper` command prints the wrapper to stdout. Redirect it into a file.

```sh
doth wrapper > doth.sh
chmod +x doth.sh
```

The wrapper behaves as follows.

- It checks for a `doth` binary on the `PATH`. It runs the binary if it is there.
- It downloads Go and installs the latest `doth` when no binary is on the `PATH`.
- It updates the installed `doth` when the latest version is newer than the running version. It rewrites itself with `doth wrapper` after the update.
- It reads a `doth.lock` file in the project root. It uses the version inside the file when present. It does not auto-update locked projects.

Run `doth` through the wrapper from the project root.

```sh
./doth.sh deploy
```

See [The Wrapper Script](./wrapper.md) for the value the wrapper adds to the workflow.

## Locking the version

A locked project does not auto-update. This is useful for production machines. Use `doth lock` to write a `doth.lock` file. The file contains the version of the running `doth` binary by default. Pass `--version` to pin to a different version.

```sh
doth lock
doth lock --version v0.0.13
```

Use `doth unlock` to remove the lock file. The wrapper resumes auto-updating.

The lock only affects the wrapper. It has no effect on `doth` binaries installed by other means.

## `doth init`

`doth init` creates a new project in the current directory. It writes a `doth.yaml`, a `modules` directory, a `.doth/` directory for local state, and a `.gitignore`. Pass `--wrapper` to also write a `doth.sh` wrapper script.

The flags are listed below.

| Flag            | Description                                                          |
|-----------------|----------------------------------------------------------------------|
| `--modules`     | Directory for modules. Defaults to `./modules`.                      |
| `--wrapper`     | Also write a `doth.sh` wrapper script.                               |
| `--destructive` | Delete and recreate an existing project. Asks for confirmation.       |
| `--dry`         | Print the actions that would be taken. Do not perform them.          |
| `--verbose`     | Print verbose output.                                                |
| `--autoconfirm` | Answer all prompts with yes.                                         |

The default `doth.yaml` looks like this.

```yaml
modulePath: "./modules"
requireConfig: false
deps:
  - name: curl
    packages:
      pacman: curl
      apt: curl
      brew: curl
dothVersionDoNotEditManually: 1
```

The default `.gitignore` ignores the `.doth/` directory.

## Running as root

`doth deploy` and `doth add` warn and ask for confirmation when run as root. The recommended approach is to run `doth` as the user that owns the target files. The confirmation can be bypassed with `--autoconfirm`.
