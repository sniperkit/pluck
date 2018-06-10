INSTALL_TARGETS += alfred

alfred: brew_cask
ifeq (, $(shell brew cask list | grep alfred))
	$(call LOG_INFO, install $@)
	brew cask install $@
	$(call LOG_SUCCESS, $@ is OK.)
else
	$(call LOG_SUCCESS, $@ is OK.)
endif