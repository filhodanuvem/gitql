export LD_LIBRARY_PATH=$(PWD)/libgit2/install/lib
export DYLD_LIBRARY_PATH=$(PWD)/libgit2/install/lib
export PKG_CONFIG_PATH=${PWD}/libgit2/install/lib/pkgconfig:/usr/local/opt/openssl/lib/pkgconfig
export C_INCLUDE_PATH=$(PWD)/libgit2/install/include
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

build-dynamic: prepare-dynamic
	@go get -v -d . 
	@echo ${LD_LIBRARY_PATH}
	@echo ${DYLD_LIBRARY_PATH}
	@ls ${DYLD_LIBRARY_PATH}
	@go build