FROM golang:1.10-alpine

WORKDIR /go/src/github.com/qlik-oss/core-athorization/
COPY . /go/src/github.com/qlik-oss/core-athorization/
RUN apk add --no-cache curl git && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh  && \
    dep ensure
CMD go test -v ./access
