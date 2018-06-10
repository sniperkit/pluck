FROM golang:alpine

COPY . /go/src/github.com/sniperkit/colly/plugins/data/collection/stats

ENV CGO_ENABLED=0
RUN apk add --no-cache git && \
    cd /go/src/github.com/sniperkit/colly/plugins/data/collection/stats && \
    go build -v -o /dogstatsd ./cmd/dogstatsd && \
    apk del git && \
    rm -rf /go/*

ENTRYPOINT ["/dogstatsd"]
