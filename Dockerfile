# run with:
# docker build -t gitql .
# docker run -it --entrypoint /bin/sh gitql:latest

FROM golang:1.18.7-alpine3.16 as builder

WORKDIR /src
COPY go.mod .
COPY go.sum .
COPY main.go autocomplete.go version.txt ./
COPY lexical lexical
COPY parser  parser
COPY runtime runtime
COPY semantical semantical
COPY utilities utilities
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/gitql

FROM alpine:3.16
COPY --from=builder /bin/gitql /bin/

ENTRYPOINT ["gitql"]