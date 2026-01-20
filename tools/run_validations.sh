#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
REPORTS_DIR="$PROJECT_ROOT/docs/implementation-plan/m0"

echo "Running validation tools..."
echo "Reports will be saved to: $REPORTS_DIR"
echo ""

# Build prompt-stack binary if not present
echo "Building prompt-stack binary..."
go build -o "$PROJECT_ROOT/prompt-stack" ./cmd/prompt-stack

# Run task sizing validation
echo "Running task sizing validation..."
cd "$PROJECT_ROOT"
"$PROJECT_ROOT/prompt-stack" validate-task-sizing --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_sizing.json" 2>&1 || {
    # Capture output even if validation fails
    "$PROJECT_ROOT/prompt-stack" validate-task-sizing --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_sizing.json" 2>&1
}

# Run enforcement validation
echo "Running enforcement validation..."
"$PROJECT_ROOT/prompt-stack" validate-enforcement --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_enforcement.json" 2>&1 || {
    "$PROJECT_ROOT/prompt-stack" validate-enforcement --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_enforcement.json" 2>&1
}

# Run constraints validation
echo "Running constraints validation..."
"$PROJECT_ROOT/prompt-stack" validate-constraints --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_constraints.json" 2>&1 || {
    "$PROJECT_ROOT/prompt-stack" validate-constraints --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_constraints.json" 2>&1
}

# Run implementation guidelines validation
echo "Running implementation guidelines validation..."
"$PROJECT_ROOT/prompt-stack" validate-implementation-guidelines --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_implementation_guidelines.json" 2>&1 || {
    "$PROJECT_ROOT/prompt-stack" validate-implementation-guidelines --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/validation_report_implementation_guidelines.json" 2>&1
}

# Run YAML validation
echo "Running YAML validation..."
"$PROJECT_ROOT/prompt-stack" validate-yaml --schema docs/ralphy-inputs.schema.json --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/yaml_validation_report.json" 2>&1 || {
    "$PROJECT_ROOT/prompt-stack" validate-yaml --schema docs/ralphy-inputs.schema.json --file "$REPORTS_DIR/ralphy_inputs.yaml" > "$REPORTS_DIR/yaml_validation_report.json" 2>&1
}

echo ""
echo "All validations completed!"
echo "Generated reports:"
ls -la "$REPORTS_DIR"/*validation*.json
