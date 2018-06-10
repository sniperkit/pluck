INSTALL_TARGETS += go

go: brew
ifeq (, $(shell which go))
	$(call LOG_INFO, install $@)
	@brew install go
else
	$(call LOG_SUCCESS, $@ is OK.)
endif