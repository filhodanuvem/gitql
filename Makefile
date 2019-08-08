all: prepare build

test: static-prepare
	@(sed -e "$$ d" static_build.sh; echo "go test -tags static -count=1 \
		. ./lexical/ ./parser/ ./semantical ./runtime") | bash

clean:
	@rm -rf ./vendor

prepare:
	@echo -n "Preparing... "
	@test -d ./vendor \
	&& echo "./vendor already exists, skip preparation." \
	|| (sed -e "$$ d" static_build.sh; echo "setup_vendor") | bash

static-prepare: prepare
	@(sed -e "$$ d" static_build.sh; echo "build_libgit2") | bash

build: static-prepare
	@echo "Building..."
	@(sed -e "$$ d" static_build.sh; echo "build_gitql") | bash

static-build: static-prepare
	@echo "Building static..."
	@env TARGET_OS_ARCH=$(TARGET_OS_ARCH) ./static_build.sh
	@echo "Ready to go!"

install:
	@bash install.sh
