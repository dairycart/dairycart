FROM golang:alpine
WORKDIR /go/src/github.com/dairycart/dairycart

ADD api .

ENV DAIRYSECRET "do-not-use-secrets-like-this-plz"

ENTRYPOINT ["go", "test", "-cover", "-tags", "test"]
