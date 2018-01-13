FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairycart

ADD integration_tests .
ADD vendor vendor

ENTRYPOINT ["go", "test", "-v"]
