build-all: build-bin build-plugin

build-bin:
	@go build -o ../../../bin/bitcq main.go

build-plugin:
	@go build -buildmode=plugin -o ../../../shared/libs/bitcq/bitcq.so main.go