# Make a release with
# make -j4 release

VERSION=$(shell git describe)
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION}"

.PHONY: build
build:
	go build ${LDFLAGS} -o ./dist/pluck ./cmd/*.go

clean:
	rm -f ./dist/pluck*

install: clean build
	@cd ./cmd && go install ${LDFLAGS}
	@cp -f ./dist/pluck ${GOPATH}/bin/pluck
	@which pluck
	@pluck --version

.PHONY: linuxarm
linuxarm:
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o ./dist/pluck_linux_arm ./cmd/*.go
	# cd dist && upx --brute pluck_linux_arm

.PHONY: linux64
linux64:
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ./dist/pluck_linux_amd64 ./cmd/*.go
	cd dist && upx --brute pluck_linux_amd64

.PHONY: windows
windows:
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ./dist/pluck_windows_amd64.exe ./cmd/*.go
	# cd dist && upx --brute pluck_windows_amd64.exe

.PHONY: osx
osx:
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ./dist/pluck_osx_amd64 ./cmd/*.go
	# cd dist && upx --brute pluck_osx_amd64

.PHONY: release
release: osx windows linux64 linuxarm
