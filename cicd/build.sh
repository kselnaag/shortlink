#!/usr/bin/env bash

go build -o ./bin/shortlink ./cmd/main.go
./bin/shortlink
