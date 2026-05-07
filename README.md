# doth

**doth** is a go-based .config manager for linux.

It allows for config-based templating, so that you can write reusable dotfiles for all your systems.

Show the battery state in the top bar on your laptop specifically? Support theming?

## Modules

Modules are subfolders in your doth project that contain a module.y\[a]ml file.
The module.yml describes what the module should do.

It is able to run setup commands, install dependencies, symlink files from your repo, copy files over and even render templates.

The module.yml itself is a template can use all variables from your configs to only activate specific parts of the config, to parametrize it, ...

## Templates

doth uses gos [text/template](https://pkg.go.dev/text/template) system. Look it up, it's powerful yet straight forward.

Basic examples will be added here later.

## Dependencies

You can define dependencies to install in your doth.yml. Using the `doth install` command, you can install them.
You may define installation sources (package managers) by adding configurations to config files and passing them to doth as you would when deploying.
You may also use the --auto flag, which will search for common package managers available on your system.
Make sure your deps use the correct keys for the source->package mapping when using auto.

These package managers are automatically detected:

- apt
- dnf
- pacman
- paru
- yay
- yum
- go
- npm
- brew
