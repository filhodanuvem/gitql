export LD_LIBRARY_PATH=$(PWD)/libgit2/install/lib
export PKG_CONFIG_PATH=$(PWD)/libgit2/install/lib/pkgconfig
URL_BASE_GIT2GO=https://github.com/libgit2/git2go/archive
GIT2GO_VERSION=master
all: prepare build

test: 
	go test ./lexical/ ./parser/ ./semantical ./runtime

clean:
	rm -rf ./git2go ./libgit2

prepare: clean
	@echo "Preparing...\n"
	wget $(URL_BASE_GIT2GO)/$(GIT2GO_VERSION).tar.gz
	tar -xvf "./$(GIT2GO_VERSION).tar.gz"
	mv ./git2go-$(GIT2GO_VERSION) ./git2go
	chmod +x ./git2go/script/build-libgit2.sh
	./git2go/script/build-libgit2.sh
	@echo "installed"

build: 
	go build
	@echo "Read to go!"
	
install:
	cp ./gitql /usr/bin/gitql
	@echo "Git is in /usr/bin/gitql"
