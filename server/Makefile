
all: deps tests generate build

.PHONY: generate
generate:
	go generate ./...
	go-bindata -o ./template/data.go -pkg template ./template/data/...

.PHONY: build
build:
	env GOGC=off go build -o ./bin/url_shortner ./cmd/main.go

.PHONY: build-alpine
build-alpine: deps generate
	env GOGC=off CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/url_shortner ./cmd/main.go

.PHONY: build
run: build
	env GOGC=off ./bin/url_shortner --config=./config.yaml

.PHONY: tests
tests:
	env GOGC=off go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: deps
deps:
	dep ensure

