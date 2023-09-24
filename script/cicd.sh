#!/usr/bin/env bash

SERVERPID=0
RESPCODE=0

function info {
    echo -e "
CI COMMANDS: style lint test build run start check check-no-lint docker-gobuilder docker-build docker-run docker-up compose-up
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
    echo -e "\n>>_AppStart_<<"
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
    echo -e "\n>>_AppClose_<<"
    kill $SERVERPID
    if [[ $? -gt 0 ]]; then checksBreaked; fi
    sleep 1
}

function containe-gobuilder {
    echo -e "\n>>_DockerGoBuilderCreate_<<"
    docker buildx build --platform linux/amd64 --no-cache -f ./script/gobuilder1.21.1.dock -t kselnaag/gobuilder:1.21.1 . --load
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function containe-app {
    echo -e "\n>>_DockerAppCreate_<<"
    docker buildx build --platform linux/amd64 --no-cache -f ./script/shortlink.dock -t kselnaag/shortlink . --load
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function container-run {
    echo -e "\n>>_DockerAppRun_<<"
    docker run -it --user 10001 -p 8080:8080/tcp kselnaag/shortlink:latest
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function container-up {
    echo -e "\n>>_DockerAppUp_<<"
    docker run --user 10001 -p 8080:8080/tcp kselnaag/shortlink:latest
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function compose-up {
    echo -e "\n>>_ComposeAllUp_<<"
    #
    if [[ $? -gt 0 ]]; then checksBreaked; fi
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
    "check")
        styleCheck
        lint
        unitTest
        build
        start
        healthCheck
        stop
        ;;
    "check-no-lint")
        styleCheck
        unitTest
        build
        start
        healthCheck
        stop
        ;;
    "docker-gobuilder")
        containe-gobuilder
        ;;
    "docker-build")
        containe-app
        ;;
    "docker-run")
        containe-app
        container-run
        ;;
    "docker-up")
        containe-app
        container-up
        ;;
    "compose-up")
        compose-up
        ;;
    *)
        info
        ;;
    esac
fi

echo -e "\n>>_Successfull_<<\n"
