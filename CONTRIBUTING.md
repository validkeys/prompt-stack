# Contributing to prompt-stack

Thank you for your interest in contributing to prompt-stack!

## Development Setup

1. Clone the repository:

```sh
git clone https://github.com/kyledavis/prompt-stack.git
cd prompt-stack
```

2. Install dependencies:

```sh
go mod download
```

3. Build the project:

```sh
make build
```

## Code Style

- Follow Go conventions and formatting (`gofmt`)
- Write clear, descriptive docstrings for exported functions
- Use table-driven tests with `t.Run()` for subtests
- Keep functions focused on a single responsibility
- Validate external inputs at entry points

## Testing

All contributions must include tests:

- Write unit tests for new functionality
- Target 80%+ test coverage
- Use table-driven test patterns
- Run `make test` before committing

Example test structure:

```go
func TestYourFunction(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "valid input",
			input: "test",
			want:  "output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YourFunction(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("YourFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("YourFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

## Commit Messages

Follow conventional commits format:

- `feat:` - New features
- `fix:` - Bug fixes
- `test:` - Test additions or changes
- `docs:` - Documentation changes

Example:

```
feat: add dry-run mode to executor command
fix: handle nil pointer in prompt validation
test: add coverage for edge cases
docs: update README with usage examples
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes with tests
4. Ensure all tests pass (`make test`)
5. Ensure linting passes (`make lint`)
6. Commit your changes with conventional commit messages
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Questions?

Feel free to open an issue for questions or discussion.
