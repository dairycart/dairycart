FROM golang:alpine
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD integration_tests .
COPY integration_tests/vendor vendor
ENTRYPOINT ["go", "test", "-v"]
