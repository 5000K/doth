---
title: Dependencies
description: Declaring dependencies and installing them with built-in or custom sources.
author: 5000K
---

doth can install the packages your configuration depends on. Dependencies are declared in the `doth.yaml` and in each `module.yaml`. `doth install` reads the declarations and runs the install commands.

## Declaring dependencies

A dependency has a name and a map of source names to package names. The source names are matched against the package sources you pass to `doth install`. The same shape is used in `doth.yaml` and in `module.yaml`.

```yaml
# doth.yaml
deps:
  - name: curl
    packages:
      pacman: curl
      apt: curl
      brew: curl
```

The name is used to deduplicate dependencies across the project. When the same name appears in multiple modules, the last declaration wins. See [YML Reference](./yml.md) for the full dependency format.

## Built-in sources

`doth install` knows a few common package managers out of the box. Pass one or more of the flags below to enable a built-in source. The command is run with `/bin/sh -c`. The `{package}` placeholder is replaced with the package name from the dependency declaration.

| Flag        | Command                                |
|-------------|----------------------------------------|
| `--apt`     | `sudo apt install --yes {package}`     |
| `--apt-get` | `sudo apt-get install --yes {package}` |
| `--dnf`     | `sudo dnf install --assumeyes {package}` |
| `--pacman`  | `sudo pacman -S --noconfirm {package}` |
| `--yay`     | `yay -S --noconfirm {package}`         |
| `--paru`    | `paru -S --noconfirm {package}`        |
| `--go`      | `go install {package}`                 |
| `--npm`     | `npm install -g {package}`             |
| `--brew`    | `yes \| brew install {package}`        |

```sh
doth install --pacman --yay
```

## Custom sources

Define custom sources in a configuration file. Pass the file to `doth install` with `--config`. The same file can also be passed to `doth deploy`. Multiple configuration files are merged in order. Later files take precedence over earlier files for the same key.

```yaml
# pacman.yaml
packageSources:
  - name: pacman
    command: sudo pacman -S {package}
```

The `{package}` placeholder is replaced with the package name from the dependency. The command runs unescaped. Be careful when adding custom sources. Built-in sources and custom sources can be mixed.

```sh
doth install --config pacman.yaml
```

The `packageSources` list is itself merged across files. Sources with the same `name` are deduplicated.

## Matching

`doth install` walks the sources in order. The first source that has an entry in the dependency's `packages` map wins. Sources without a matching entry are skipped. Dependencies with no matching source are skipped and a message is printed.

Custom sources from configuration files are tried first. Built-in sources come after them. The built-in sources follow a fixed order. The order is `apt`, `apt-get`, `dnf`, `pacman`, `yay`, `paru`, `go`, `npm`, `brew`. The order of built-in flags on the command line does not matter.

## Confirmation

`doth install` always asks for confirmation before running commands. The shell is given access to your environment. Pass `--dry` to print the commands that would run without executing them.

The `--silent` flag runs each command without piping its output. Use it for unattended scripts.

## Examples

A `doth.yaml` declaring two dependencies.

```yaml
deps:
  - name: ripgrep
    packages:
      pacman: ripgrep
      apt: ripgrep
      brew: ripgrep
  - name: i3
    packages:
      pacman: i3-wm
```

Install both on a pacman system.

```sh
doth install --pacman
```

A `pacman.yaml` that defines a custom source. It is useful when you want a non default install command.

```yaml
packageSources:
  - name: pacman
    command: sudo pacman -S --needed {package}
```

A `module.yaml` that declares its own dependencies. Module dependencies are merged with the top level dependencies of `doth.yaml`.

```yaml
# modules/i3/module.yaml
target: ~/.config/i3/
files:
  - name: config
    strategy: render
deps:
  - name: i3-gaps
    packages:
      aur: i3-gaps
```
