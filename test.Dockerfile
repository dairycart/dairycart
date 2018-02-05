FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairycart

# this Dockerfile should be executed with the upper directory's context
ADD . .

ENTRYPOINT ["go", "test", "-cover", "github.com/dairycart/dairycart/api", "-parallel=1"]
