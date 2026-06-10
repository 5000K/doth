---
title: Modules
description: Writing modules, choosing file strategies, and using templates.
author: 5000K
---

# Modules

A module is a folder inside the `modules` directory of your project. The folder contains a `module.yaml` and the files and folders the module deploys. The `module.yaml` describes where the files go and how. The `modules` directory is configurable through the top level `modulePath` field of `doth.yaml`.

## The shape of a module

A module is a folder with a `module.yaml` and the files it owns.

```
modules/
└── bash/
    ├── module.yaml
    └── .bashrc
```

The folder name is the module name. It must be a valid directory name. Empty names, `.` and `..`, and names containing path separators are rejected. You may use `~` for paths relative to the current users home directory.

## Adding a module

`doth add` creates a new module from an existing folder of files. It copies the files into the new module folder. It writes a `module.yaml` that uses the `copy` strategy for each imported file.

```sh
doth add --name bash --target ~/.config/bash
```

The `--name` flag is the module's folder name. The folder appears under `modulePath`. The `--target` flag is the source path to import from. It should be a directory of files. The `module.yaml` is set to deploy into that same directory by default.

The flags are listed below.

| Flag              | Description                                                                 |
| ----------------- | --------------------------------------------------------------------------- |
| `--name`          | Internal name of the module. Required.                                      |
| `--target`        | Source path of the files to import. Required.                               |
| `--glob`          | Glob pattern of files to include. Relative to the target path.              |
| `--skip-existing` | Do not copy files from the target into the module.                          |
| `--destructive`   | Delete and recreate the module if it already exists. Asks for confirmation. |
| `--dry`           | Print the actions that would be taken. Do not perform them.                 |
| `--verbose`       | Print verbose output.                                                       |
| `--autoconfirm`   | Answer all prompts with yes.                                                |

The `--target` path is also written to the module's `target` field. Files in the module are deployed into that path. Edit the field when you want the module's deploy target to differ from its import source.

### Manual module creation

You can also create a module by hand. Make the folder, drop your files in, and write a `module.yaml`.

```
mkdir -p modules/bash
cp ~/.bashrc modules/bash/.bashrc
```

```yaml
# modules/bash/module.yaml
target: ~/
files:
  - name: .bashrc
    strategy: copy
```

## The module file

A `module.yaml` has four fields.

| Field    | Type   | Description                                                                                |
| -------- | ------ | ------------------------------------------------------------------------------------------ |
| `target` | string | Base path the module's files are deployed into. Each file lands at `target/<file name>`.   |
| `skip`   | bool   | When `true`, the module is ignored by `doth deploy`.                                       |
| `files`  | list   | The file entries to deploy. See below.                                                     |
| `deps`   | list   | Dependencies required by the module. See [Dependencies](./dependencies.md) for the format. |

A file entry under `files` has three fields.

| Field      | Type   | Description                                                                                                     |
| ---------- | ------ | --------------------------------------------------------------------------------------------------------------- |
| `name`     | string | The file name or relative path inside the module folder. May contain glob patterns.                             |
| `strategy` | string | Required. One of `copy`, `link`, or `render`.                                                                   |
| `target`   | string | Optional override for the deployment path. The file lands at this exact path instead of `module.target/<file>`. |

## File strategies

### copy

`copy` reads the file from the module folder and writes a copy to the target location. Existing files at the target are replaced. Directories are walked recursively. Each file inside the directory is copied individually.

### link

`link` creates a symbolic link at the target location. The link points into the module folder. Edits made in the module folder are visible at the target immediately. Existing files at the target are removed and replaced with a link.

### render

`render` reads the file as a Go `text/template`. It executes the template with the values from the configuration files. It writes the result to the target location. Existing files at the target are replaced.

Templates use Go's `text/template` syntax. Map values are accessed through the `index` function. Use `default` to provide a fallback.

```jsonc
// modules/waybar/config.jsonc
{
    // ...
    "battery": {{ index . "show-battery" | default false }},
    "theme":   "{{ index . "theme" | default "light" }}" // could also be simplified to {{ .theme | default "light" }}, since the key 'theme' does not contain a dash -
}
```

### Configuration

The configuration file is a YAML or JSON object. Its top level keys become the template's data context.

For the waybar example above, imagine this configuration file:

```yaml
# laptop.yaml
show-battery: true
theme: dark
```

Pass a configuration file to `doth deploy`.

```sh
doth deploy --config laptop.yaml
```

See [the Go template documentation](https://pkg.go.dev/text/template) for the full template syntax.

## Skipping a module

Set `skip: true` in a `module.yaml` to keep the module in the repository without deploying it. This is useful for keeping alternate configurations around.

```yaml
# modules/work-laptop/module.yaml
skip: true
target: ~/
files:
  - name: .bashrc
    strategy: copy
```
