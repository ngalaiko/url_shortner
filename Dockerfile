FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/ngalayko/url_shortner/server
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u golang.org/x/tools/cmd/goimports
RUN go get -u github.com/jteeuwen/go-bindata/...
RUN go get -u gopkg.in/reform.v1/...
ADD ./server ./
RUN make build-alpine

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/ngalayko/url_shortner/server/bin/url_shortner .
COPY --from=builder /go/src/github.com/ngalayko/url_shortner/server/template/static static

VOLUME /data/shortner

ADD docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
