# run with:
# docker build -t gitql .
# docker run -it --entrypoint /bin/sh gitql:latest

FROM golang:1.15.2-alpine3.12 as builder

WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/gitql

FROM alpine:3.12
RUN apk add -U git
COPY --from=builder /bin/gitql /bin/

ENTRYPOINT ["gitql"]