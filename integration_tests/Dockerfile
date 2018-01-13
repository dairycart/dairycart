FROM golang:alpine
WORKDIR /go/src/github.com/dairycart/dairycart

ADD . .
ADD ../storage .

ENTRYPOINT ["go", "test", "-v"]
