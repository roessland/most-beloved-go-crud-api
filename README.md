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
- PostgreSQL lib: pgx. Everyone likes pgx.
- DB migrations: Atlas. I've used golang-migrate before, but Atlas promises a
nice declarative way to automatically create migrations.

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

## Feedback to FL0
- Web UI could be snappier. A lot of delays when clicking around.
- Make it more clear how to cache Docker build steps. My build takes 1 minute
longer than necessary since dependences are downloaded again every time.
- "DB name same as application name" use-case is very common. Show error
message before I submit the form.
- Provide a nice way to enable running migrations as part of deployment,
instead of during process startup, to enable 

## Backlog
- Add rest of the handlers
- Add sqlc
- Add migrations and SQL
- Add pg DB in FL0.

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
- Always run Gin in release mode when running as a Docker container. Go does
not need to be compiled for release.
- Put go mod download before `COPY . .` to avoid re-downloading dependencies if
go.mod and go.sum are unchanged.

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
ENV GIN_MODE=release
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


### Move routes to a separate package
To keep it clean when adding more handlers. Remember to rename to capital case
so it's exported from the handlers package.
```golang
// handlers/handlers.go
package handlers
func Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
```

### Run Gin in release mode by default
Add `ENV GIN_MODE=release` to the last stage in the Dockerfile.
Run `docker build -t most-beloved-go-crud-api .; docker run -e PORT=80 -it --rm
-p 8080:80 -e PORT=80 most-beloved-go-crud-api` and then `curl localhost:8080`
to see output. By default, nothing is logged at startup during release mode.

Verify release mode is working by deploying and looking at logs:

```shell
git add handlers/
git commit -am "make handlers package"
```

Now, in FL0 -> Environment variables, add `GIN_MODE=debug`, click Save, then
redeploy the latest deployment (under sandwich menu). View the logs. Verify
that we are now back in debug mode.

