# syntax=docker/dockerfile:1

FROM golang:1.20-alpine as build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
WORKDIR /app/cmd/todo
RUN CGO_ENABLED=0 GOOS=linux go build -o main
RUN chmod +x ./main
EXPOSE 9090

FROM alpine:latest AS build-release-stage
WORKDIR /app

COPY --from=build-stage /app/cmd/todo/main .
RUN chmod +x ./main
EXPOSE 9090
CMD ["./main"]
