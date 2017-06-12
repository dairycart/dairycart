FROM golang:alpine
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD api .
COPY api/vendor vendor
ENTRYPOINT ["go", "test", "-cover"]
