# colly-plugin-example

## Purpose

Showcase of a web scraper used as a go plugin

## Build status [![Build Status](https://travis-ci.org/prusya/colly-plugin-example.svg?branch=master)](https://travis-ci.org/prusya/colly-plugin-example)

## Description

Use https://github.com/sniperkit/colly/pkg in a plugin

## Requirements

go 1.8+

## Installation

```bash
go get -t -v github.com/snperkit/colly/addons/plugins/...
```

## Compilation

Inside your GOPATH directory

```bash
cd $GOPATH/src/github.com/snperkit/colly/addons/plugins/bitcq
go build -buildmode=plugin -o ../../shared/libs/bitcq.so ./bitcq/main.go
go build .
```

## Testing

```bash
go test -v ./...
```