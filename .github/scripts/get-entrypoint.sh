#!/bin/bash

FILE_BASENAME=$(basename "$1" .go)

# Path to the YAML file (relative to repository root)
YAML_PATH=".github/entrypoints.yaml"

# Debug: Print current directory and YAML path
echo "Current directory: $(pwd)"
echo "Looking for YAML at: $YAML_PATH"

# Verify YAML file exists
if [ ! -f "$YAML_PATH" ]; then
  echo "Error: YAML file not found at $YAML_PATH"
  exit 1
fi

# Extract entrypoint
ENTRYPOINT=$(yq -r ".$FILE_BASENAME" "$YAML_PATH")

# Debug: Print what we're looking for and what we found
echo "Looking for key: $FILE_BASENAME"
echo "Found entrypoint: $ENTRYPOINT"

if [ -z "$ENTRYPOINT" ] || [ "$ENTRYPOINT" == "null" ]; then
  echo "Error: Entrypoint not found for file: $FILE_BASENAME"
  echo "Available keys in YAML:"
  yq -r 'keys | .[]' "$YAML_PATH"
  exit 1
fi

echo "entrypoint=$ENTRYPOINT" >> $GITHUB_OUTPUT