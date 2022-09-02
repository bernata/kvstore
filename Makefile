OUTPUT_DIR=./bin
SERVICE_SOURCE_CODE=./cmd/kvservice
CLIENT_SOURCE_CODE=./cmd/kvclient

.PHONY: all
all: build_server build_client
	echo "all"

.PHONY: build_server
build_server: lint
	go build --mod=vendor -o ${OUTPUT_DIR}/kvservice ${SERVICE_SOURCE_CODE}

.PHONY: build_client
build_client: lint
	go build --mod=vendor -o ${OUTPUT_DIR}/kvclient ${CLIENT_SOURCE_CODE}

test: build_server build_client
	go test --mod=vendor --race -v -coverprofile=coverage.out ./...

integration_test: build_server build_client
	./integration_test.sh

.PHONY: run
run: build_server
	./bin/kvservice -p 8282


.PHONY: lint
lint:
	golangci-lint run -E gosec,gocyclo,goerr113,goimports,nestif,nilerr,predeclared,revive,rowserrcheck,stylecheck,tparallel,unconvert,wastedassign ./...

