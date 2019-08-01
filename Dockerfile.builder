FROM liushuyu/osxcross

RUN wget -qO- https://dl.google.com/go/go1.13beta1.linux-amd64.tar.gz | tar xvzf - -C /usr/lib && ln -sf /usr/lib/go/bin/go /usr/bin/go && ln -sf /usr/lib/go/bin/gofmt /usr/bin/gofmt && apt update && apt install -y binutils-mingw-w64-i686 binutils-mingw-w64-x86-64 g++-mingw-w64-i686 gcc-mingw-w64 gcc-mingw-w64-base gcc-mingw-w64-x86-64 gcc-multilib ninja-build cmake build-essential file pkg-config
