#!/usr/bin/env bash

set -ex

ZLIB_URL=https://github.com/madler/zlib.git
ZLIB_VER=v1.2.11
OPENSSL_URL=https://github.com/openssl/openssl.git
OPENSSL_VER=OpenSSL_1_1_0k
LIBSSH2_URL=https://github.com/libssh2/libssh2.git
LIBSSH2_VER=libssh2-1.8.2
CURL_URL=https://github.com/curl/curl.git
CURL_VER=curl-7_65_1
HTTPPARSER_URL=https://github.com/nodejs/http-parser.git
HTTPPARSER_VER=v2.9.2
LIBGIT2_URL=https://github.com/libgit2/libgit2.git
LIBGIT2_VER=v0.28.2

GIT2GO_URL=https://github.com/libgit2/git2go.git

# export CC="ccache gcc"

NPROC=$(nproc 2>/dev/null || sysctl -n hw.physicalcpu 2>/dev/null || printenv NUMBER_OF_PROCESSORS)
ROOT=$PWD
BASE=$ROOT/static-build

export CFLAGS="-I$ROOT/static-build/install/include"
export LDFLAGS="-L$ROOT/static-build/install/lib"
export PKG_CONFIG_PATH="$BASE/install/lib/pkgconfig"

sed -ie 's,giterr_,git_error_,g' `find . -name '*.go'`

mkdir -p $BASE $BASE/bld $BASE/install
clone(){
	git clone --depth 1 -b $ZLIB_VER $ZLIB_URL $ROOT/vendor/zlib || :
	git clone --depth 1 -b $OPENSSL_VER $OPENSSL_URL $ROOT/vendor/openssl || :
	git clone --depth 1 -b $LIBSSH2_VER $LIBSSH2_URL $ROOT/vendor/libssh2 || :
	git clone --depth 1 -b $CURL_VER $CURL_URL $ROOT/vendor/curl || :
	git clone --depth 1 -b $LIBGIT2_VER $LIBGIT2_URL $ROOT/vendor/libgit2 || :
	git clone --depth 1 -b $HTTPPARSER_VER $HTTPPARSER_URL $ROOT/vendor/http-parser || :
}
build_musl(){
	rm -rf $ROOT/vendor/musl
	mkdir -p $ROOT/vendor/musl && curl -sL https://www.musl-libc.org/releases/musl-1.1.22.tar.gz | tar xvzf - --strip=1 -C $ROOT/vendor/musl
	ln -sf $ROOT/vendor/musl $BASE/bld/musl && cd $ROOT/vendor/musl
	./configure --prefix=$BASE/install # --disable-shared
	make -j $NPROC
	make -j $NPROC install
}

# export CC="ccache $BASE/install/bin/musl-gcc"
# export LDFLAGS="-L"

build_zlib(){
	ln -sf $ROOT/vendor/zlib $BASE/bld/zlib && pushd $BASE/bld/zlib 
	./configure --prefix=$BASE/install # --static
	make -j $NPROC
	make -j $NPROC install
	popd
}
build_openssl(){
	mkdir -p $BASE/bld/openssl && pushd $BASE/bld/openssl
	# CC="ccache $BASE/install/bin/musl-gcc" 
	$ROOT/vendor/openssl/config --prefix=$BASE/install --with-zlib-lib=$BASE/install/lib/ --with-zlib-include=$BASE/install/lib/include/ zlib  # no-tests  no-shared 
	make -j $NPROC
	make -j $NPROC install_sw
	popd
}
build_libssh2(){
	ln -sf $ROOT/vendor/libssh2 $BASE/bld/libssh2 && pushd $BASE/bld/libssh2
	./buildconf
	./configure --prefix=$BASE/install # --disable-shared
	make -j $NPROC
	make -j $NPROC install
	popd
}
build_curl(){
	ln -sf $ROOT/vendor/curl $BASE/bld/curl && pushd $BASE/bld/curl
	./buildconf
	./configure --with-ssl=$BASE/install --with-libssh2=$BASE/install --prefix=$BASE/install # --disable-shared
	make -j $NPROC
	make -j $NPROC install
	popd
}
build_http_parser(){
	ln -sf $ROOT/vendor/http-parser $BASE/bld/http-parser && pushd $BASE/bld/http-parser
	# skip installing shared libs
	# make -j $NPROC
	# make -j $NPROC PREFIX=$BASE/install install
	make -j $NPROC package
	cp -v libhttp_parser.a $BASE/install/lib/libhttp_parser.a
	popd
}
build_libgit2_shared(){
	mkdir -p $BASE/bld/libgit2-shared && pushd $BASE/bld/libgit2-shared
	cmake -G Ninja -D CMAKE_C_FLAGS="-fPIC -Wno-stringop-truncation" -D CMAKE_BUILD_TYPE=RelWithDebInfo -D BUILD_SHARED_LIBS=ON -DCMAKE_INSTALL_PREFIX="${BASE}/install" $ROOT/vendor/libgit2 # -DCMAKE_CXX_COMPILER_LAUNCHER=ccache 
	cmake --build .
	cmake --build . --target install
	popd
}
build_libgit2_static(){
	mkdir -p $BASE/bld/libgit2-static && pushd $BASE/bld/libgit2-static
	cmake -G Ninja -D CMAKE_PREFIX_PATH="$BASE/install" -D CMAKE_C_FLAGS="-I$ROOT/static-build/install/include -Wno-stringop-truncation" -D CMAKE_BUILD_TYPE=RelWithDebInfo -D BUILD_SHARED_LIBS=OFF -DCMAKE_INSTALL_PREFIX="${BASE}/install" $ROOT/vendor/libgit2 # -DCMAKE_CXX_COMPILER_LAUNCHER=ccache 
	cmake --build .
	cmake --build . --target install
	popd
}

    clone
    build_zlib
    build_openssl
    build_libssh2
    build_curl
    build_http_parser
