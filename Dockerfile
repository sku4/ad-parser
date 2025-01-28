FROM golang:1.23.5-alpine3.21 AS builder

RUN go version

COPY . /ad-parser/
WORKDIR /ad-parser/

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go mod download
RUN go build -o ./.bin/ad-parser -tags=go_tarantool_ssl_disable ./cmd/ad/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /ad-parser/.bin/ad-parser .
COPY --from=builder /ad-parser/configs/config.yml configs/config.yml

CMD /app/ad-parser
