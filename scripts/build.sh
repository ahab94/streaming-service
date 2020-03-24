#!/bin/bash

PROJECT_FILE="${PROJECT_FILE:-project.yml}"

target="cmd/stream-server"
echo "building: $target"

import_path=$(yaml read "${PROJECT_FILE}" metadata.import)
go build -mod=vendor -o "bin/$(basename "$target")" "${import_path}/$target"
