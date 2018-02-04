FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairycart

# this Dockerfile should be executed with the upper directory's context
ADD api api
ADD vendor vendor
ADD storage storage

ENTRYPOINT ["go", "test", "-cover", "-tags", "test", "github.com/dairycart/dairycart/api"]
