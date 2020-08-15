FROM golang:1.15-alpine

WORKDIR /go/src/github.com/qlik-oss/core-authorization/
COPY . /go/src/github.com/qlik-oss/core-authorization/
RUN apk add --no-cache curl git gcc musl-dev && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh  && \
    dep ensure
CMD go test -v ./access
