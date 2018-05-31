FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairycart

# this is meant to be run in the upper context
ADD . .

ENTRYPOINT ["go", "test", "github.com/dairycart/dairycart/cmd/integration_tests/v1"]
