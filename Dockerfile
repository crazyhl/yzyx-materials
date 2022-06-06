# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY . ./
RUN go mod download


RUN go build -o /server ./cmd/server/main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /server /server
COPY --from=build /app/env.prod /env

EXPOSE 8080

# USER nonroot:nonroot

ENTRYPOINT ["/server"]