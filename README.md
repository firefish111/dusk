# Dusk

This is **the** package manager for [The Night Programming language](https://github.com/DynamicSquid/night).
This provides a seamless interface for configuring and managing packages.

## Features

Here are the current features:

- `dusk add` installs the package(s) of your choice, just list them after the command.
- `dusk del` deletes the package(s) provided.
- `dusk upd` updates the package(s) provided.
- `dusk inf` gets information about the given package(s) from the server.
- `dusk ls` lists all installed packages.
- **IN ALPHA** `dusk find` finds packages containing the given parameters
---

## Getting Started

Dusk is already distributed through Night. You can find installation instructions at the [Night Website](https://night-website.dynamicsquid.repl.co/index.html).

You can also install Dusk from source if you'd wish. Here's how:

### Installing from Source

1. Install the Golang compiler, using the instructions at [the official Golang website](https://golang.org/doc/install)

2. Clone this repo with:
   `git clone https://github.com/firefish111/dusk.git`
   or using the GitHub CLI:
   `gh repo clone firefish111/dusk`

3. Enter the source folder: `cd dusk`

4. Compile dusk using `go build dusk.go`

And that's it!
