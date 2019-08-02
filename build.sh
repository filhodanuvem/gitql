#!/usr/bin/env bash

set -ex

setup_vendor(){
  go mod download
  if ! [[ -d vendor ]]; then
    go mod vendor
  fi

  git clone --depth 1 -b $LIBGIT2_VER $LIBGIT2_URL $LIBGIT2_PATH || :
}

build_libgit2_linux(){
  cmake \
  -G Ninja \
  -DCMAKE_C_FLAGS=-fPIE \
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

build_libgit2_windows(){
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

build_libgit2_darwin(){
  cmake -DTHREADSAFE=ON \
  -DBUILD_CLAR=OFF \
  -DBUILD_SHARED_LIBS=OFF \
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

build_libgit2(){
  mkdir -vp $LIBGIT2_BUILD $INSTALL
  pushd $LIBGIT2_BUILD

  case "$TARGET_OS_ARCH" in
  linux-amd64*)
    export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
    build_libgit2_linux
  ;;
  darwin-amd64*)
    export GOOS=darwin GOARCH=amd64 CC=x86_64-apple-darwin18-clang
    FLAGS=""
    export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
    build_libgit2_darwin
  ;;
  windows-amd64*)
    export GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-clang
    FLAGS="-lws2_32"
    export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
    build_libgit2_windows
  ;;
  windows-386*)
    export GOOS=windows GOARCH=386 CC=i686-w64-mingw32-clang
    FLAGS="-lws2_32"
    export CGO_LDFLAGS="${INSTALL}/lib/libgit2.a -L${INSTALL}/include ${FLAGS}"
    build_libgit2_windows
  ;;
  esac

  cmake --build . -- -j8
  cmake --build . --target install

  popd
}

build_gitql(){
  case "$TARGET_OS_ARCH" in
  linux-amd64*)
    go build -v -tags static -ldflags "-extldflags '-static'" .
  ;;
  darwin-amd64*)
    # MacOS doesn’t support fully static binaries, see 
    # https://stackoverflow.com/questions/3801011/ld-library-not-found-for-lcrt0-o-on-osx-10-6-with-gcc-clang-static-flag
    # this is the best we could possibly do
    go build -v -tags static .
  ;;
  windows-amd64*)
    go build -v -tags static -ldflags "-extldflags '-static'" .
  ;;
  windows-386*)
    go build -v -tags static -ldflags "-extldflags '-static'" .
  ;;
  esac
}

main(){
  export TARGET_OS_ARCH="${1:-linux-amd64}"
  export GO111MODULE="on" GOFLAGS="-mod=vendor" CGO_ENABLED=1
  export LIBGIT2_URL=https://github.com/libgit2/libgit2.git
  export LIBGIT2_VER=v0.28.2

  export NPROC=$(nproc 2>/dev/null || echo 4)
  export GIT2GO_PATH=$PWD/vendor/github.com/libgit2/git2go
  export LIBGIT2_PATH=$GIT2GO_PATH/vendor/libgit2
  export LIBGIT2_BUILD=$LIBGIT2_PATH/static-build
  export INSTALL=$GIT2GO_PATH/static-build/install

  setup_vendor
  build_libgit2
  build_gitql
}

main "$@"
