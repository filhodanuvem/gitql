.PHONY: all clean prepare-vendor prepare-libgit2 test static-build build install

all: test build

clean:
	@rm -rf ./vendor

prepare-vendor:
	@echo "Preparing vendor..."
	@test -d ./vendor \
	&& echo "./vendor already exists, skip preparation." \
	|| (sed -e "$$ d" static_build.sh; echo "setup_vendor") | bash

prepare-libgit2: prepare-vendor
	@echo "Preparing libgit2..."
	@(sed -e "$$ d" static_build.sh; echo "build_libgit2") | bash

test: prepare-libgit2
	@echo "Testing..."
	@(sed -e "$$ d" static_build.sh; echo "go test -count=1 \
		. ./lexical/ ./parser/ ./semantical ./runtime") | bash

static-build:
	@echo "Building static binary..."
	@env TARGET_OS_ARCH=$(TARGET_OS_ARCH) ./static_build.sh
	@echo "Ready to go!"

build: static-build

install:
	@install -m 755 -v gitql /usr/local/bin/gitql
	@git config --global alias.ql '! gitql'
	@echo "You can now use: git ql 'query here'"
