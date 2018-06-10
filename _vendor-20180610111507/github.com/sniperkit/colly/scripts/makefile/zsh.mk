INSTALL_TARGETS += zsh

zsh: brew
ifeq (zsh, $(shell which zsh))
	$(call LOG_INFO, install $@ ...)
	brew install zsh
	$(call LOG_SUCCESS, $@ is OK.)
else
	$(call LOG_SUCCESS, $@ is OK.)
endif
