OUTPUT_DIR=./bin
LINUX_OUTPUT_DIR=./linux/bin
SERVICE_SOURCE_CODE=./cmd/kvservice
CLIENT_SOURCE_CODE=./cmd/kvclient

.PHONY: all
all: build_server build_client build_server_centos
	echo "all"

.PHONY: build_server
build_server: lint
	go build --mod=vendor -o ${OUTPUT_DIR}/kvservice ${SERVICE_SOURCE_CODE}

.PHONY: build_server_centos
build_server_centos: lint
	GOOS=linux GOARCH=amd64 go build --mod=vendor -o ${LINUX_OUTPUT_DIR}/kvservice ${SERVICE_SOURCE_CODE}

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

.PHONY: docker
docker: build_server_centos
	echo "Making Docker"
	docker build --no-cache -t kvservice -f Dockerfile ${LINUX_OUTPUT_DIR}

dockerrun:
	echo "Run Docker"
	docker run --rm -p 8282:8282 kvservice

.PHONY: lint
lint:
	golangci-lint run -E gosec,gocyclo,goerr113,goimports,nestif,nilerr,predeclared,revive,rowserrcheck,stylecheck,tparallel,unconvert,wastedassign ./...

.PHONY: clean
clean:
	rm -f bin/*
	rm -f linux/bin/*
