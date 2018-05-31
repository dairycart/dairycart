FROM dairycart/dairydev:latest
WORKDIR /go/src/github.com/dairycart/dairycart

# this is meant to be run in the upper context
ADD . .
RUN go build -o /storage_gen github.com/dairycart/dairycart/cmd/gen

WORKDIR /go/src/github.com/dairycart/dairycart/storage/v1/database

ENTRYPOINT ["/storage_gen"]
