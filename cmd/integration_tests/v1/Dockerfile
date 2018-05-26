FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairycart

# this is meant to be run in the upper context
ADD cmd/integration_tests/v1 .
ADD vendor vendor

ENTRYPOINT ["go", "test", "-v"]
