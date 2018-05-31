FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairycart

ADD . .

ENTRYPOINT ["go", "test", "-cover", "github.com/dairycart/dairycart/api/v1", "-parallel=1"]
