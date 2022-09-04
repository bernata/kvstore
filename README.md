[![Go](https://github.com/bernata/kvstore/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/bernata/kvstore/actions/workflows/go.yml)
![Coverage](https://img.shields.io/badge/Coverage-92.9%25-brightgreen)

# kvstore
In memory key value store backed by a simple map protected by a read-write lock.

## Description
In memory key value store backed by a simple map protected by a read-write lock.
This repository consists of 2 components:
- kvservice: in memory key-value service with http endpoints for: write, read, delete
- kvclient: convenience client to invoke http endpoints from the command line

Currently, there is no authentication to the kvservice so all keys/values are accessible
to everyone; and can be overwritten by anyone.

NOTE: keys are logged; but values are never logged.

## Requirements
- golang 1.18
- golangci-lint 1.46.2+
- gnumake 3.81+
- docker 20.10+

## Getting Started
```
git clone https://github.com/bernata/kvstore.git
cd kvstore
make all
```

In one terminal type:
```
make run
```

In another terminal type:
```shell
./bin/kvclient write --key k1 --value v1
./bin/kvclient get --key k1
./bin/kvclient delete --key k1
```

## Build
`make all`

Builds both kvclient and kvservice in `bin/`

## Test
`make test`

Runs unit tests for both kvclient and kvservice.

## Run
`make run`

Runs kvservice locally on port 8282

## Docker
`make docker` - builds a centos image of kvservice

`make dockerrun` - runs the docker image with port mapping 8282:8282

It can be convenient to run the docker container and use the client to test against the service.
```
make docker
make dockerrun
make build_client
./bin/kvclient write --key "mykey" --value "myvalue"
```

## Examples
```shell
  > make run   # runs kvservice server at port 8282 in foreground
  
  > ./bin/kvclient write --key "mykey" --value "myvalue"
    {}
    
  > ./bin/kvclient get --key "mykey"
    {"key": "mykey", "value": "myvalue"}
    
  > ./bin/kvclient write --key "mykey" --value "myvalue2"
    {}
    
  > ./bin/kvclient get --key "mykey"
    {"key": "mykey", "value": "myvalue2"}
    
  > ./bin/kvclient delete --key "mykey"
    {}
    
  > ./bin/kvclient get --key "mykey"
    kvclient: error: [404]: [404]: key 'mykey' not found"
```

## Deployment
- Run terraform to provision infrastructure/monitoring/alerts -- see runbook
- Run pipeline to deploy image -- see runbook

## Monitoring
- Dashboards links -- see runbook

## Endpoints
- GET `/v1/ping`
  Returns 200 OK

- GET `/v1/readiness`
  Returns 200 OK

- GET `/v1/keys/{key}`
    ```
  key can be any url encoded string upto 250 bytes
  Response is 200 with value of key OR
  Response is 404 if there is no such key
  key can be a path like: "a/b/c"
    ```
  
- DELETE `/v1/keys/{key}`
    ```
  key can be any url encoded string upto 250 bytes
  Response is 200 if the key no longer exists
    ``` 

- POST `/v1/keys/{key}`
    ```
  {
      "value": "data"
  }
  key can be any url encoded string upto 250 bytes
  value can be any data string upto 1MB
  Response is 200 if the key is written
    ``` 

## Client
```shell
 >  make build_client
 > ./bin/kvclient --help
```

## Server Source Code
- `internal/kv/` - source code for the key value pair; a simple map with a mutex
- `internal/httpserver/` - http endpoints that translate calls to the kv store.
  Every endpoint has a "decoder", "encoder", "endpoint" function. The "decoder" is for decoding  
  the JSON off the wire into a go request structure. The "encoding" is for encoding the go
  structure to JSON to deliver on the wire. The "endpoint" accepts a request struct and returns
  a response struct according to its business logic [in this case a straight up call to kv store].
- `cmd/kvservice/` - location of main.go for server

## Client Source Code
- `internal/clientcommands/` - source code for a command line tool that invokes a http client to call the kvservice.
  Commands are: get, write, delete
- `cmd/kvclient/` - location of main.go for client

## Common Source Code
- `apiclient/` - hand coded http calls to kvservice. The request/response structures are
  used by the server to encode/decode on the wire. They are also used by the command line tool
  to marshal command line parameters into a request structure; and output responses.
