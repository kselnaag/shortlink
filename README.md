<p align="left">
  <img src="https://img.shields.io/static/v1?label=test&message=Project&color=ffa757&style=plastic" alt="test Project">
	<img src="https://img.shields.io/static/v1?label=build%20by&message=Go&color=ffa757&style=plastic" alt="build by Go">
	<img src="https://img.shields.io/static/v1?label=make%20with&message=Fun&color=ffa757&style=plastic" alt="make with Fun">
	<a href="https://github.com/kselnaag/shortlink/blob/master/LICENSE" title="LICENSE"><img src="https://img.shields.io/github/license/kselnaag/shortlink?style=plastic" alt="license"></a>
	<a href="https://github.com/kselnaag/shortlink/actions" title="Workflows"><img src="https://img.shields.io/github/actions/workflow/status/kselnaag/shortlink/go.yml?branch=master&style=plastic" alt="tests checks"></a>
</p>

### **SHORTLINK** 📏 Let us make your links shorter in easy way !
----

## **📜 Description**
This is a test project to generate the short link from the long link you already have. We want to be able to:
<img style="margin-right: 60px;" align="right" width="40%" title ="Go go, Power Gophers !" alt="#POWERGOPHERS" src="./asset/gogogophers.png"/>

- get the short link from the long link
- save the result to database
- redirect from the short link to the long link destination
- get simple UI as HTML page
- check if the long link HTTP valid

## **🔬 Analysis**
We choose Monolith as system arch pattern and Rich Domain Model as software arch pattern. Let us look at some architect points:
<img style="margin-top: 0px; margin-right: 100px;" align="right" width="29%" title ="Architecture" alt="#ARCPIC" src="./asset/arch.png"/>

- `Adapters`
  - HTTP transport
  - SQL or NoSQL database
  - JSON logger
  - file + env config
- `Interfaces`
  - http
  - db
  - log
  - service
- `Controllers`
  - http
  - db
- `Services`
  - get healthcheck
  - get html UI
  - redirect from short link to long link 
  - search the short link if you have a long link
  - search the long link if you have a short link
  - get ALL link pairs presented in db
  - check if long link HTTP available
- `Models`
  - link pair (string, string)
  - compute short link from long link
  - unite short link and long link
  - check if pair is valid

## **💡 Solution notes**
<img style="margin-right: 100px;" align="right" width="14%" title ="Prod 🙃" alt="#PROD" src="./asset/production.png"/>

- layering and hexagonal arch
- standart go project layout (more or less)
- pre-commit hooks and github actions (CI) + docker-compose (CD) + minikube (prod🙃)
- integration and unit tests with mocks included
- HTTPfast or HTTPnet server adapters
- logZero, logSlog or logFprintf logger adapters
- postgreSQL, redis, mongoDB or tarantool db adapters
- OpenAPI Swagger docs `http://localhost:8080/oapi`

## **🛠️ Libs and tools**
<img style="margin-right: 0px;" align="right" width="50%" title ="Domain" alt="#DOMAIN" src="./asset/domain.png"/>

- `Libs (https://github.com/)`
  - caarlos0/env v3.5.0
  - joho/godotenv v1.5.1
  - rs/zerolog v1.29.0
  - gin-gonic/gin v1.9.0
  - gin-contrib/static v0.0.1
  - valyala/fasthttp v1.45.0
  - gofiber/fiber v2.42.0
  - stretchr/testify v1.8.2
  - gavv/httpexpect v2.15.0
  - jackc/pgx v5.4.3
- `Tools`
  - golang v1.21
  - golangci-lint v1.51.2
  - curl v7.68.0
  - docker v24.0.6 + docker-compose v1.25.0
  - minikube

## **📊 Metrics**
<img style="margin-right: 0px;" align="right" width="60%" title="Package graph" alt="#METRICS" src="./asset/metricsgraph.png"/>
<p align="left">
	<img src="https://img.shields.io/github/languages/code-size/kselnaag/shortlink?style=plastic" title="src files size" alt="src files size">
	<img src="https://img.shields.io/github/repo-size/kselnaag/shortlink?style=plastic" title="repo size" alt="repo size">
