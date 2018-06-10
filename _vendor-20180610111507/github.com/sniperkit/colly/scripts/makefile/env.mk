ifeq ($(OS),)
  BUILD_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
  OS := $(BUILD_OS)
endif

ifeq ($(CPU),)
  CPU := $(shell uname -m | sed -e 's/i[345678]86/i386/')
endif

PLATFORM = $(CPU)-$(OS)

ifeq ($(OS), sunos)
  OS = solaris
endif
