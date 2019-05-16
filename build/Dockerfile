FROM golang:latest

MAINTAINER xcat.org


VOLUME ["/src","/build"]
RUN mkdir  -p /$GOROOT/go/src/github.ibm.com/yangsbj/
RUN ln -s /$GOROOT/go/src/github.ibm.com/yangsbj/ /src
RUN mkdir -p /build
RUN cd  /$GOROOT/go/src/github.ibm.com/yangsbj/
RUN go get github.com/hashicorp/terraform
RUN go get github.com/thedevsaddam/gojsonq
#ENTRYPOINT go build -o /build/terraform-provider-xcat

