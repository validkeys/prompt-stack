// go.mod defines the Go module for the task sizing validation tool.
//
// Module Path
//
//   github.com/kyledavis/prompt-stack/tools/task_sizing
//
// Version
//
//   Go 1.20 or later required
//
// Dependencies
//
//   - gopkg.in/yaml.v3 (v3.0.1)
//     YAML parsing library with full YAML 1.2 support
//     Features:
//       - Robust YAML parsing with proper error handling
//       - Support for YAML 1.2 specification
//       - Well-maintained with active community
//
// Maintenance
//
//   To update dependencies:
//     cd tools/task_sizing
//     go get -u ./...
//     go mod tidy
//     go mod verify
//
//   To verify dependency integrity:
//     go mod verify
//
// Security
//
//   Dependencies are sourced from well-maintained, actively maintained projects
//   with security policies in place:
//
//   - yaml.v3: https://github.com/go-yaml/yaml
//
//   Run `go list -json -m all` to view all dependencies and their versions.
//
module github.com/kyledavis/prompt-stack/tools/task_sizing

go 1.20

require gopkg.in/yaml.v3 v3.0.1
