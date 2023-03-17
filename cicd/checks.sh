#!/usr/bin/env bash

echo -e "\n>>_StyleChecking_<<"
go fmt ./...
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_Linting_<<"
golangci-lint run ./...
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_UnitTests_<<"
set -o pipefail
go test -vet=off -count=1 -race ./... | { grep -v 'no test files'; true; }
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_ShortLink build_<<"
go build -o ./bin/shortlink ./cmd/main.go
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_ServerStart_<<"
./bin/shortlink 1>/dev/null &
SERVERPID=$!
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_HealthCheck_<<"
sleep 2
RESPCODE=`curl -i http://localhost:8080/check/ping 2>/dev/null | head -n 1 | cut -d$' ' -f2`
if [ "$RESPCODE" != "200" ]; then 
    exit 1;
fi

echo -e "\n>>_CloseServer_<<\n"
kill $SERVERPID
