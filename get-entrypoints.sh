#!/bin/bash

FILE_BASENAME=$(basename "$1" .go)

ENTRYPOINT=$(yq -r ".$FILE_BASENAME" ./entrypoints.yaml)

if [ -z "$ENTRYPOINT" ] || [ "$ENTRYPOINT" == "null" ]; then
  echo "Entrypoint not found for file: $FILE_BASENAME"
  exit 1
fi

echo "entrypoint=$ENTRYPOINT" >> $GITHUB_OUTPUT