### Add Create endpoint (no DB)
We want to parse the POST body JSON to a struct. In web framework-speak this is
called model binding. [Gin docs for model binding and
validation](https://gin-gonic.com/docs/examples/binding-and-validation/).

Add a new Create route. Remember to register it.
Add input and result structs for model binding.

```golang
  // main.go
	r.POST("/", handlers.Create)
```

```golang
package handlers

type CreateQuoteInput struct {
	Book  string `json:"book" binding:"required"`
	Quote string `json:"quote" binding:"required"`
}

type CreateQuoteResult struct {
	UUID  string `json:"uuid"`
	Book  string `json:"book"`
	Quote string `json:"quote"`
}

func Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// Create creates a new quote
func Create(c *gin.Context) {
	var input CreateQuoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := createQuote(input)

	c.JSON(http.StatusCreated, result)
}

func createQuote(request CreateQuoteInput) CreateQuoteResult {
	// TODO:: Save to DB
	return CreateQuoteResult{
		UUID:  "1234-TODO-1234",
		Book:  request.Book,
		Quote: request.Quote,
	}
}
```

Test it:
```shell
$ go run main.go
$ curl -X POST localhost:8080 -d {}
{"error":"Key: 'CreateQuoteInput.Book' Error:Field validation for 'Book' failed on the 'required' tag\nKey: 'CreateQuoteInput.Quote' Error:Field validation for 'Quote' failed on the 'required' tag"}%

$ curl -X POST localhost:8080 -d '{"book": "LOTR", "quote": "Precioussss"}'
{"uuid":"1234-TODO-1234","book":"LOTR","quote":"Precioussss"}%
```

### Add a database and connect to it
Open your FL0 project and add a new PostgreSQL database. Name it the same as
your project. Realize that FL0 prohibits this. In anger, append `-db`. Success.

Take a look at the connection settings. There is no need to copy them into the
environment variables by hand. Instead, go the container environment variables
and use the Import database credentials button. Save.

This will add a bunch of variables:
- DATABASE_URL
- PGHOST
- PGDATABASE
- PGPORT
- PGUSER
- PGPASSWORD
- PGSSLMODE (=`require`)

Add pgx dependency.
```shell
go get github.com/jackc/pgx/v5
```

Connect to the database and verify that connection is working:
```golang
func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	dbHealthCheck(conn)

	r := setupRouter()
	r.Run()
}

func dbHealthCheck(conn *pgx.Conn) {
	healthCheckCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result int
	err := conn.QueryRow(healthCheckCtx, "SELECT 1 + 1;").Scan(&result)
	if err != nil || result != 2 {
		panic(err)
	}
}
```

Copy the database URL and export it to a local shell variable, so we can run our local instance against it.
```shell
export DATABASE_URL='postgres://fl0user:password@my-instance.eu-central-1.aws.neon.fl0.io:5432/most-beloved-go-crud-api-db?sslmode=require'
go run main.go
```

Verify that it also works when deployed to FL0.
```shell
git commit -am "connect to db"
git push
```

### Setup sqlc and write some queries

Install the CLI tool. https://docs.sqlc.dev/en/stable/overview/install.html

```shell
brew install sqlc # MacOS
```

Add sqlc config file `sqlc.yaml`

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "db"
        sql_package: "pgx/v5"
        out: "db"
```

Follow the Atlas guide for sqlc to create migrations:
https://atlasgo.io/guides/frameworks/sqlc-versioned


Add some queries:

```sql
-- query.sql
-- name: CreateQuote :one
INSERT INTO quotes (id, book, quote, inserted_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
```

Add the initial migration:
```sql
-- schema.sql
CREATE TABLE IF NOT EXISTS  quotes (
  id UUID PRIMARY KEY,
  book varchar NOT NULL,
  quote TEXT NOT NULL,
  inserted_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (book, quote)
);
```

Generate a go-package named `db`.

```shell
sqlc generate
```

#### Sidenote: Migration strategy
Running migrations on application startup is an anti-pattern, since the
migrations can take a long time, and health checks will fail for that duration.
Migrations should instead be run in the deployment pipeline, before the new
application version is deployed. This does not seem to be supported by FL0 in
the free plan, since it's only possible to run one container, so I will run
migrations on startup regardless, since it's just a simple demo application.

### Setup migrations with Atlas
Atlas provides a declarative approach where we define the desired DB state, and
Atlas comes up with a migration we can apply.

The target state can be described with HCL (HashiCorp Configuration Language).
Considering the recent HashiCorp licensing issues I'm hesitant to use HCL, but
HCL's licensing seems to remain MPL (Mozilla Public License).

```shell
brew install ariga/tap/atlas
```

Run the initial migration:
```shell
atlas schema apply -u "$DATABASE_URL" --to file://schema.sql --dev-url "$DATABASE_URL" -s public
```

Note that Atlas requires a dev-database to parse SQL code, and we simply reused
the dev-database for this. Atlas also has a method to spin up a local DB using
a Docker container, but this didn't work for me. (`Error: cannot diff a schema
with a database connection: "" <> "public"`).

Login with `psql` and verify that the changes were applied:

```shell
$ psql "$DATABASE_URL"
most-beloved-go-crud-api-db=> \d quotes
                          Table "public.quotes"
   Column    |           Type           | Collation | Nullable | Default
-------------+--------------------------+-----------+----------+---------
 id          | uuid                     |           | not null |
 book        | character varying        |           | not null |
 quote       | text                     |           | not null |
 inserted_at | timestamp with time zone |           | not null |
 updated_at  | timestamp with time zone |           | not null |
Indexes:
    "quotes_pkey" PRIMARY KEY, btree (id)
    "quotes_book_quote_key" UNIQUE, btree (book, quote)
```
Success!

### Inject DB into handlers, UUID shenanigans
Best practice is to inject dependencies. This can can be done with a closure,
or by making the handler functions into struct methods.

Use a more ergonomic UUID type:
```shell
go get github.com/google/uuid
```

sqlc.yaml. Remember to `sqlc generate` after this.
```sqlc
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "db"
        sql_package: "pgx/v5"
        out: "db"
        overrides:
        - db_type: "uuid"
          go_type: "github.com/google/uuid.UUID"
        - db_type: "uuid"
          go_type: "github.com/google/uuid.NullUUID"
```

Inject DB into handlers:

handlers.go
```golang
import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/roessland/most-beloved-go-crud-api/db"
)

//...

// Quotes holds resources used by handlers
type Quotes struct {
	Queries *db.Queries
}

// Create creates a new quote
func (quotes *Quotes) Create(c *gin.Context) {
	var input CreateQuoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := quotes.create(c, input)
	if err != nil {
		// TODO: Don't return the actual error in release mode
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// create creates a new quote
func (quotes *Quotes) create(ctx context.Context, input CreateQuoteInput) (*CreateQuoteResult, error) {
	dbParams := db.CreateQuoteParams{
		ID:         uuid.New(),
		Book:       input.Book,
		Quote:      input.Quote,
		InsertedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	dbQuote, err := quotes.Queries.CreateQuote(ctx, dbParams)
	if err != nil {
		return nil, err
	}
	return &CreateQuoteResult{
		UUID:  dbQuote.ID.String(),
		Book:  dbQuote.Book,
		Quote: dbQuote.Quote,
	}, nil
}
```

main.go:
```golang
// ...
func setupRouter(queries *db.Queries) *gin.Engine {
	r := gin.Default()

	quotes := &handlers.Quotes{
		Queries: queries,
	}

	r.GET("/", handlers.Health)
	r.POST("/", quotes.Create)

	return r
}

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)
	queries := db.New(conn)

	dbHealthCheck(conn)

	r := setupRouter(queries)
	r.Run()
}
// ...
```
