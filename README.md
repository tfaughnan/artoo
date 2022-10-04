# artoo

Barebones IRC bot / astromech droid written in Go.

## Compilation

`make`

## Configuration

Configuration is read from the first found file in the following list:

1. File passed via the `-c` flag
2. `~/.config/artoo.toml`
3. `/etc/artoo.toml`

The file `artoo.toml.example` is provided for reference.
