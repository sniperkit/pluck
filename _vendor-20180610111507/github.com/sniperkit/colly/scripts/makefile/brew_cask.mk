INSTALL_TARGETS += brew_cask

brew_cask: brew
ifeq (, $(shell brew list | grep cask))
	$(call LOG_INFO, install $@ ...)
	brew tap caskroom/cask
	brew install cask
	$(call LOG_SUCCESS, $@ is OK.)
else
	$(call LOG_SUCCESS, $@ is OK.)
endif