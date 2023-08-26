# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/todo
RUN CGO_ENABLED=0 GOOS=linux go build -o main

WORKDIR /app/cmd/health
RUN CGO_ENABLED=0 GOOS=linux go build -o health

FROM alpine:latest AS build-release-stage
WORKDIR /app

COPY --from=build-stage /app/cmd/todo/main .
COPY --from=build-stage /app/cmd/health/health .
RUN chmod +x ./main
RUN chmod +x ./health
EXPOSE 9090
ENTRYPOINT ["./main"]
