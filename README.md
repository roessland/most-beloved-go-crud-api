# Most beloved Quotes API
Simple service that serves quotes. Made using the most beloved technology
listed in Stack Overflow Developer Survey 2023. Inspired by [Dreams of Code:
Building the most loved CRUD app](https://www.youtube.com/watch?v=3oN70MzSDfY),
but with Go instead of Rust.

## Tech
- Go: Not the most beloved in the survey, but my favourite, and close enough.
- PostgreSQL: Most beloved database
- Not-an-ORM: ( sqlc )[https://sqlc.dev/]: Generates type-safe code from SQL queries.
- Cloud: Using [FL0](https://www.fl0.com/)'s free tier to try it out.
- Docker: Since it's so popular
- Web framework: Gin. Runner ups: Chi, standard library + gorilla/mux. But I've
never used Gin, so I'm trying it out. Seems like the only "real" framework for
Go APIs, as many people prefer to just use the standard library.

## Features
- POST /quotes - Add quote
- GET /quotes - Gets all quotes
- PUT /quotes/{id} - Replaces quote
- DELETE /quotes/{id} - Deletes quote

Data model:
```json
{
  "book": "The Hobbit",
  "quote": "My precious...",
}
```

## Backlog
Switch to release mode when running in cloud.

## Dev log
Make new project:
```shell
mkdir most-beloved-go-crud-api
cd most-beloved-go-crud-api
go mod init github.com/roessland/most-beloved-go-crud-api
```

Make main executable:
```shell
touch main.go
```

Make GitHub repository.
```shell
# use web UI
# add MIT license, no readme, Go .gitignore
```

Initial commit:
```shell
git init
git branch -m master main
git add .
git commit -m "add readme and go.mod"
git remote add origin git@github.com:roessland/most-beloved-go-crud-api.git
git pull origin main
git push
```

Add Gin dependency:
```shell
go get -u github.com/gin-gonic/gin
````

Add web server [boilerplate](https://raw.githubusercontent.com/gin-gonic/examples/master/basic/main.go) to `main.go`:
```golang
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen on PORT env var, or on 3000 if not set
	r.Run()
}
```

Run it and test it:
```shell
$ go run main.go
[GIN-debug] Listening and serving HTTP on 0.0.0.0:3000

$ curl localhost:3000
OK
```

Sign up for FL0 (I promise, this is not an ad): 
- https://www.fl0.com/
- Continue with GitHub
- Workspace "roessland"
- Project "most-beloved-go-crud-api"
- Deploy an existing GitHub repo
- Authorize FL0 connector to only this repo
- Take note that FL0 requires the app to listen on 8080, not 3000 as we have defined. It provides a PORT environment variable that we can use.
- Advanced options -> Dockerfile -> Keep everything default/unchanged
- Deploy fails since we have no Dockerfile

### Add Dockerfile
Follow docs: https://docs.fl0.com/docs/builds/dockerfile/go

Changes:
- Go 1.21 instead of Go 1.19 as in docs
- Add gcompat to enable running glibc programs.

```dockerfile
ARG APP_NAME=most-beloved-go-crud-api

# Build stage
FROM golang:1.21 as build
ARG APP_NAME
ENV APP_NAME=$APP_NAME
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /$APP_NAME

# Production stage
FROM alpine:latest as production
RUN apk add gcompat
ARG APP_NAME
ENV APP_NAME=$APP_NAME
WORKDIR /root/
COPY --from=build /$APP_NAME ./
CMD ./$APP_NAME
```

Usage:
```shell
docker build -t most-beloved-go-crud-api .
```

Run it:
```shell
docker run -p 3000:3080 -e PORT=3000 most-beloved-go-crud-api
```

View it on http://localhost:3000

### Run it in the cloud
```shell
git push
```
Follow the build at FL0's web UI: https://app.fl0.com/roessland/most-beloved-go-crud-api/dev/most-beloved-go-crud-api/deployments

Open the URL near the top to see it live: https://griffith-koala-cded.1.ie-1.fl0.io

