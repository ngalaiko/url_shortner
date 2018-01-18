
all: deps generate build

generate:
	go generate ./...
	go-bindata -o ./template/data.go -pkg template ./template/data/...

build:
	env GOGC=off go build -o ./bin/url_shortner ./cmd/main.go

build-alpine: deps generate
	env GOGC=off CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/url_shortner ./cmd/main.go

run: build
	env GOGC=off ./bin/url_shortner --config=./config.yaml

tests:
	env GOGC=off go test ./...

fmt:
	go fmt ./...

deps:
	dep ensure

.PHONY: generate
