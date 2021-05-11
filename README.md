# 🌎 The planets api 🪐

The planets API was made to manage planets, here you can create, find, list, and delete planets.

<img width="45%" src="https://github.com/ashleymcnamara/gophers/blob/master/GOPHER_STAR_WARS.png"/>

#### Applications descriptions
This project was divided into three applications:
- command-API: Responsible only for writing to the database.
- processor: Responsible for processing the planets created async.
- query-API: Responsible for reading the planets.

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

