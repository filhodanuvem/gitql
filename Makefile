test: 
	go test ./lexical/ ./parser/ ./semantical ./runtime

install: 
	git clone https://github.com/libgit2/git2go.git
	ls ./git2go/script
	chmod +x ./git2go/script/build-libgit2.sh
	./git2go/script/build-libgit2.sh
	export LD_LIBRARY_PATH=./libgit2/install/lib
	export PKG_CONFIG_PATH=./libgit2/install/lib/pkgconfig
	go get 