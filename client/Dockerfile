FROM golang:1.18 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client main.go

FROM alpine

WORKDIR /

COPY --from=build /app/client /
