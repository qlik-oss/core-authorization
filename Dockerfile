FROM golang:1.10-alpine
WORKDIR /go/src/github.com/qlik-oss/core-athorization/
COPY . /go/src/github.com/qlik-oss/core-athorization/
CMD go test -v ./access
