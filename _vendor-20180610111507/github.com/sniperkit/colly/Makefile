APP_PREFIX := "colly"
APP_SUFFIX := ""
APP_NAME := "$(APP_PREFIX)$(SEPARATOR)$(APP_DIRNAME)$(APP_SUFFIX)"

PWD = $(shell pwd)
SEPARATOR := "-"
INSTALL_DIR := "/usr/local/bin/"
# DIST_DIRS := find * -type d -exec
DIST_DIR := "../../../dist"
BIN_DIR := "../../../bin"

# determine platform
ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  PLATFORM := OSX
else
  PLATFORM := Linux
endif

# dep 
# dep init

## glide
# yes no | glide create
# glide install --strip-vendor

PLATFORM_VERSION ?= $(shell uname -r)
PLATFORM_ARCH ?= $(shell uname -m)
PLATFORM_INFO ?= $(shell uname -a)
PLATFORM_OS ?= $(shell uname -s)

APP_ROOT := $(shell pwd)
APP_DIRNAME := $(shell basename `pwd`)
APP_PKG_URI ?= $(shell pwd | sed "s\#$(GOPATH)/src/\#\#g")
APP_PKG_URI_ARRAY ?= $(shell pwd | sed "s\#$(GOPATH)/src/\#\#g" | tr "/" "\n")

GO_EXECUTABLE ?= $(shell which go)
GO_VERSION ?= $(shell $(GO_EXECUTABLE) version)

GOX_EXECUTABLE ?= $(shell which gox)
GOX_VERSION ?= "master"

GLIDE_EXECUTABLE ?= $(shell which glide)
GLIDE_VERSION ?= $(shell $(GLIDE_EXECUTABLE) --version)

GODEPS_EXECUTABLE ?= $(shell which dep)
GODEPS_VERSION ?= $(shell $(GODEPS_EXECUTABLE) version | tr -s ' ')

GIT_EXECUTABLE ?= $(shell which git)
GIT_VERSION ?= $(shell $(GIT_EXECUTABLE) version)

SEPARATOR := "-"
INSTALL_DIR := "/usr/local/bin/"
# DIST_DIRS := find * -type d -exec
DIST_DIR := "../../../dist"
BIN_DIR := "../../../bin"


APP_ROOT := $(shell pwd)
APP_DIRNAME := $(shell basename `pwd`)
APP_PKG_URI ?= $(shell pwd | sed "s\#$(GOPATH)/src/\#\#g")
APP_PKG_URI_ARRAY ?= $(shell pwd | sed "s\#$(GOPATH)/src/\#\#g" | tr "/" "\n")
APP_PKG_DOMAIN ?= "$(word 1, $(APP_PKG_URI_ARRAY))"
APP_PKG_OWNER ?= "$(word 2, $(APP_PKG_URI_ARRAY))"
APP_PKG_NAME ?= "$(word 3, $(APP_PKG_URI_ARRAY))"
APP_PKG_URI_ROOT ?= "$(APP_PKG_DOMAIN)/$(APP_PKG_OWNER)/$(APP_PKG_NAME)"
APP_PKG_LOCAL_PATH ?= "$(GOPATH)/src/$(APP_PKG_URI_ROOT)"
APP_SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')
APP_PREFIX := $(shell basename $(APP_PKG_LOCAL_PATH))
APP_SUFFIX := ""
APP_NAME := "$(APP_PREFIX)$(SEPARATOR)$(APP_DIRNAME)$(APP_SUFFIX)"

VERSION ?= $(shell git describe --tags)
VERSION_INCODE = $(shell perl -ne '/^var version.*"([^"]+)".*$$/ && print "v$$1\n"' main.go)
VERSION_INCHANGELOG = $(shell perl -ne '/^\# Release (\d+(\.\d+)+) / && print "$$1\n"' CHANGELOG.md | head -n1)

VCS_GIT_REMOTE_URL = $(shell git config --get remote.origin.url)
VCS_GIT_VERSION ?= $(VERSION)

# NLIST = $(shell for x in {1..$(words $(APP_PKG_URI_ARRAY))}; do echo $$x; done)
# LIST = $(foreach x,$(NLIST), do_something_with_$(x)_and_$(word $(x),$(APP_PKG_URI_ARRAY)))