#   build_libgit2_shared
    build_libgit2_static


# https://github.com/cloudson/gitql/blob/master/install.sh
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux
CGO_ENABLED=1 go install -tags "static" -ldflags "-extldflags -static" .
CGO_ENABLED=1 go install -tags "static" -ldflags "-extldflags -static" github.com/navigaid/gitql
;;
    Darwin*)    machine=Mac
cat <<'EOF' >git_static.go
// +build static

package git

/*
#cgo windows CFLAGS: -I${SRCDIR}/static-build/install/include/
#cgo windows LDFLAGS: -L${SRCDIR}/static-build/install/lib/ -lgit2 -lwinhttp
#cgo !windows CFLAGS: -I${SRCDIR}/static-build/install/include/
#cgo !windows LDFLAGS: -L${SRCDIR}/static-build/install/lib/ ${SRCDIR}/static-build/install/lib/libcrypto.a ${SRCDIR}/static-build/install/lib/libz.a ${SRCDIR}/static-build/install/lib/libgit2.a ${SRCDIR}/static-build/install/lib/libhttp_parser.a ${SRCDIR}/static-build/install/lib/libssl.a ${SRCDIR}/static-build/install/lib/libssh2.a ${SRCDIR}/static-build/install/lib/libcurl.a -liconv -framework CoreFoundation -framework Security
// -l:libz.a -l:libgit2.a -l:ibhttp_parser.a -l:libssl.a -l:libssh2.a -l:libcurl.a -liconv -framework CoreFoundation -framework Security
//#cgo !windows LDFLAGS: -L${SRCDIR}/static-build/install/lib/ ${SRCDIR}/static-build/install/lib/libcrypto.a ${SRCDIR}/static-build/install/lib/libz.a ${SRCDIR}/static-build/install/lib/libgit2.a ${SRCDIR}/static-build/install/lib/libhttp_parser.a ${SRCDIR}/static-build/install/lib/libssl.a ${SRCDIR}/static-build/install/lib/libssh2.a ${SRCDIR}/static-build/install/lib/libcurl.a -liconv -framework CoreFoundation -framework Security
//-lgit2 -lhttp_parser -lz -lcurl -lssh2 -lcrypto -lssl
//#cgo !windows pkg-config: --static ${SRCDIR}/static-build/install/lib/pkgconfig/libgit2.pc
#include <git2.h>

#if LIBGIT2_VER_MAJOR != 0 || LIBGIT2_VER_MINOR != 28
# error "Invalid libgit2 version; this git2go supports libgit2 v0.28"
#endif

*/
import "C"
EOF
CGO_ENABLED=1 go install -tags "static" .
CGO_ENABLED=1 go install -tags "static" github.com/navigaid/gitql
;;
    CYGWIN*)    machine=Cygwin
 cp ./libgit2/install/lib/lib*  /usr/local/lib/
 ldconfig /usr/local/lib >/dev/null 2>&1 || echo "ldconfig not found">/dev/null
 cp ./gitql /usr/local/bin/gitql
 ln -s -f /usr/local/bin/gitql /usr/local/bin/git-ql
 echo "Git is in /usr/local/bin/gitql"
 echo "You can also use: git ql 'query here'"
;;
    MINGW*)     machine=MinGw
 # Everything here is ugly. But MinGW seem to have no ldconfig, 
 # no respect to /usr/local/bin, /usr/local/lib
 cp ./libgit2/install/bin/libgit2.dll  /usr/bin/
 cp ./gitql.exe /usr/bin/gitql.exe
 ln -s -f /usr/bin/gitql.exe /usr/bin/git-ql
 echo "Gitql.exe is in /usr/bin/gitql.exe"
 echo "You can also use: git ql 'query here'"
;;
    *)          machine="UNKNOWN:${unameOut}"
esac

echo $unameOut
gitql="$(go list -f '{{ .Target }}' github.com/navigaid/gitql)"
ldd $gitql || otool -L $gitql || true
