#!/usr/bin/env bash

set -ex

ZLIB_URL=https://github.com/madler/zlib.git
ZLIB_VER=v1.2.11
OPENSSL_URL=https://github.com/openssl/openssl.git
OPENSSL_VER=OpenSSL_1_1_1c
LIBSSH2_URL=https://github.com/libssh2/libssh2.git
LIBSSH2_VER=libssh2-1.8.2
CURL_URL=https://github.com/curl/curl.git
CURL_VER=curl-7_65_1
HTTPPARSER_URL=https://github.com/nodejs/http-parser.git
HTTPPARSER_VER=v2.9.2
LIBGIT2_URL=https://github.com/libgit2/libgit2.git
LIBGIT2_VER=v0.28.2

GIT2GO_URL=https://github.com/libgit2/git2go.git

export CC="ccache gcc"

NPROC=$(nproc 2>/dev/null || echo 1)
ROOT=$PWD
BASE=$ROOT/static-build
mkdir -p $BASE $BASE/bld $BASE/install
{
	git clone --depth 1 -b $ZLIB_VER $ZLIB_URL $ROOT/vendor/zlib || :
	git clone --depth 1 -b $OPENSSL_VER $OPENSSL_URL $ROOT/vendor/openssl || :
	git clone --depth 1 -b $LIBSSH2_VER $LIBSSH2_URL $ROOT/vendor/libssh2 || :
	git clone --depth 1 -b $CURL_VER $CURL_URL $ROOT/vendor/curl || :
	git clone --depth 1 -b $LIBGIT2_VER $LIBGIT2_URL $ROOT/vendor/libgit2 || :
	git clone --depth 1 -b $HTTPPARSER_VER $HTTPPARSER_URL $ROOT/vendor/http-parser || :
}
musl(){
	rm -rf $ROOT/vendor/musl
	mkdir -p $ROOT/vendor/musl && curl -sL https://www.musl-libc.org/releases/musl-1.1.22.tar.gz | tar xvzf - --strip=1 -C $ROOT/vendor/musl
	ln -sf $ROOT/vendor/musl $BASE/bld/musl && cd $ROOT/vendor/musl
	./configure --prefix=$BASE/install # --disable-shared
	make -j $NPROC
	make -j $NPROC install
}

# export CC="ccache $BASE/install/bin/musl-gcc"
# export LDFLAGS="-L"

{
	ln -sf $ROOT/vendor/zlib $BASE/bld/zlib && cd $BASE/bld/zlib 
	./configure --prefix=$BASE/install # --static
	make -j $NPROC
	make -j $NPROC install
}
{
	mkdir -p $BASE/bld/openssl && cd $BASE/bld/openssl
	# CC="ccache $BASE/install/bin/musl-gcc" 
	$ROOT/vendor/openssl/config --prefix=$BASE/install --with-zlib-lib=$BASE/install/lib/ --with-zlib-include=$BASE/install/lib/include/ zlib no-tests # no-shared 
	make -j $NPROC
	make -j $NPROC install_sw
}
{
	ln -sf $ROOT/vendor/libssh2 $BASE/bld/libssh2 && cd $BASE/bld/libssh2
	./buildconf
	./configure --prefix=$BASE/install --disable-shared
	make -j $NPROC
	make -j $NPROC install
}
{
	ln -sf $ROOT/vendor/curl $BASE/bld/curl && cd $BASE/bld/curl
	./buildconf
	./configure --with-ssl=$BASE/install --with-libssh2=$BASE/install --prefix=$BASE/install --disable-shared
	make -j $NPROC
	make -j $NPROC install
}
{
	ln -sf $ROOT/vendor/http-parser $BASE/bld/http-parser && cd $BASE/bld/http-parser
	# skip installing shared libs
	# make -j $NPROC
	# make -j $NPROC PREFIX=$BASE/install install
	make -j $NPROC package
	cp -v libhttp_parser.a $BASE/install/lib/libhttp_parser.a
}
shared(){
	mkdir $BASE/bld/libgit2-shared && cd $BASE/bld/libgit2-shared
	cmake -G Ninja -D CMAKE_C_FLAGS="-fPIC -Wno-stringop-truncation" -D CMAKE_BUILD_TYPE=RelWithDebInfo -D BUILD_SHARED_LIBS=ON -DCMAKE_INSTALL_PREFIX="${BASE}/install" -DCMAKE_CXX_COMPILER_LAUNCHER=ccache $ROOT/vendor/libgit2
	cmake --build .
	cmake --build . --target install
}
{
	mkdir $BASE/bld/libgit2-static && cd $BASE/bld/libgit2-static
	cmake -G Ninja -D CMAKE_C_FLAGS="-fPIC -Wno-stringop-truncation" -D CMAKE_BUILD_TYPE=RelWithDebInfo -D BUILD_SHARED_LIBS=OFF -DCMAKE_INSTALL_PREFIX="${BASE}/install" -DCMAKE_C_COMPILER_LAUNCHER=ccache $ROOT/vendor/libgit2
	cmake --build .
	cmake --build . --target install
}

# export GOPATH="/go" 
# export PATH="$GOPATH/bin:/usr/local/go/bin:$PATH" 
# cd "$GOPATH/src/github.com/libgit2/git2go"
pwd
cd $ROOT
go install -tags "static" -ldflags "-extldflags '-static'" ./... || :
go run -tags "static" -ldflags "-extldflags '-static'" ./script/check-MakeGitError-thread-lock.go || :
