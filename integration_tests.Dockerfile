FROM golang:alpine
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD tests .
COPY vendor vendor
ENTRYPOINT ["go", "test", "-v"]