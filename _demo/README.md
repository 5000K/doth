# doth demo

## Run it

Execute `./doth.sh deploy -c config.yaml -d` to see what it would do.

Execute `./doth.sh deploy -c config.yaml` to actually run it (or `./doth.sh deploy -c config.yaml -v` for a verbose version).

You don't need anything to run this demo besides the demo folder itself. doth.sh is the so called "doth wrapper", which will automatically setup the most current version of doth for you.

## What it does

The demo has a single module (modules/example), which has six configured steps in it's module.yaml:

```yaml
files:
  - name: copy.yml # copies copy.yml
    strategy: copy
  - name: link.yml # creates a symlink for link.yml
    strategy: link
  - name: render.yml # renders render.yml using the provided config
    strategy: render
  - name: copy-folder # copy an entire folder
    strategy: copy
  - name: link-folder # link an entire folder
    strategy: link
  - name: render-folder # render an entire folder using the provided config
    strategy: render
```

The output folder is configured to be ./target/example - so take a look into this folder after you ran the demo!
