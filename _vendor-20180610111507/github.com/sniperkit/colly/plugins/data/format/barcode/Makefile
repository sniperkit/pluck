# Makefile 

# local - vars
CURRENT_DIR = $(shell pwd)



ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  PLATFORM := MacOSX
  PLATFORM_ALIAS := osx
  DEP_MANAGER := $(shell which brew)
else
  PLATFORM := Linux
  PLATFORM_ALIAS := nix
endif

# makefile - vars
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))

# makefile - targets
all: deps

deps: install-zbar

update-dep-osx:
	@brew update
	@brew reinstall zbar

install-zbar:
ifeq ($(PLATFORM_ALIAS), "osx")
	$(DEP_MANAGER) install hg; 
	$(DEP_MANAGER) install zbar;
else
	cd $(CURRENT_DIR);
	./scripts/build_zbar.sh
endif

coverage:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out


