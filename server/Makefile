
all: deps generate build

generate:
	go generate ./generate/...

build: generate
	go build -o ./bin/url_shortner ./cmd/main.go

build-alpine: deps generate
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/url_shortner ./cmd/main.go

run: build
	./bin/url_shortner --config=./config.yaml

tests:
	go test ./...

fmt:
	go fmt ./...

deps:
	dep ensure

docker-build:
	docker build . -t url_shortner

docker-run: docker-build
	docker run -p 8080:8080 -v $(shell pwd)/data:/tmp/data url_shortner

.PHONY: generate
