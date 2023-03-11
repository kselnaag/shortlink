#!/usr/bin/env bash

echo -e "\n>>_ShortLink build_<<\n"
go build -o ./bin/shortlink ./cmd/main.go
if [[ $? -gt 0 ]]; then exit 1; fi

./bin/shortlink
