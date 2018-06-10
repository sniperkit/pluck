INSTALL_TARGETS += tmux

tmux: brew
ifeq (, $(shell which tmux))
	$(call LOG_INFO, install $@ ...)
	brew install tmux
	$(call LOG_SUCCESS, $@ is OK.)
else
	$(call LOG_SUCCESS, $@ is OK.)
endif