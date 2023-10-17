#!/usr/bin/env bash

SERVERPID=0
RESPCODE=0

function info {
    echo -e "\n
CI/CD COMMANDS: style lint test ttTest build run start check check-no-lint 
                docker-gobuilder docker-build docker-run docker-up compose-up
                metrics metrics-graph
EXAMLPE:        ./script/cicd.sh build
\n"
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
    go test -short -tags go_tarantool_ssl_disable -vet=off -count=1 -race ./... | { grep -v 'no test files'; true; }
    if [[ $? -gt 0 ]]; then checksBreaked; fi   
}

function intergTTtest {     # &"C:\Program Files\Go\bin\go.exe" test -v -tags go_tarantool_ssl_disable -vet=off -count=1 -run ^TestTarantool$ shortlink/internal/adapter/db
    echo -e "\n>>_TarantoolTest_<<"
    go test -v -tags go_tarantool_ssl_disable -vet=off -count=1 -run ^TestTarantool$ shortlink/internal/adapter/db
    if [[ $? -gt 0 ]]; then checksBreaked; fi   
}

function build {
    echo -e "\n>>_Build_<<"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags go_tarantool_ssl_disable -ldflags='-w -s -extldflags "-static"' -a -o ./bin/shortlink ./cmd/main.go
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function run {
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

function metrics-graph {                                    # https://www.yworks.com/yed-live/
    echo -e "\n>>_MetricsGraph_<<"
    image_packages . ./asset/metrics.graphml
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function loc {
    echo -e "\n>>_LinesOfCode_<<"
    FPATH=./script/metrics/sl.loc
    gcloc . > $FPATH
    cat $FPATH | grep Language
    if [[ $? -gt 0 ]]; then checksBreaked; fi
    cat $FPATH | grep Golang
    cat $FPATH | grep Bash
}

function cycl {
    echo -e "\n>>_CyclomaticComplexity_<<"
    FPATH=./script/metrics/sl.cycl
    gocyclo -avg . > $FPATH
    cat $FPATH | grep Average
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function cogn {
    echo -e "\n>>_CognitiveComplexity_<<"
    FPATH=./script/metrics/sl.cogn
    gocognit -avg . > $FPATH
    cat $FPATH | grep Average
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function testcov {
    echo -e "\n>>_TestCoverage_<<"
    FTEMP=./script/metrics/coverage.test
    FOUT=./script/metrics/sl.cov
    go test -short -tags go_tarantool_ssl_disable -vet=off -count=1 -coverprofile=$FTEMP ./...
    go tool cover -func=$FTEMP > $FOUT
    cat $FOUT | grep total | awk '{print "TOTAL:", $3}'
    if [[ $? -gt 0 ]]; then checksBreaked; fi
}

function complex {
    echo -e "\n>>_ComplexityMetrics_<<"
    FPATH=./script/metrics/sl.compl
    go vet -vettool=$(which complexity) -loc -cyclo -halst -maint ./... 2>$FPATH    
    cat $FPATH | grep locSum | awk 'BEGIN{sum=0}{sum+=$2}END{print "TOTAL LoC:", sum}'
    if [[ $? -gt 0 ]]; then checksBreaked; fi
    cat $FPATH | grep cycloAvg | awk 'BEGIN{sum=0;div=0}{sum+=$3;div+=$4}END{if (div==0) div=1; print "TOTAL cycloAvg:", sum/div}'
    cat $FPATH | grep halstVolAvg | awk 'BEGIN{sum=0;div=0}{sum+=$3;div+=$4}END{if (div==0) div=1; print "TOTAL halstVolAvg:", sum/div}'
    cat $FPATH | grep halstDiffAvg | awk 'BEGIN{sum=0;div=0}{sum+=$3;div+=$4}END{if (div==0) div=1; print "TOTAL halstDiffAvg:", sum/div}'
    cat $FPATH | grep maintAvg | awk 'BEGIN{sum=0;div=0}{sum+=$3;div+=$4}END{if (div==0) div=1; print "TOTAL maintAvg:", sum/div}'    
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
    "ttTest")
        intergTTtest
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
        testcov
        loc
        cycl
        cogn        
        complex
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
        # docker run -d --name SLsrv --user 10001 -p 8080:8080/tcp kselnaag/shortlink
        # docker run -d --name SLpg -p 5432:5432 -e POSTGRES_DB=shortlink -e POSTGRES_USER=login -e POSTGRES_PASSWORD=password postgres:16.0-alpine3.18
        # docker run -d --name SLmg -p 27017:27017 -e MONGO_INITDB_DATABASE=shortlink -e MONGO_INITDB_ROOT_USERNAME=login -e MONGO_INITDB_ROOT_PASSWORD=password mongo:7.0.2 
        # docker run -d --name SLrd -p 6378:6379 -e REDIS_ARGS="--requirepass password" redis:7.2.1-alpine3.18
        # docker run -d --name SLtt -p 3301:3301 -e TARANTOOL_USER_NAME=login -e TARANTOOL_USER_PASSWORD=password tarantool/tarantool:2.10.8-gc64-amd64

        # docker buildx build --platform linux/amd64 --no-cache --build-arg="TNT_VER=2.11.1" --build-arg="NPROC=1" --build-arg="ENABLE_BUNDLED_LIBYAML=ON" --build-arg="LUAJIT_DISABLE_SYSPROF=OFF" --build-arg="GC64=ON" --build-arg="ROCKS_INSTALLER=luarocks" -f ./dockerfiles/alpine_3.15 -t kselnaag/tarantool:2.11.1 . --load
        # go test -v -tags go_tarantool_ssl_disable -vet=off -count=1 -race -timeout 30s -run ^TestTarantool$ shortlink/internal/adapter/db
        # pprof
        # go tool pprof shortlink http://localhost:8080/debug/pprof/profile
        # ab -n 100000 -c10 http://localhost:8080/
        ;;
    *)
        info
        ;;
    esac
fi

echo -e "\n>>_Successfull_<<\n"
