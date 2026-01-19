# prompt-stack

AI-assisted development workflow tool with Plan/Build modes for generating and validating Ralphy YAML files.

## Quickstart

Build the binary:

```sh
make build
```

Run the program:

```sh
./dist/prompt-stack --help
```

## Features

- **Plan mode**: Generate implementation plans from requirements using AI assistance
- **Build mode**: Build project components based on implementation plan tasks
- **Validate**: Validate implementation plans against schema and quality standards
- **Review**: Review implementation progress and quality metrics
- **Init**: Interactive requirements gathering for new milestones

## Usage

### Initialize a new milestone

Run an interactive interview to gather milestone requirements:

```sh
./dist/prompt-stack init
./dist/prompt-stack init --output-dir docs/implementation-plan/m1
```

### Generate implementation plans

```sh
./dist/prompt-stack plan --input docs/requirements.md --output docs/implementation-plan/m0/
```

### Validate implementation plans

```sh
./dist/prompt-stack validate --input docs/implementation-plan/m0/final_implementation-plan.yaml
```

### Build from implementation plan

```sh
./dist/prompt-stack build --plan docs/implementation-plan/m0/final_implementation-plan.yaml
```

### Review implementation progress

```sh
./dist/prompt-stack review
```

## Project Structure

```
prompt-stack/
├── cmd/prompt-stack/       # CLI commands
├── pkg/                 # Reusable packages
│   ├── executor/       # Ralphy executor integration
│   └── prompt/         # Interactive prompt system
├── docs/               # Documentation
├── tests/              # Integration tests
└── .prompt-stack/         # Tool configuration and vendored scripts
```

## Development

### Build

```sh
make build
```

### Run tests

```sh
make test
```

### Format code

```sh
make fmt
```

### Lint

```sh
make lint
```

### Clean build artifacts

```sh
make clean
```

## Documentation

- [Architecture](docs/architecture.md) - System architecture overview
- [Commands](docs/commands.md) - Detailed command documentation
- [Best Practices](docs/best-practices.md) - Development guidelines

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to this project.

## License

TBD
