#!/usr/bin/env bash

set -ex
export GO111MODULE="on"
export GOFLAGS="-mod=vendor -tags=static"
export CGO_ENABLED=1

export GIT2GO_PATH="${PWD}/vendor/github.com/libgit2/git2go"
export LIBGIT2_PATH="${GIT2GO_PATH}/vendor/libgit2"
export LIBGIT2_BUILD="${GIT2GO_PATH}/static-build/${TARGET_OS_ARCH}"
export LIBGIT2_STATIC_PREFIX="${GIT2GO_PATH}/static-build/install"

setup_vendor(){
  LIBGIT2_URL=https://github.com/libgit2/libgit2.git
  LIBGIT2_VER=v0.28.2

  go mod download

  if ! [[ -d vendor ]]; then
    go mod vendor
  fi

  git -c advice.detachedHead=false clone --quiet --depth 1 -b $LIBGIT2_VER $LIBGIT2_URL $LIBGIT2_PATH || :
}

build_libgit2_generic(){
  cmake \
  -G Ninja \
  -DTHREADSAFE=ON \
  -DBUILD_CLAR=OFF \
  -DBUILD_SHARED_LIBS=OFF \
  -DUSE_EXT_HTTP_PARSER=OFF \
  -DUSE_BUNDLED_ZLIB=ON \
  -DUSE_HTTPS=OFF \
  -DUSE_SSH=OFF \
  -DUSE_ICONV=OFF \
  -DCMAKE_C_FLAGS=-fPIE \
  -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
  -DCMAKE_INSTALL_PREFIX="${LIBGIT2_STATIC_PREFIX}" \
  ${LIBGIT2_PATH}
}

build_libgit2_linux(){
  cmake \
  -G Ninja \
  -DTHREADSAFE=ON \
  -DBUILD_CLAR=OFF \
  -DBUILD_SHARED_LIBS=OFF \
  -DUSE_EXT_HTTP_PARSER=OFF \
  -DUSE_BUNDLED_ZLIB=ON \
  -DUSE_HTTPS=OFF \
  -DUSE_SSH=OFF \
  -DUSE_ICONV=OFF \
  -DCMAKE_SYSTEM_NAME=Linux \
  -DCMAKE_C_COMPILER=${CC} \
  -DCMAKE_C_FLAGS=-fPIE \
  -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
  -DCMAKE_INSTALL_PREFIX="${LIBGIT2_STATIC_PREFIX}" \
  ${LIBGIT2_PATH}
}

build_libgit2_windows(){
  cmake \
  -G Ninja \
  -DTHREADSAFE=ON \
  -DBUILD_CLAR=OFF \
  -DBUILD_SHARED_LIBS=OFF \
  -DUSE_EXT_HTTP_PARSER=OFF \
  -DUSE_BUNDLED_ZLIB=ON \
  -DUSE_HTTPS=OFF \
  -DUSE_SSH=OFF \
  -DUSE_ICONV=OFF \
  -DWINHTTP=OFF \
  -DCMAKE_SYSTEM_NAME=Windows \
  -DCMAKE_C_COMPILER=${CC} \
  -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
  -DCMAKE_INSTALL_PREFIX="${LIBGIT2_STATIC_PREFIX}" \
  -DWIN32=ON \
  -DMINGW=ON \
  -DCMAKE_SIZEOF_VOID_P=8 \
  ${LIBGIT2_PATH}
}

build_libgit2_darwin(){
  cmake \
  -G "Ninja" \
  -DTHREADSAFE=ON \
  -DBUILD_CLAR=OFF \
  -DBUILD_SHARED_LIBS=OFF \
  -DUSE_EXT_HTTP_PARSER=OFF \
  -DUSE_BUNDLED_ZLIB=ON \
  -DUSE_HTTPS=OFF \
  -DUSE_SSH=OFF \
  -DUSE_ICONV=OFF \
  -DCMAKE_SYSTEM_NAME=Darwin \
  -DCMAKE_OSX_DEPLOYMENT_TARGET="10.14" \
  -DCMAKE_C_COMPILER=${CC} \
  -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
  -DCMAKE_INSTALL_PREFIX="${LIBGIT2_STATIC_PREFIX}" \
  -DCMAKE_OSX_SYSROOT="${OSX_SDK}" \
  ${LIBGIT2_PATH}
}

build_libgit2(){
  mkdir -vp $LIBGIT2_BUILD
  pushd $LIBGIT2_BUILD

  case "$TARGET_OS_ARCH" in
  "")
    FLAGS=""
    export CGO_LDFLAGS="${LIBGIT2_STATIC_PREFIX}/lib/libgit2.a -L${LIBGIT2_STATIC_PREFIX}/include ${FLAGS}"
    build_libgit2_generic
  ;;
  linux/amd64*)
    export GOOS=linux GOARCH=amd64 CC=clang
    FLAGS=""
    export CGO_LDFLAGS="${LIBGIT2_STATIC_PREFIX}/lib/libgit2.a -L${LIBGIT2_STATIC_PREFIX}/include ${FLAGS}"
    build_libgit2_linux
  ;;
  darwin/amd64*)
    export GOOS=darwin GOARCH=amd64 CC=x86_64-apple-darwin18-clang OSX_SDK=/opt/osxcross/SDK/MacOSX10.14.sdk/
    FLAGS=""
    export CGO_LDFLAGS="${LIBGIT2_STATIC_PREFIX}/lib/libgit2.a -L${LIBGIT2_STATIC_PREFIX}/include ${FLAGS}"
    build_libgit2_darwin
  ;;
  windows/amd64*)
    export GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-clang
    FLAGS="-lws2_32"
    export CGO_LDFLAGS="${LIBGIT2_STATIC_PREFIX}/lib/libgit2.a -L${LIBGIT2_STATIC_PREFIX}/include ${FLAGS}"
    build_libgit2_windows
  ;;
  windows/386*)
    export GOOS=windows GOARCH=386 CC=i686-w64-mingw32-clang
    FLAGS="-lws2_32"
    export CGO_LDFLAGS="${LIBGIT2_STATIC_PREFIX}/lib/libgit2.a -L${LIBGIT2_STATIC_PREFIX}/include ${FLAGS}"
    build_libgit2_windows
  ;;
  esac

  cmake --build . -- -j$(nproc 2>/dev/null || sysctl -n hw.ncpu) && cmake --build . --target install &>/dev/null

  popd
}

build_gitql(){
  case "$TARGET_OS_ARCH" in
  windows/amd64*)
  ;&
  windows/386*)
  ;&
  linux/amd64*)
    go build -v -ldflags '-extldflags -static' .
  ;;
  darwin/amd64*)
    # MacOS doesn't support fully static binaries, see
    # https://stackoverflow.com/questions/3801011/ld-library-not-found-for-lcrt0-o-on-osx-10-6-with-gcc-clang-static-flag
    # this is the best we could possibly do
  ;&
  "")
  ;&
  *)
    go build -v .
  ;;
  esac
}

main(){
  setup_vendor
  build_libgit2
  build_gitql
}

main "$@"
