#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
REPORTS_DIR="$PROJECT_ROOT/docs/implementation-plan/m0"

echo "Running validation tools..."
echo "Reports will be saved to: $REPORTS_DIR"
echo ""

# Run task sizing validation
echo "Running task sizing validation..."
cd "$SCRIPT_DIR"
go run task_sizing/validate.go -file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_sizing.json" 2>&1 || {
    # Capture output even if validation fails
    go run task_sizing/validate.go -file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_sizing.json" 2>&1
}

# Run enforcement validation
echo "Running enforcement validation..."
cd "$SCRIPT_DIR"
go run validate_enforcement/validate.go -file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_enforcement.json" 2>&1 || {
    # Capture output even if validation fails
    go run validate_enforcement/validate.go -file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_enforcement.json" 2>&1
}

# Run constraints validation
echo "Running constraints validation..."
cd "$SCRIPT_DIR"
go run validate_constraints/validate.go "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_constraints.json" 2>&1 || {
    # Capture output even if validation fails
    go run validate_constraints/validate.go "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_constraints.json" 2>&1
}

# Run implementation guidelines validation
echo "Running implementation guidelines validation..."
cd "$SCRIPT_DIR"
go run validate_implementation_guidelines/validate.go "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_implementation_guidelines.json" 2>&1 || {
    # Capture output even if validation fails
    go run validate_implementation_guidelines/validate.go "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_implementation_guidelines.json" 2>&1
}

# Run YAML validation
echo "Running YAML validation..."
cd "$SCRIPT_DIR"
go run validate_yaml.go "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/yaml_validation_report.json" 2>&1 || {
    # Capture output even if validation fails
    go run validate_yaml.go "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/yaml_validation_report.json" 2>&1
}

echo ""
echo "All validations completed!"
echo "Generated reports:"
ls -la "$REPORTS_DIR"/*validation*.json