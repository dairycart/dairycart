FROM golang:alpine
WORKDIR /go/src/github.com/dairycart/dairycart

ADD . .

ENTRYPOINT ["go", "test", "-v"]
