#!/usr/bin/env bash

echo -e "\n>>_StyleChecking_<<"
go fmt ./...
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_Linting_<<"
golangci-lint run ./...
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_UnitTests_<<"
go test -vet=off -count=1 -race ./...
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_ShortLink build_<<\n"
go build -o ./bin/shortlink ./cmd/main.go
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_HelthChecks_<<\n"
echo -e "empty\n"
