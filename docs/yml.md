---
title: YML Reference
description: The full format reference for doth.yaml, module.yaml, and configuration files.
author: 5000K
---

# YML Reference

This page documents the YAML formats used by doth. It is a reference. See [[Setup]], [[Modules]], and [[Dependencies]] for usage guides.

## `doth.yaml`

The top level configuration of a project. Located at the project root as `doth.yaml` or `doth.yml`.

| Field                          | Type   | Description                                                                                          |
|--------------------------------|--------|------------------------------------------------------------------------------------------------------|
| `modulePath`                   | string | Path to the directory that contains modules. Relative to the project root. Defaults to `./modules`. |
| `deps`                         | list   | Top level dependencies. See [[Dependencies]] for the format.                                         |
| `requireConfig`                | bool   | When `true`, `doth deploy` should refuse to run without configuration files. Currently not enforced. |
| `dothVersionDoNotEditManually` | int    | The format version of the file. Do not edit by hand. Currently not checked.                          |

## `module.yaml`

Located inside each module folder as `module.yaml` or `module.yml`.

| Field    | Type   | Description                                                                       |
|----------|--------|-----------------------------------------------------------------------------------|
| `target` | string | Base path the module's files are deployed into. Each file lands at `target/<file>`. |
| `skip`   | bool   | When `true`, the module is ignored by `doth deploy`.                              |
| `files`  | list   | File entries that describe what to deploy. See below.                             |
| `deps`   | list   | Module level dependencies. Merged with the top level `deps`.                     |

## File entry

Each entry under `files` describes a single file or glob pattern.

| Field      | Type   | Description                                                                                          |
|------------|--------|------------------------------------------------------------------------------------------------------|
| `name`     | string | The file name or relative path inside the module folder. May contain glob patterns.                  |
| `strategy` | string | Required. One of `copy`, `link`, or `render`. Empty or unknown values cause an error on deploy.     |
| `target`   | string | Override for the deployment path. The file lands at this path instead of `module.target/<file>`.   |

## Strategies

| Value    | Behavior                                                                                              |
|----------|-------------------------------------------------------------------------------------------------------|
| `copy`   | Copies the file from the module folder to the target location. Directories are walked recursively.     |
| `link`   | Creates a symbolic link at the target location that points into the module folder.                    |
| `render` | Executes the file as a Go `text/template` with the values from the deployment config files.           |

## Dependency

A single entry in a `deps` list.

| Field      | Type   | Description                                                                                  |
|------------|--------|----------------------------------------------------------------------------------------------|
| `name`     | string | A human readable name. Used to deduplicate dependencies across the project.                 |
| `packages` | map    | A map from source name to the package name in that source. The first matching source wins.  |

## Configuration file

A configuration file is passed to `doth deploy` or `doth install` with the `--config` flag. The file is YAML or JSON. Multiple files are merged in order. Later files take precedence over earlier files for the same key. Maps are merged recursively. Lists are concatenated and deduplicated. Scalars are replaced.

The full list of recognized fields follows. Unrecognized fields are preserved in the merged result. They are not used by `doth` directly.

### `packageSources`

Used by `doth install`. A list of package source definitions.

| Field     | Type   | Description                                                                                  |
|-----------|--------|----------------------------------------------------------------------------------------------|
| `name`    | string | The source name. Referenced from a dependency's `packages` map.                              |
| `command` | string | The shell command. The `{package}` placeholder is replaced with the package name from the dependency. |

### Template data

Used by `doth deploy`. The top level keys of the configuration file become the template's data context. Templates access them through Go's `index` function.

```yaml
# laptop.yaml
show-battery: true
theme: dark
```

```jsonc
// modules/waybar/config.jsonc
{
    "battery": {{ index . "show-battery" | default false }},
    "theme":   "{{ index . "theme" | default "light" }}"
}
```

```sh
doth deploy --config laptop.yaml
```

The dot is the data context. To access a nested value, walk the structure with `index`. Use `default` to provide a fallback. See [the Go template documentation](https://pkg.go.dev/text/template) for the full syntax.
