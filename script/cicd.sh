#!/usr/bin/env bash

SERVERPID=0
RESPCODE=0

function info {
    echo -e "
CI COMMANDS: style lint test build run start check check-no-lint 
             docker-gobuilder docker-build docker-run docker-up compose-up
             metrics metrics-graph
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
    docker run -it --name SLsrv --user 10001 -p 8080:8080/tcp kselnaag/shortlink
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function container-up {
    echo -e "\n>>_DockerAppUp_<<"
    docker run -d --name SLsrv --user 10001 -p 8080:8080/tcp kselnaag/shortlink
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function container-start {
    echo -e "\n>>_DockerAppStart_<<"
    docker start SLsrv
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function container-stop {
    echo -e "\n>>_DockerAppStop_<<"
    docker stop SLsrv
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function compose-up {
    echo -e "\n>>_ComposeAllUp_<<"
    #
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function metrics-graph { # https://www.yworks.com/yed-live/
    echo -e "\n>>_MetricsGraph_<<"
    image_packages . ./script/sl.graphml
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function loc {
    echo -e "\n>>_LinesOfCode_<<"
    FPATH=./script/sl.loc
    gcloc . > $FPATH
    cat $FPATH | grep Language
    if [[ $? -gt 0 ]]; then checksBreaked; fi
    cat $FPATH | grep Golang
    cat $FPATH | grep Bash
}

function cycl {
    echo -e "\n>>_CyclomaticComplexity_<<"
    FPATH=./script/sl.cycl
    gocyclo -avg . > $FPATH
    cat $FPATH | grep Average
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function cogn {
    echo -e "\n>>_CognitiveComplexity_<<"
    FPATH=./script/sl.cogn
    gocognit -avg . > $FPATH
    cat $FPATH | grep Average
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function testcov {
    echo -e "\n>>_TestCoverage_<<"
    FTEMP=./script/coverage.test
    FOUT=./script/sl.cov
    go test -vet=off -count=1 -coverprofile=$FTEMP ./...
    go tool cover -func=$FTEMP > $FOUT
    cat $FOUT | grep total
     if [[ $? -gt 0 ]]; then checksBreaked; fi
    rm $FTEMP
}

function maintain { # !
    echo -e "\n>>_MaintainabilityIndex_<<"
    FPATH=./script/sl.maint
    go vet -vettool=$(which complexity) --maintunder 100 ./... > $FPATH
    cat $FPATH | grep 'index='
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
    "metrics-graph")
        metrics-graph
        ;;
    "metrics")
        loc
        cycl
        cogn
        testcov
        ;;    
    "dock") # experiments
        PWD=`pwd`
        SRV=SLsrv
        DB=SLpg
        # build
        docker cp "$PWD/bin/shortlink" $SRV:"/shortlink"
        docker cp "$PWD/config/shortlink.env" $SRV:"/shortlink.env"
        docker start $SRV
        docker logs $SRV
        docker stop $SRV
        # docker run -it --name SLsrv --user 10001 -p 8080:8080/tcp kselnaag/shortlink
        # docker run -d --name SLpg -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=example postgres:16.0-alpine3.18
        # docker start SLpostgres && docker exec -it SLpostgres ls -la /var/lib/postgresql
        ;;
    *)
        info
        ;;
    esac
fi

echo -e "\n>>_Successfull_<<\n"

: <<COMMENT
METRICS:
+
0. Package graph
1. Lines of Code
2. Cyclomatic Complexity
3. Cognitive Complexity
12. Test Coverage
19. Code size + Repository size

...
4. Halstead Complexity
15. Maintainability index
22. Commit itme (Checks time + Build time)
24. Code duplication

?
5. Coupling
7. Hits of Code (Code Churn)
8. LCOMx
9. PCC
10. LCC
11. MMAC and NHD
13. DSQI
14. Instruction path length
16. Function Points

-
6. GitHub stars
17. Mutation Coverage
18. Number of methods, classes, etc.
20. Forks and pull requests !
21. Bugs !
23. Algorithmic complexity Big-O

COMMENT

