FROM golang:alpine
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD integration_tests .
COPY vendor vendor
ENTRYPOINT ["go", "test", "-v", "-bench=.", "-benchmem"]
