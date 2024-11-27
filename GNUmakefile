default: fmt lint install generate

all: automatic_generator generate_sdk fmt lint install generate

automatic_generator:
	cd generator; go run *.go

generate_sdk:
	cd sdk; go run github.com/Khan/genqlient

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
