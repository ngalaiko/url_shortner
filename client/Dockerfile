FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/ngalayko/url_shortner/client
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/mailru/easyjson/...
RUN go get -u golang.org/x/tools/cmd/goimports
ADD . ./
RUN make build-alpine

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/ngalayko/url_shortner/client/config.yaml .
COPY --from=builder /go/src/github.com/ngalayko/url_shortner/client/bin/url_shortner .
RUN date > /created_date

ENTRYPOINT ["./url_shortner", "--config=config.yaml"]
