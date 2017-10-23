FROM golang:alpine
WORKDIR /go/src/github.com/dairycart/dairycart

ADD api .
ADD api/storage api/storage

ENV DB_TO_USE "postgres"
ENV DAIRYSECRET "do-not-use-secrets-like-this-plz"

ENTRYPOINT ["go", "test", "-cover", "-tags", "test"]
