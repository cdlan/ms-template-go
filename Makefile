MODULE_NAME := ms-template-go
IMAGE_NAME := ms-template-go

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/server cmd/server/server.go

run:
	go run cmd/server/server.go

docker:
	docker build -t $(IMAGE_NAME):local -f build/Dockerfile .

.PHONY: docs
docs:
	docker build -t $(IMAGE_NAME)/docs:main -f build/docs.Dockerfile .
	docker run --rm --name server-docs -p 9090:80 $(IMAGE_NAME)/docs:main

grpc:
	scripts/gen_grpc_classes.sh

test:
	echo "### UNIT TEST ###"
	$(MAKE) unit-test

	echo "### INTEGRATION TEST ###"
	$(MAKE) integration-test

unit-test:
	go test -v -run TestUnitTestSuite ./...

integration-test:
	go test -v -run TestIntegrationTestSuite ./...