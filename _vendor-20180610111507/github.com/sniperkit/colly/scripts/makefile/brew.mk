INSTALL_TARGETS += brew

brew:
ifeq (, $(shell which brew))
	$(call LOG_INFO, install $@ ...)
	@/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
	$(call LOG_SUCCESS, $@ is OK.)
else
	$(call LOG_SUCCESS, $@ is OK.)
endif