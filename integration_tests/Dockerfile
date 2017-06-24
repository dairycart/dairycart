FROM golang:alpine
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD . .
ENTRYPOINT ["go", "test", "-v", "-bench=.", "-benchmem"]
