INSTALL_TARGETS += vim dein.vim

lua: brew
ifeq (, $(shell which lua))
	brew install lua
endif

vim: brew lua
ifeq (, $(shell brew list macvim))
	$(call LOG_INFO, install $@ ...)
	brew install macvim --with-lua --with-override-system-vim --with-python3 > /dev/null
	$(call LOG_SUCCESS, $@ is OK.)
else
	$(call LOG_SUCCESS, $@ is OK.)
endif

dein.vim: git vim simple_deploy
ifeq (, $(shell find ~/.config/nvim -iname dein))
	$(call LOG_INFO, install $@ ...)
	@bash etc/init/common/dein.vim.sh $(HOME)/.config/nvim/dein
	@vim -c "call dein#install() | q"
	$(call LOG_SUCCESS, $@ is OK.)
else
	$(call LOG_SUCCESS, $@ is OK.)
endif