#!/usr/bin/env bash

set -ex

export GO111MODULE="on" GOFLAGS="-mod=vendor" CGO_ENABLED=1
LIBGIT2_URL=https://github.com/libgit2/libgit2.git
LIBGIT2_VER=v0.28.2

NPROC=$(nproc 2>/dev/null || echo 4)
GIT2GO_PATH=$PWD/vendor/github.com/libgit2/git2go
LIBGIT2_PATH=$GIT2GO_PATH/vendor/libgit2
LIBGIT2_BUILD=$LIBGIT2_PATH/static-build
INSTALL=$GIT2GO_PATH/static-build/install

setup_vendor(){
  go mod download
  if ! [[ -d vendor ]]; then
    go mod vendor
  fi

  git clone --depth 1 -b $LIBGIT2_VER $LIBGIT2_URL $LIBGIT2_PATH || :
}

build_linux(){
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
}

build_mingw(){
  cmake \
  -G Ninja \
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
}

build_osxcross(){
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
}

build(){
  mkdir -vp $LIBGIT2_BUILD $INSTALL
  pushd $LIBGIT2_BUILD

  case "$OS_NAME" in
  linux*)
    export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
    build_linux
  ;;
  win32*)
    export GOOS=windows GOARCH=386 CC=i686-w64-mingw32-clang
    FLAGS="-lws2_32"
    export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
    build_mingw
  ;;
  win64*)
    export GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-clang
    FLAGS="-lws2_32"
    export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
    build_mingw
  ;;
  osxcross*)
    export GOOS=darwin GOARCH=amd64 CC=x86_64-apple-darwin18-clang
    FLAGS=""
    export CFLAGS="-mmacosx-version-min=10.14"
    export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
    build_osxcross
  ;;
  *)
    echo '[ERROR_UNKNOWN_PLATFORM] please set OS_NAME to one of linux|win32|win64|osxcross'
    export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
    build_linux
  ;;
  esac

  cmake --build . -- -j8
  cmake --build . --target install

  popd
}

main(){
  setup_vendor
  build
  go build -v --tags static .
}

main