# @for key in "$${!APP_PKG_URI_ARRAY[@]:1:3}"; do echo "foo_$${key}_$${APP_PKG_URI_ARRAY[$${key}]}"; done
# @echo "$${!APP_PKG_URI_ARRAY[@]:0:1}"
# @echo "$${APP_PKG_URI_ARRAY[@]:0:1}"
# @IFS='/' && for VALUE in $${APP_PKG_URI_ARRAY[@]:0:1}; do echo "-VALUE=$${VALUE}"; $(call process_portal, $${VALUE}); done

print-%: ; @echo $*=$($*)

default: info

print-srcs: clear ## print the list of source files for this program
	@echo "APP_SRCS:"
	@echo "$(APP_SRCS)"

info: clear help ## build the linux binary
	@echo "\033[36mGLOBAL VARIABLES:\033[0m"
	@echo "$(SEPARATOR) APP_NAME: $(APP_NAME)"
	@echo "$(SEPARATOR) APP_ROOT: $(APP_ROOT)"
	@echo "$(SEPARATOR) APP_DIRNAME: $(APP_DIRNAME)"
	@echo "$(SEPARATOR) APP_PREFIX: $(APP_PREFIX)"
	@echo "$(SEPARATOR) APP_SUFFIX: $(APP_SUFFIX)"
	@echo "$(SEPARATOR) APP_PKG_URI: $(APP_PKG_URI)"
	@echo "$(SEPARATOR) APP_PKG_URI_ARRAY: $(APP_PKG_URI_ARRAY)"
	@echo "$(SEPARATOR) APP_PKG_URI_ROOT: $(APP_PKG_URI_ROOT)"
	@echo "$(SEPARATOR) APP_PKG_LOCAL_PATH: $(APP_PKG_LOCAL_PATH)"
	@echo "\033[36mVERSION VARIABLES:\033[0m"
	@echo "$(SEPARATOR) VERSION: $(VERSION)"
	@echo "$(SEPARATOR) VERSION_INCODE: $(VERSION_INCODE)"
	@echo "$(SEPARATOR) VERSION_INCHANGELOG: $(VERSION_INCHANGELOG)"
	@echo "\033[36mVCS VARIABLES:\033[0m"
	@echo "$(SEPARATOR) VCS_GIT_REMOTE_URL: $(VCS_GIT_REMOTE_URL)"
	@echo "$(SEPARATOR) VCS_GIT_VERSION: $(VCS_GIT_VERSION)"
	@echo "\033[36mDEPLOY VARIABLES:\033[0m"
	@echo "$(SEPARATOR) BUILD & INSTALL:"
	@echo "$(SEPARATOR) BIN_DIR: $(BIN_DIR)"
	@echo "$(SEPARATOR) DIST_DIR: $(DIST_DIR)"		
	@echo "\033[36mPLATFORM VARIABLES:\033[0m"
	@echo "$(SEPARATOR) PLATFORM: $(PLATFORM), PLATFORM_OS: $(PLATFORM_OS), PLATFORM_VERSION: $(PLATFORM_VERSION), PLATFORM_ARCH: $(PLATFORM_ARCH)"
	@echo "$(SEPARATOR) PLATFORM_INFO: $(PLATFORM_INFO)"
	@echo "\033[36mENV VARIABLES:\033[0m"
	@echo "$(SEPARATOR) GO_EXECUTABLE: $(GO_EXECUTABLE), GO_VERSION: $(GO_VERSION)"
	@echo "$(SEPARATOR) GLIDE_EXECUTABLE: $(GLIDE_EXECUTABLE), GLIDE_VERSION: $(GLIDE_VERSION)"
	@echo "$(SEPARATOR) GOX_EXECUTABLE: $(GOX_EXECUTABLE), GOX_VERSION: $(GOX_VERSION)"
	@echo "$(SEPARATOR) GODEPS_EXECUTABLE: $(GODEPS_EXECUTABLE), GODEPS_VERSION:  $(GODEPS_VERSION)"
	@echo "$(SEPARATOR) GIT_EXECUTABLE: $(GIT_EXECUTABLE), GIT_VERSION: $(GIT_VERSION)"
	@echo ""

clear: ## clear terminal screen
	@clear

build: ## build the executable for the current workstation
	${GO_EXECUTABLE} build -o bin/$(APP_DIRNAME) -ldflags "-X main.version=${VERSION}" main.go

