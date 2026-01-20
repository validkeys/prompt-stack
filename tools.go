//go:build tools

// Package tools pins tool-only dependencies used during development.
package tools

import (
	_ "github.com/santhosh-tekuri/jsonschema/v5"
	_ "gopkg.in/yaml.v3"
	_ "sigs.k8s.io/yaml"
)
