<p align="left">
	<img src="https://img.shields.io/github/languages/code-size/kselnaag/shortlink?style=plastic" title="src files size" alt="src files size">
	<img src="https://img.shields.io/github/repo-size/kselnaag/shortlink?style=plastic" title="repo size" alt="repo size">
	<a href="https://github.com/kselnaag/shortlink/blob/master/LICENSE" title="LICENSE"><img src="https://img.shields.io/github/license/kselnaag/shortlink?style=plastic" alt="license"></a>
	<a href="https://github.com/kselnaag/shortlink/actions" title="Workflows"><img src="https://img.shields.io/github/actions/workflow/status/kselnaag/shortlink/go.yml?branch=master&style=plastic" alt="tests checks"></a>
</p>
<p align="left">
  <img src="https://img.shields.io/static/v1?label=test&message=Project&color=ffa757&style=plastic" alt="test Project">
	<img src="https://img.shields.io/static/v1?label=build%20by&message=Go&color=ffa757&style=plastic" alt="build by Go">
	<img src="https://img.shields.io/static/v1?label=make%20with&message=Fun&color=ffa757&style=plastic" alt="make with Fun">
</p>

### **SHORTLINK** 📏 Let us make your links shorter in easy way !
----

## **🧾Description**
This is a test project to generate the short link from the long link you already have. We want to be able to:
<img style="margin-right: 60px;" align="right" width="40%" alt="#POWERGOPHERS" src="./asset/gogogophers.png"/>

- get the short link from the long link
- save the result to database
- redirect from the short link to the long link destination
- get simple UI as HTML page
- check if the long link HTTP valid

## **📊Analysis**
We choose Monolith as system arch pattern and Rich Domain Model as software arch pattern. Let us look at some architect points:
<img style="margin-top: 0px; margin-right: 100px;" align="right" width="29%" alt="#ArchPic" src="./asset/arch.png"/>

- `Domain Adapters`
  - HTTP transport
  - SQL database
  - JSON logger
  - file + env config
- `Use Cases`
  - get healthcheck
  - get html UI
  - redirect from short link to long link 
  - search the short link if you have a long link
  - search the long link if you have a short link
  - get ALL link pairs presented in db
  - check if long link HTTP available
- `Domain Rules`
  - compute short link from long link
  - unite short link and long link
  - check if pair is valid
- `Domain Models`
  - link pair (string, string)

## **💡Solution notes**
<img style="margin-right: 100px; transform: rotate(03.7deg);" align="right" width="14%" alt="#Prod" src="./asset/production.png"/>

- clean arch + DDD aproach
- standart go project layout (more or less)
- pre-commit hooks and github actions (CI) + podman-compose (CD) + minikube (prod🙃)
- tests with mocks included
- tarantool migrations and TTL records included

## **🛠️Libs and tools**
<img style="margin-right: 0px;" align="right" width="50%" alt="#DOMAIN" src="./asset/domain.png"/>

- `Libs (github.com)`
  - caarlos0/env v3.5.0
  - joho/godotenv v1.5.1
  - rs/zerolog v1.29.0
  - gin-gonic/gin v1.9.0
  - gin-contrib/static v0.0.1
  - valyala/fasthttp v1.45.0
  - gofiber/fiber/v2 v2.42.0
  - stretchr/testify v1.8.2
  - gavv/httpexpect v2.15.0
- `Tools`
  - golangci-lint
  - curl 
  - podman + podman-compose
  - minikube

## **⚙️HowTo**
<img style="margin-right: 0px;" align="right" width="30%" alt="#CICD" src="./asset/cicd.png"/>

- check if `podman` and `podman-compose` has been installed
- clone the project
- run everything with `./script/run.sh`
- go to `http://localhost:8080` in your browser and try it

## **🦋The beauty is like this and nothing more**

----
### **🔗LINKS**
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

