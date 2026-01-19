style-anchor — minimal Go single-binary example

Quickstart

Build the binary:

```sh
make build
```

Install locally (copies binary to ~/bin):

```sh
make install
```

Run the program:

```sh
./bin/mytool --name Alice
# or after install
mytool --name Alice
```

What this repo demonstrates

- Single-binary CLI layout under `cmd/mytool`
- Lean `Makefile` with `build`, `test`, `install`, and `clean` targets
- Minimal public package (`pkg/greeter`) with unit tests
- An install script in `bin/install.sh` to mirror single-binary distribution patterns

Developer notes

- Use `go test ./...` for running tests.
- Keep changes small and document intent in README or CHANGELOG.
