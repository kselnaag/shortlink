#!/usr/bin/env bash

function checksBreaked {
   echo -e "\n>>_ChecksBreaked_<<"
   exit 1 
}

echo -e "\n>>_StyleChecking_<<"
go fmt ./...
if [[ $? -gt 0 ]]; then checksBreaked; fi

echo -e "\n>>_Linting_<<"
golangci-lint run ./...
if [[ $? -gt 0 ]]; then checksBreaked; fi

echo -e "\n>>_UnitTests_<<"
set -o pipefail
go test -vet=off -count=1 -race ./... | { grep -v 'no test files'; true; }
if [[ $? -gt 0 ]]; then checksBreaked; fi

echo -e "\n>>_Build_<<"
go build -o ./bin/shortlink ./cmd/main.go
if [[ $? -gt 0 ]]; then checksBreaked; fi

echo -e "\n>>_ServiceStart_<<"
./bin/shortlink 1>/dev/null &
SERVERPID=$!
if [[ $? -gt 0 ]]; then checksBreaked; fi
sleep 1

echo -e "\n>>_HealthCheck_<<"
RESPCODE=`curl -i http://localhost:8080/check/ping 2>/dev/null | head -n 1 | cut -d$' ' -f2`
if [ "$RESPCODE" != "200" ]; then checksBreaked; fi

echo -e "\n>>_ServiceClose_<<\n"
kill $SERVERPID
if [[ $? -gt 0 ]]; then checksBreaked; fi
sleep 1

echo -e "\n>>_ChecksSuccessfull_<<\n"
