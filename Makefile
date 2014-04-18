export LD_LIBRARY_PATH=$(PWD)/libgit2/install/lib
export DYLD_LIBRARY_PATH=$(PWD)/libgit2/install/lib
export PKG_CONFIG_PATH=$(PWD)/libgit2/install/lib/pkgconfig
URL_BASE_GIT2GO=https://github.com/libgit2/git2go/archive
GIT2GO_VERSION=master
all: prepare build

test: 
	go test ./lexical/ ./parser/ ./semantical ./runtime

clean:
	rm -rf ./libgit2

prepare: clean
	@echo "Preparing...\n"
	chmod +x $(GOPATH)/src/github.com/libgit2/git2go/script/build-libgit2.sh
	$(GOPATH)/src/github.com/libgit2/git2go/script/build-libgit2.sh

build: 
	go build
	@echo "Ready to go!"
	
install:
	cp ./gitql /usr/local/bin/gitql
	ln -s /usr/local/bin/gitql /usr/local/bin/git-ql
	@echo "Git is in /usr/local/bin/gitql"
	@echo "You can also use: git ql -q 'query here'"
