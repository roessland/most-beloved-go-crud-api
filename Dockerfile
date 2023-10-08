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
