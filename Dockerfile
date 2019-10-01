FROM golang:1.12-alpine

WORKDIR /etc
ADD . . 
RUN apk update && \ 
    apk add git make clang build-base pkgconfig openssh cmake libssh2 libssh2-dev openssl bash ninja
RUN go get -u -d -v .
RUN ./install.sh 
RUN ./gitql -v
RUN echo "INSTALLED " $TARGET_OS_ARCH
