#!/bin/sh
set -x
set -e

# Set temp environment vars

## local vars
export PARALLEL_JOBS=${PARALLEL_JOBS:-"4"}

## vcs vars
export TINI_VCS_REPO_URI="${TINI_VCS_REPO_URI:-"github.com/krallin/tini"}"
export TINI_VCS_REPO_URL="${TINI_VCS_REPO_URL:-"https://${TINI_VCS_REPO_URI}"}"
export TINI_VCS_REPO_BRANCH=${TINI_VCS_REPO_BRANCH:-"master"}
export TINI_VCS_CLONE_DEPTH=${TINI_VCS_CLONE_DEPTH:-"1"}
export TINI_VCS_LOCAL_DIR=${TINI_VCS_LOCAL_DIR:-"${GOPATH}/src/${TINI_VCS_REPO_URI}"}
export TINI_VCS_BUILD_DIR=${TINI_VCS_BUILD_DIR:-"${GOPATH}/src/${TINI_VCS_REPO_URI}/build"}

## cmake vars
export TINI_CMAKE_BUILD_TYPE=${TINI_CMAKE_BUILD_TYPE:-"Release"}

# Install build deps
apk --no-cache --no-progress add cmake gcc musl-dev libgit2-dev@testing

git clone --recursive --depth=${TINI_VCS_CLONE_DEPTH} --branch=${TINI_VCS_REPO_BRANCH} ${TINI_VCS_REPO_URL} ${TINI_VCS_LOCAL_DIR}

# go to tini's clone dir, build and install the binaries
mkdir -p ${TINI_VCS_BUILD_DIR}
cd ${TINI_VCS_BUILD_DIR}
cmake -DCMAKE_BUILD_TYPE=${TINI_CMAKE_BUILD_TYPE} ..
make -j${PARALLEL_JOBS}
make install 

# Cleanup
rm -r ${TINI_VCS_LOCAL_DIR}

# Remove build deps
apk --no-cache --no-progress del go gcc musl-dev libgit2-dev