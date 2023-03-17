<p align="left">
  <img src="https://img.shields.io/static/v1?label=test&message=Project&color=ffa757&style=plastic" alt="test Project">
	<img src="https://img.shields.io/static/v1?label=build%20by&message=Go&color=ffa757&style=plastic" alt="build by Go">
	<img src="https://img.shields.io/static/v1?label=make%20with&message=Fun&color=ffa757&style=plastic" alt="make with Fun">
</p>

### **SHORTLINK** 📏 Let us make your links shorter in easy way !
----

## **🧾Description**
This is a test project to generate the short link from the long link you already have. We want to be able to:
<img style="margin-top: 20px; margin-right: 60px;" align="right" width="40%" alt="#POWERGOPHERS" src="./asset/gogogophers.png"/>

- get the short link from the long link
- save the result to database
- redirect from the short link to the long link destination
- get simple UI as HTML page
- get all link pairs we already have
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
  - link pair

## **💡Solution notes**
<img style="margin-right: 100px; transform: rotate(03.7deg);" align="right" width="14%" alt="#Prod" src="./asset/production.png"/>

- DDD aproach
- standard go project layout (more or less)
- pre-commit hooks and github actions (CI) + podman-compose (tests) + minikube (prod 🙃)
- tarantool migrations and TTL records included
- tests (with mocks and thunderclient) included

## **🛠️Libs and tools**
<img style="margin-right: 0px;" align="right" width="30%" alt="#CICD" src="./asset/cicd.png"/>

- libs
- podman + podman compose
- minikube
- d

## **⚙️HowTo**

- check if `podman` and `podman-compose` has been installed
- clone the project
- run everything with `./script/run.sh`
- go to `localhost:8080` in your browser and try it

## **🦋A picture**

----
### **🔗LINKS**
| [gin](https://github.com/gin-gonic/gin "https://github.com/gin-gonic/gin")
| [gin docs](https://gin-gonic.com/docs/ "https://gin-gonic.com/docs/")
| [fiber](https://github.com/gofiber/fiber "https://github.com/gofiber/fiber")
| [fiber docs](https://docs.gofiber.io "https://docs.gofiber.io")
|

| [tarantool](https://github.com/tarantool/tarantool "https://github.com/tarantool/tarantool")
| [tarantool docs](https://www.tarantool.io/ru/doc/ "https://www.tarantool.io/ru/doc")
| 