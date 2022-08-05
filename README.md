# Upvote Challenge

# Context
This project was developed to fulfill the technical requests from a job interview.
Proposed challenge:
<br>
"
The Technical Challenge consists of creating an API with Golang using gRPC with stream pipes that exposes an upvote service endpoints.
Technical requirements:
- Keep the code in Github
API:
- The API must guarantee the typing of user inputs. If an input is expected as a string, it can only be received as a string.
- The structs used with your mongo model should support Marshal/Unmarshal with bson, json and struct
- The API should contain unit test of methods it uses
<br>
Extra:
- Deliver the whole solution running in some free cloud service
"

## Used stacks

> MongoDB, Golang, gRPC, Cobra, Docker, protobuf

## Dependencies

- Go
- docker
- docker-compose

## Running database and Server

```
docker-compose up -d
```

## Running Client

```
docker-compose run client
/client [command]
```

## Running Tests
:warning: **Tests just works if the server and db is running up**
  ```
  cd go/server/
  go test -v
  ```