[![Test](https://github.com/kaduartur/planet/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/kaduartur/planet/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/kaduartur/planet/branch/main/graph/badge.svg?token=XJU7N08559)](https://codecov.io/gh/kaduartur/planet)
[![CodeFactor](https://www.codefactor.io/repository/github/kaduartur/planet/badge?s=c776547d6f360ee54ba5ccf8fda1580978c2648d)](https://www.codefactor.io/repository/github/kaduartur/planet)


# 🌎 The planets api 🪐

The planets API was made to manage planets, here you can create, find, list, and delete planets.

<img width="45%" src="https://github.com/ashleymcnamara/gophers/blob/master/GOPHER_STAR_WARS.png"/>

#### Application architecture
This project was divided into three applications:
- command-api: Responsible only for writing to the database.
- processor: Responsible for processing the planets created async.
- query-api: Responsible for reading the planets.

![planet_architecture](https://user-images.githubusercontent.com/17505818/117897434-2ecf2580-b299-11eb-9b23-95e1bdce1aed.png)


#### Requirements to run:
* 🐳 Docker compose
* <img width="15px" src="https://golang.org/favicon.ico"/> Go
* <img width="15px" src="https://www.postman.com/favicon-32x32.png?v=6fa10b9ee2b6e5dcec30e5027a14e7a4"/> Postman
* :ox: GNU Make

### How to run?
- Run `docker-compose up` command inside the project folder.
- Run `make run-command-api` to start the command server.
- Run `make run-processor` to start the processor server.
- Run `make run-query-api` to start the query server.

The [postman](https://documenter.getpostman.com/view/2956294/TzRUA718
) link to see the API documentation.

### How to run tests?
- Run `docker-compose up` command inside the project folder.
- Run `make run-test`

