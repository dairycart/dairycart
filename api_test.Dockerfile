FROM golang:alpine
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD api .
COPY vendor vendor
ENTRYPOINT ["go", "test", "-v", "-cover"]