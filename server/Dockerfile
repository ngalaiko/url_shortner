FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/ngalayko/url_shortner/server
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u golang.org/x/tools/cmd/goimports
RUN go get -u github.com/jteeuwen/go-bindata/...
ADD . ./
RUN make build-alpine

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/ngalayko/url_shortner/server/config.yaml .
COPY --from=builder /go/src/github.com/ngalayko/url_shortner/server/bin/url_shortner .
RUN date > /created_date

ENTRYPOINT ["./url_shortner", "--config=config.yaml"]
