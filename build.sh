#!/usr/bin/env bash

set -ex

if [[ $OS_NAME == linux ]]; then
  apt update
  apt install -y binutils gcc clang ninja-build cmake build-essential file
fi

if [[ $OS_NAME == win32 ]]; then
  apt update
  apt install -y binutils-mingw-w64-i686 binutils-mingw-w64-x86-64 g++-mingw-w64-i686 gcc-mingw-w64 gcc-mingw-w64-base gcc-mingw-w64-x86-64 gcc-multilib ninja-build cmake build-essential file
fi

if [[ $OS_NAME == win32 ]]; then
  apt update
  apt install -y binutils-mingw-w64-i686 binutils-mingw-w64-x86-64 g++-mingw-w64-i686 gcc-mingw-w64 gcc-mingw-w64-base gcc-mingw-w64-x86-64 gcc-multilib ninja-build cmake build-essential file
fi

if [[ $OS_NAME == osxcross ]]; then
  apt update
  apt install -y ninja-build cmake build-essential file pkg-config vim
  wget -qO- https://dl.google.com/go/go1.13beta1.linux-amd64.tar.gz > go.tgz
  tar xf go.tgz -C /usr/lib
  export PATH=/usr/lib/go/bin/:$PATH
fi

LIBGIT2_URL=https://github.com/libgit2/libgit2.git
LIBGIT2_VER=v0.28.2

NPROC=$(nproc 2>/dev/null || sysctl -n hw.physicalcpu 2>/dev/null || printenv NUMBER_OF_PROCESSORS)
GIT2GO_PATH=$PWD/vendor/github.com/libgit2/git2go
LIBGIT2_PATH=$GIT2GO_PATH/vendor/libgit2
LIBGIT2_BUILD=$LIBGIT2_PATH/static-build
INSTALL=$GIT2GO_PATH/static-build/install

export GO111MODULE="on" GOFLAGS="-mod=vendor" CGO_ENABLED=1

go mod download
if ! [[ -d vendor ]]; then
  go mod vendor
fi

git clone --depth 1 -b $LIBGIT2_VER $LIBGIT2_URL $LIBGIT2_PATH || :

build_libgit2(){
mkdir -vp $LIBGIT2_BUILD/{,bld,install}
mkdir -p $LIBGIT2_BUILD/bld/libgit2-static && pushd $LIBGIT2_BUILD/bld/libgit2-static
cmake \
        -G Ninja \
        -DCMAKE_C_FLAGS=-fPIC \
        -DUSE_EXT_HTTP_PARSER=OFF \
        -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
        -DCMAKE_C_COMPILER=clang \
        -DUSE_BUNDLED_ZLIB=ON \
        -DUSE_HTTPS=OFF \
        -DUSE_SSH=OFF \
        -DBUILD_SHARED_LIBS=OFF \
        -DCMAKE_INSTALL_PREFIX="${INSTALL}" \
        -DBUILD_CLAR=OFF \
        -DTHREADSAFE=ON \
        ${LIBGIT2_PATH}
cmake --build . -- -j8
cmake --build . --target install
popd
}

build_libgit2_mingw(){
mkdir -vp $LIBGIT2_BUILD/{,bld,install}
mkdir -p $LIBGIT2_BUILD/bld/libgit2-static && pushd $LIBGIT2_BUILD/bld/libgit2-static
cmake \
        -G Ninja \
        -DCMAKE_C_FLAGS=-fPIC \
        -DUSE_EXT_HTTP_PARSER=OFF \
        -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
        -DWINHTTP=OFF \
        -DUSE_BUNDLED_ZLIB=ON \
        -DUSE_HTTPS=OFF \
        -DUSE_SSH=OFF \
        -DCMAKE_C_COMPILER=$CC \
        -DBUILD_SHARED_LIBS=OFF \
        -DCMAKE_INSTALL_PREFIX="${INSTALL}" \
        -DBUILD_CLAR=OFF \
        -DTHREADSAFE=ON \
        -DCMAKE_SYSTEM_NAME=Windows \
        -DWIN32=ON \
        -DMINGW=ON \
        -DCMAKE_SIZEOF_VOID_P=8 \
        ${LIBGIT2_PATH}
cmake --build . -- -j8
cmake --build . --target install
popd
}

build_osxcross(){
mkdir -vp $LIBGIT2_BUILD/{,bld,install}
mkdir -p $LIBGIT2_BUILD/bld/libgit2-static && pushd $LIBGIT2_BUILD/bld/libgit2-static
cmake -DTHREADSAFE=ON \
      -DBUILD_CLAR=OFF \
      -DBUILD_SHARED_LIBS=OFF \
      -DCMAKE_C_FLAGS="-fPIC -mmacosx-version-min=10.14" \
      -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
      -DCMAKE_INSTALL_PREFIX="${INSTALL}" \
      -DWINHTTP=OFF \
      -DUSE_BUNDLED_ZLIB=ON \
      -DUSE_HTTPS=OFF \
      -DUSE_SSH=OFF \
      -DCURL=OFF \
      -G "Ninja" \
      -DCMAKE_SIZEOF_VOID_P=8 \
      -DCMAKE_C_COMPILER=$CC \
      -DCMAKE_OSX_SYSROOT=/opt/osxcross/SDK/MacOSX10.14.sdk/ \
      -DCMAKE_SYSTEM_NAME=Darwin \
      ${LIBGIT2_PATH}
sed -i -e 's,-I/usr/include,,g' build.ninja
sed -i -e 's,-isystem /usr/include,,g' build.ninja
cmake --build . -- -j8
cmake --build . --target install
popd
}

if [[ $OS_NAME == linux ]]; then
  export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
  build_libgit2
fi

if [[ $OS_NAME == win32 ]]; then
  export GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc
  FLAGS="-lws2_32"
  export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
  build_libgit2_mingw
fi

if [[ $OS_NAME == win64 ]]; then
  export GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc
  FLAGS="-lws2_32"
  export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
  build_libgit2_mingw
fi

if [[ $OS_NAME == osxcross ]]; then
  export GOOS=darwin GOARCH=amd64 CC=x86_64-apple-darwin18-clang
  FLAGS=""
  export CFLAGS="-mmacosx-version-min=10.14"
  export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
  build_libgit2_osxcross
fi


pwd
go build -v --tags static .
echo $?
exit

case "$OS_NAME" in
    linux*)
go install -tags "static" -ldflags "-extldflags -static" github.com/libgit2/git2go
GO111MODULE=off go install -tags "static" -ldflags "-extldflags -static" github.com/navigaid/gitql
;;
    osxcross*)
go install -tags "static" github.com/libgit2/git2go
GO111MODULE=off go install -tags "static" github.com/navigaid/gitql
;;
    *)
echo UNKNOWN_PLATFORM
;;
esac

gitql="$(GO111MODULE=off go list -f '{{ .Target }}' github.com/navigaid/gitql)"
ldd $gitql || otool -L $gitql || true