</p>

- Code size + Repository size
- Package graph (image_packages)
- Unit test coverage (go tool cover)
- Lines of code (gcloc, complexity)
- Cyclomatic complexity (gocyclo, complexity)
- Cognitive Complexity (gocognit)
- Halstead complexity (complexity)
- Maintainability index (complexity)
- `Metrics repos (https://github.com/)`
  - ManyakRus/image_packages
  - JoaoDanielRufino/gcloc
  - fzipp/gocyclo
  - uudashr/gocognit
  - kselnaag/go-complexity-analysis (fork)

```
kselnaag:shortlink$ ./script/cicd.sh

CI/CD COMMANDS: style lint test build run start check check-no-lint 
                docker-gobuilder docker-build docker-run docker-up compose-up
                metrics metrics-graph
EXAMLPE:        ./script/cicd.sh build

kselnaag:shortlink$ ./script/cicd.sh metrics

>>_TestCoverage_<<
?       shortlink/cmd                   [no test files]
?       shortlink/internal              [no test files]
?       shortlink/internal/adapter/cfg  [no test files]
ok      shortlink/internal/adapter/db   0.026s  coverage: 35.4% of statements
?       shortlink/internal/apptype      [no test files]
?       shortlink/internal/control      [no test files]
?       shortlink/web                   [no test files]
ok      shortlink/internal/adapter/http 0.094s  coverage: 38.6% of statements
ok      shortlink/internal/adapter/log  0.022s  coverage: 77.0% of statements
ok      shortlink/internal/model        0.031s  coverage: 100.0% of statements
ok      shortlink/internal/service      0.030s  coverage: 69.8% of statements
TOTAL: 64.2%

>>_LinesOfCode_<<
  Language   | Files | Lines | Blank lines | Comments | Code lines
  Golang     |    33 |  1698 |         205 |       52 |       1441
  Bash       |     1 |   294 |          32 |        6 |        256

>>_CyclomaticComplexity_<<
Average: 2.26

>>_CognitiveComplexity_<<
Average: 1.38

>>_ComplexityMetrics_<<
TOTAL LoC: 1138
TOTAL cycloAvg: 2.31034
TOTAL halstVolAvg: 323.898
TOTAL halstDiffAvg: 18.3464
TOTAL maintAvg: 64.5287

>>_Successfull_<<
```

More details in: ./script/metrics/

## **⚙️ HowTo**
<img style="margin-right: 0px;" align="right" width="30%" title="CI/CD workflow" alt="#CICD" src="./asset/cicd.png"/>

- check if `docker` and `docker-compose` has been installed
- clone the project `git clone https://github.com/kselnaag/shortlink`
- dive into a folder `cd shortlink`
- run everything with `./script/cicd.sh compose-up`
- go to `http://localhost:8080` in your browser and try it

## **🦋 The beauty is like this and nothing more**

----
### **🔗 LINKS**
| [gin](https://github.com/gin-gonic/gin "https://github.com/gin-gonic/gin")
| [gin docs](https://gin-gonic.com/docs/ "https://gin-gonic.com/docs/")
|             
| [zerolog](https://github.com/rs/zerolog "https://github.com/rs/zerolog")
|                    
| [tarantool](https://hub.docker.com/r/tarantool/tarantool "https://hub.docker.com/r/tarantool/tarantool")
| [tarantool docs](https://www.tarantool.io/ru/doc/ "https://www.tarantool.io/ru/doc")
| 

| [fiber](https://github.com/gofiber/fiber "https://github.com/gofiber/fiber")
| [fiber docs](https://docs.gofiber.io "https://docs.gofiber.io")
|         
| [godotenv](github.com/joho/godotenv "github.com/joho/godotenv")
| [env](https://github.com/caarlos0/env "https://github.com/caarlos0/env")
|        
| [redis](https://hub.docker.com/_/redis "https://hub.docker.com/_/redis")
| [redis docs](https://redis.io/docs/ "https://redis.io/docs/")
|
