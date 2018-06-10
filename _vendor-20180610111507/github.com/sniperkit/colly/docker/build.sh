#!/bin/sh
set -x
set -e

# Set temp environment vars
export GOPATH=/tmp/go
export PATH=${PATH}:${GOPATH}/bin
# export BUILDPATH=${GOPATH}/src/github.com/sniperkit/colly
export BUILD_PATH=${APP_SRC:"${GOPATH}/src/github.com/sniperkit/colly"}

# Install build deps
apk --no-cache --no-progress add go gcc musl-dev libgit2-dev@testing

# Init go environment to build git2etcd
mkdir -p $(dirname ${BUILDPATH})
ln -s /app ${BUILD_PATH}
cd ${BUILD_PATH}
go get -v
go build

# Cleanup GOPATH
rm -r ${GOPATH}

# Remove build deps
apk --no-cache --no-progress del go gcc musl-dev libgit2-dev