FROM golang:1.18 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o post_server ./server/main.go

FROM scratch

WORKDIR /

COPY --from=build /app/post_server /

EXPOSE 3001

CMD [ "/post_server" ]