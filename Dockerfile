# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY index.gohtml ./

RUN go build -o /unixtime

# Use a distroless image to save space
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /unixtime /unixtime
COPY --from=build-stage /app/index.gohtml /index.gohtml

USER nonroot:nonroot

EXPOSE 9000

ENTRYPOINT ["/unixtime"]