install-local: ## local install the executable in /usr/local/bin
	${GO_EXECUTABLE} install -ldflags "-X main.version=${VERSION}" main.go

install: build ## install the executable in /usr/local/bin and ./bin directories
	install -d ${DESTDIR}$(INSTALL_DIR)
	install -m 755 ./bin/$(APP_DIRNAME) ${DESTDIR}$(INSTALL_DIR)$(APP_DIRNAME)

test: ## start tests available for this program
	${GO_EXECUTABLE} test . ./gb ./path ./action ./tree ./util ./godep ./godep/strip ./gpm ./cfg ./dependency ./importer ./msg ./repo ./mirrors

integration-test: ## start integration test
	${GO_EXECUTABLE} build

fmt: ## format code in source code files recursively.
	gofmt -s -l -w $(APP_SRCS)

dep-ensure: ## ensure/install program dependencies
	dep ensure -v

clean: ## clean all build files
	rm -rf ./dist
	rm -rf ./bin

bootstrap-dist:
	${GO_EXECUTABLE} get -u github.com/Masterminds/gox

build-darwin: bootstrap-dist ## cross-compile the program for apple/darwin based operating systems.
	$(GOX_EXECUTABLE) -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="darwin" \
	-arch="amd64 386" \
	-output="$(DIST_DIR)/{{.OS}}-{{.Arch}}/{{.Dir}}" .

build-linux: bootstrap-dist ## cross-compile the program for linux based operating systems.
	$(GOX_EXECUTABLE) -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux" \
	-arch="amd64 386" \
	-output="$(DIST_DIR)/{{.OS}}-{{.Arch}}/{{.Dir}}" .

build-win: bootstrap-dist ## cross-compile the program for windows based operating systems.
	$(GOX_EXECUTABLE) -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="windows" \
	-arch="amd64 386" \
	-output="$(DIST_DIR)/{{.OS}}-{{.Arch}}/{{.Dir}}" .

build-all: bootstrap-dist ## cross-compile the program for linux, darwin, windows, freebsd, openbsd, netbsd operating systems.
	$(GOX_EXECUTABLE) -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux darwin windows freebsd openbsd netbsd" \
	-arch="amd64 386 armv5 armv6 armv7 arm64 s390x" \
	-output="$(DIST_DIR)/{{.OS}}-{{.Arch}}/{{.Dir}}" .

dist: build-all ## build all dist version of the program and archive it
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(APP_NAME)-${VERSION}-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(APP_NAME)-${VERSION}-{}.zip {} \; && \
	cd ..

verify-version: ## verify app version
	@if [ "$(VERSION_INCODE)" = "v$(VERSION_INCHANGELOG)" ]; then \
		echo "$(APP_NAME): $(VERSION_INCHANGELOG)"; \
	elif [ "$(VERSION_INCODE)" = "v$(VERSION_INCHANGELOG)-dev" ]; then \
		echo "$(APP_NAME) (development): $(VERSION_INCHANGELOG)"; \
	else \
		echo "Version number in main.go does not match CHANGELOG.md"; \
		echo "main.go: $(VERSION_INCODE)"; \
		echo "CHANGELOG : $(VERSION_INCHANGELOG)"; \
		exit 1; \
	fi

EXAMPLE_GOLANGLIBS_SRC ?= "./_examples/custom/golanglibs/sitemap"
EXAMPLE_GOLANGLIBS_BIN ?= "$(CURDIR)/bin/golanglibs"

examples-list: ## display the list of examples to build 
	@echo "not ready yet"

build-example-golanglib: ## build example golanglib crawler
	@rm -fR $(EXAMPLE_GOLANGLIBS_BIN)
	@${GO_EXECUTABLE} build -o $(EXAMPLE_GOLANGLIBS_BIN) $(EXAMPLE_GOLANGLIBS_SRC)
	@$(EXAMPLE_GOLANGLIBS_BIN) --help


build-plugins: ## build colly plugin example `bitcq`
	@${GO_EXECUTABLE} build -buildmode=plugin -o ./shared/libs/bitcq.so ./addons/plugins/bitcq/main.go

.PHONY: build test install clean bootstrap-dist build-all dist integration-test verify-version info array-test uri all help dep-ensure fmt clear build-plugins examples-list build-plugins

help: ## display available makefile targets for this project
	@echo "\033[36mMAKEFILE TARGETS:\033[0m"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "- \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)