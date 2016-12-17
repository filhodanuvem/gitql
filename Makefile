export LD_LIBRARY_PATH=$(PWD)/libgit2/install/lib
export DYLD_LIBRARY_PATH=$(PWD)/libgit2/install/lib
export PKG_CONFIG_PATH=$(PWD)/libgit2/install/lib/pkgconfig
URL_BASE_GIT2GO=https://github.com/libgit2/git2go/archive
GIT2GO_VERSION=master
all: prepare build

test: 
	@go test ./lexical/ ./parser/ ./semantical ./runtime

clean:
	@rm -rf ./libgit2

prepare: clean
	@echo "Preparing...\n"
	@chmod +x $(GOPATH)/src/github.com/cloudson/git2go/script/build-libgit2.sh
	@$(GOPATH)/src/github.com/cloudson/git2go/script/build-libgit2.sh

build: 
	@echo "Building..."
	@go build
	@echo "Ready to go!"

install:
	@cp ./libgit2/install/lib/lib*  /usr/local/lib/
	@ldconfig /usr/local/lib >/dev/null 2>&1 || echo "ldconfig not found">/dev/null
	@cp ./gitql /usr/local/bin/gitql
	@ln -s -f /usr/local/bin/gitql /usr/local/bin/git-ql
	@echo "Git is in /usr/local/bin/gitql"
	@echo "You can also use: git ql 'query here'"
