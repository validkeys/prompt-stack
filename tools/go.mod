// go.mod defines the Go module for the prompt-stack tools.
//
// Module Path
//
//   github.com/kyledavis/prompt-stack/tools
//
// Version
//
//   Go 1.20 or later required
//
// Dependencies
//
//   - github.com/santhosh-tekuri/jsonschema/v5 (v5.1.0)
//     Enterprise-grade JSON Schema validator supporting Draft 7, 2019-09, and 2020-12
//     Features:
//       - Fast validation with minimal overhead
//       - Detailed error reporting with path and keyword information
//       - $ref resolution for schema composition
//       - Custom formats and keyword validation
//
//   - sigs.k8s.io/yaml (v1.3.0)
//     Kubernetes project's YAML-to-JSON converter
//     Features:
//       - Robust handling of complex YAML types (maps, sequences, scalars)
//       - Preserves YAML comments during conversion
//       - Well-tested with enterprise workloads
//       - Active maintenance and security updates
//
// Maintenance
//
//   To update dependencies:
//     cd tools
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
//   - jsonschema: https://github.com/santhosh-tekuri/jsonschema
//   - k8s.io/yaml: https://github.com/kubernetes-sigs/yaml
//
//   Run `go list -json -m all` to view all dependencies and their versions.
//
module github.com/kyledavis/prompt-stack/tools

go 1.20

require (
	github.com/santhosh-tekuri/jsonschema/v5 v5.1.0
	gopkg.in/yaml.v3 v3.0.1
	sigs.k8s.io/yaml v1.3.0
)

require gopkg.in/yaml.v2 v2.4.0 // indirect
