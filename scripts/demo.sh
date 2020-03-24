#!/bin/bash

printf "Starting Streaming Server\n\n"

go run cmd/stream-server/main.go &
SERVER_STATUS=$?

if [[ $SERVER_STATUS -ne 0 ]]; then
  printf "Failed to Run Streaming Server. Process Exiting"
  exit 1
fi

sleep 1

go run cmd/user-demo/main.go
CLI_STATUS=$?

if [[ $CLI_STATUS -ne 0 ]]; then
  printf "Failed to Run User demo command. Process Exiting"
  exit 1
fi
