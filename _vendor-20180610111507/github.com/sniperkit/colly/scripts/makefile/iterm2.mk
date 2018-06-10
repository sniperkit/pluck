INSTALL_TARGETS += iterm2

iterm2: brew_cask
ifeq (, $(shell brew cask list | grep iterm2))
	$(call LOG_INFO, install $@ ...)
	brew cask install $@
	$(call LOG_SUCCESS, $@ is OK.)
else
	$(call LOG_SUCCESS, $@ is OK.)
endif