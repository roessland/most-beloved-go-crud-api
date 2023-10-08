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
- Web framework: Gin. Runner ups: Chi, standard library + gorilla/mux. But I've never used Gin, so I'm trying it out.

## Features
- POST /quotes - Add quote
- GET /quotes - Gets all quotes
- PUT /quotes/{id} - Replaces quote
- DELETE /quotes/{id} - Deletes quote

Data model:
```go
{
  "book": "The Hobbit",
  "quote": "My precious...",
}
```

## Dev log
Make new project:
```shell
mkdir most-beloved-go-crud-api
cd most-beloved-go-crud-api
go mod init github.com/roessland/most-beloved-go-crud-api
```

Make main executable:
```
touch main.go
```

Make GitHub repository.
```
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
