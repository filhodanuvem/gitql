export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$(PWD)/libgit2/install/lib
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$(PWD)/libgit2/install/lib
# export PKG_CONFIG_PATH=/usr/local/opt/openssl/lib/pkgconfig
export C_INCLUDE_PATH=$C_INCLUDE_PATH:$(PWD)/libgit2/install/include
URL_BASE_GIT2GO=https://github.com/libgit2/git2go/archive
GIT2GO_VERSION=master
GOPATH=$(shell go env GOPATH)

all: build

test: 
	@go test -count=1 -v ./parser/ ./lexical/ ./utilities/ ./semantical/ 

build: 
	@echo "Building..."
	@bash install.sh

clean:
	@rm -rf ./libgit2
	@rm -rf install-libgit2.sh

prepare-dynamic: clean
	@echo "Preparing...\n"
	@rm go.mod go.sum || echo 0
	@curl https://raw.githubusercontent.com/cloudson/git2go/original_libgit2/script/install-libgit2.sh >> install-libgit2.sh
	@chmod +x ./install-libgit2.sh
	@bash ./install-libgit2.sh
	@ls ${PWD}/libgit2/install/lib/pkgconfig
	@echo $(PKG_CONFIG_PATH)

build-dynamic: prepare-dynamic
	@pwd
	@ls ${PWD}/libgit2/install/lib/pkgconfig
	@echo $(PKG_CONFIG_PATH)
	@go get -v -d . 
	@go build