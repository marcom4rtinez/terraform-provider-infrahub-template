default: fmt lint install generate

all: automatic_generator generate_sdk fmt lint install generate

generate_deploy: automatic_generator generate_sdk fmt lint release upload_registry clean celebrate

celebrate:
	clear;echo "\n\nPackaging finished:  📦📦 \nDeployment finished: 🎉🎉\n\n\nNow the real work starts! 📝\n\n\nProceed to use your Provider in your main.tf 💼"

clean:
	rm -Rf dist/; rm registry-manifest.json

release:
	goreleaser release --skip=publish --clean; goreleaser release

upload_registry:
	curl -X POST -L 'http://localhost:8080/v1/providers/marcom4rtinez/infrahub-main/upload' -H 'Content-Type: application/json' -d @registry-manifest.json

automatic_generator:
	cd generator; go run *.go

generate_sdk:
	cd sdk; curl -o schema.graphql http://localhost:8000/schema.graphql; go run github.com/Khan/genqlient

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
