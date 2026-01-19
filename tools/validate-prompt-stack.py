#!/usr/bin/env python3
"""
Simple prompt-stack config validator.
Usage: ./validate-prompt-stack.py prompt-stack.yaml
"""

import sys
import json
import yaml
from jsonschema import validate, ValidationError

SCHEMA_PATH = "docs/prompt-stack.schema.json"


def main():
    if len(sys.argv) < 2:
        print("Usage: validate-prompt-stack.py <path-to-prompt-stack.yaml>")
        sys.exit(2)
    cfg_path = sys.argv[1]
    with open(cfg_path, "r", encoding="utf-8") as f:
        cfg = yaml.safe_load(f)
    with open(SCHEMA_PATH, "r", encoding="utf-8") as f:
        schema = json.load(f)
    try:
        validate(instance=cfg, schema=schema)
    except ValidationError as e:
        print("Validation failed:")
        print(e)
        sys.exit(1)
    print("Validation passed (advisory).")
    # Basic advisory checks
    for role, defn in cfg.get("roles", {}).items():
        for m in defn.get("candidates", []):
            if m not in cfg.get("models", {}):
                print(f'Warning: role "{role}" references unknown model "{m}"')
    print("Done.")


if __name__ == "__main__":
    main()
