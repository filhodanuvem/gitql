export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$(PWD)/libgit2/install/lib
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$(PWD)/libgit2/install/lib
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$(PWD)/libgit2/install/lib/pkgconfig
export C_INCLUDE_PATH=$C_INCLUDE_PATH:$(PWD)/libgit2/install/include
URL_BASE_GIT2GO=https://github.com/libgit2/git2go/archive
GIT2GO_VERSION=master
GOPATH=$(shell go env GOPATH)

all: prepare build

test: 
	@go test -count=1 -v ./parser/ ./lexical/ ./utilities/ ./semantical/ 

clean:
	@rm -rf ./libgit2

prepare: clean
	@echo "Preparing...\n"
	@chmod +x $(GOPATH)/src/github.com/cloudson/git2go/script/build-libgit2.sh
	@$(GOPATH)/src/github.com/cloudson/git2go/script/build-libgit2.sh

build: 
	@echo "Building..."
	@bash install.sh