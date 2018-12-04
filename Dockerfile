FROM golang:1.11-alpine

WORKDIR /go/src/github.com/qlik-oss/core-authorization/
COPY . /go/src/github.com/qlik-oss/core-authorization/
RUN apk add --no-cache curl git && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh  && \
    dep ensure
CMD go test -v ./access
