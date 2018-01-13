FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairycart

# this Dockerfile should be executed with the upper directory's context
ADD api .
ADD storage .

ENV DB_TO_USE "postgres"
ENV DAIRYSECRET "do-not-use-secrets-like-this-plz"

ENTRYPOINT ["go", "test", "-cover", "-tags", "test"]
