#!/usr/bin/env bash

SERVERPID=0
RESPCODE=0

function info {
    echo -e "
CI COMMANDS: style lint test build run start checknolint check
EXAMLPE:     ./script/cicd.sh run
"
    exit 0
}

function checksBreaked {
    echo -e "\n>>_ChecksBreaked_<<\n"
    exit 1 
}

function styleCheck {
    echo -e "\n>>_StyleChecking_<<"
    go fmt ./...
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function lint {
    echo -e "\n>>_Linting_<<"
    golangci-lint run ./... --config=./config/.golangci.yaml
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function unitTest {
    echo -e "\n>>_UnitTests_<<"
    set -o pipefail
    go test -vet=off -count=1 -race ./... | { grep -v 'no test files'; true; }
    if [[ $? -gt 0 ]]; then checksBreaked; fi   
}

function build {
    echo -e "\n>>_Build_<<"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o ./bin/shortlink ./cmd/main.go
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function run {
    go build -o ./bin/shortlink ./cmd/main.go
    ./bin/shortlink
    exit 0
}

function start {
    echo -e "\n>>_ServiceStart_<<"
    ./bin/shortlink 1>/dev/null &
    SERVERPID=$!
    if [[ $? -gt 0 ]]; then checksBreaked; fi
    sleep 2
}

function healthCheck {
    echo -e "\n>>_HealthCheck_<<"
    RESPCODE=`curl -i http://localhost:8080/check/ping 2>/dev/null | head -n 1 | cut -d$' ' -f2`
    if [ "$RESPCODE" != "200" ]; then checksBreaked; fi
}

function stop {
    echo -e "\n>>_ServiceClose_<<"
    kill $SERVERPID
    if [[ $? -gt 0 ]]; then checksBreaked; fi
    sleep 1
}

if [[ $# -ne 1 ]]; then info; else
    case $1 in
    "style")
        styleCheck
        ;;
    "lint")
        lint
        ;;
    "test")
        unitTest
        ;;
    "build")
        build
        ;;
    "run")
        run
        ;;
    "start")
        build
        start
        healthCheck
        stop
        ;;
    "checknolint")
        styleCheck
        unitTest
        build
        start
        healthCheck
        stop
        ;;
    "check")
        styleCheck
        lint
        unitTest
        build
        start
        healthCheck
        stop
        ;;
    *)
        info
        ;;
    esac
fi

echo -e "\n>>_ChecksSuccessfull_<<\n"
