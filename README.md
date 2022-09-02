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
- go-swagger [only to build api client from yaml]

## go-swagger
brew tap go-swagger/go-swagger
brew install go-swagger

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
`make docker`

Builds docker image for kvservice. 
Run docker: `docker run -p 8282:8282 kvservice:latest`

## Examples
```shell
  > make run   # runs kvservice server at port 8282 in foreground
  
  > ./bin/kvclient store --key "mykey" --value "myvalue"
  
  > ./bin/kvclient read --key "mykey"
    myvalue
    
  > ./bin/kvclient store --key "mykey" --value "myvalue2"
  
  > ./bin/kvclient read --key "mykey"
    myvalue2  
    
  > ./bin/kvclient delete --key "mykey"
 
  > ./bin/kvclient read --key "mykey"
    ERROR: no key `mykey` found
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
  key can be any url encoded string upto 250 characters
  Response is 200 with value of key OR
  Response is 404 if there is no such key
  TODO: make sure key can be like: "a/b/c"
    ```
  
- DELETE `/v1/keys/{key}`
    ```
  key can be any url encoded string upto 250 characters
  Response is 200 if the key no longer exists
    ``` 

- POST `/v1/keys/{key}?`
    ```
  {
      "value": "base64"
  }
  key can be any url encoded string upto 250 characters
  value can be any base64 encoded string upto 1MB
  Response is 200 if the key no longer exists
    ``` 