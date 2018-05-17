FROM golang:latest
WORKDIR /go/src/github.com/dairycart/postgres

ADD . .

CMD go build -buildmode=plugin -o /output/result.so