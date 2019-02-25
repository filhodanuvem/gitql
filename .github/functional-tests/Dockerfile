FROM "golang:latest"

RUN apt-get update && \ 
    apt-get install -y cmake
RUN go get -u -d github.com/cloudson/gitql
RUN cd $GOPATH/src/github.com/cloudson/gitql && \
    make && \
    make install
ENTRYPOINT bash $GOPATH/src/github.com/cloudson/gitql/gitql