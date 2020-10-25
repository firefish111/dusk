# dusk

The [Night](https://github.com/DynamicSquid/night) package manager. It provides a seemless interface for configuring and managing packages.

To view more information about Dusk, including the commands, be sure to check out the [website](https://night-website.dynamicsquid.repl.co/index.html) as well.

---

### Getting Started with Dusk

Dusk is already distributed through Night. You can find installation information on the [Night website](https://night-website.dynamicsquid.repl.co/index.html).

You can also install Dusk from source if you'd wish.

### Building from Source

1. Install the *Go* compiler

2. Clone this repository

```
git clone https://github.com/firefish111/dusk.git
cd dusk
```

3. Compile Dusk

```
go build -ldflags \"-s -w\" src/dusk.go
```

And you're done!
