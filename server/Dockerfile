FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/ngalayko/url_shortner/server
RUN go get -u github.com/golang/dep/cmd/dep
ADD . ./
RUN make build-alpine

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/ngalayko/url_shortner/server/config.yaml .
COPY --from=builder /go/src/github.com/ngalayko/url_shortner/server/bin/url_shortner .
RUN date > /created_date
EXPOSE 80

CMD ["./url_shortner", "--config=config.yaml"]
