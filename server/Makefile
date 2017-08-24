
all: deps run

build:
	go build -o ./bin/url_shortner ./cmd/url_shortner.go

run: build
	./bin/url_shortner

deps:
	dep ensure
