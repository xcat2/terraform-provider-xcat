FROM golang:latest

MAINTAINER xcat.org


#:VOLUME ["/build"]
RUN mkdir  -p /go/src/github.ibm.com/yangsbj/
RUN ln -s /go/src/github.ibm.com/yangsbj/ /src
RUN mkdir -p /build
RUN cd  /go/src/github.ibm.com/yangsbj/
RUN go get github.com/hashicorp/terraform
RUN go get github.com/tidwall/gjson
ENTRYPOINT cd /go/src/github.ibm.com/yangsbj/terraform-provider-xcat && go build -o /build/terraform-provider-xcat

