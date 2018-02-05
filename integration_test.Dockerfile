# build stage
FROM golang:latest AS build-stage
WORKDIR /go/src/github.com/dairycart/dairycart

# this Dockerfile should be executed with the upper directory's context
ADD . .
RUN go build -o /dairycart

# final stage
FROM ubuntu:latest

ADD test_dairyconfig.toml dairyconfig.toml
COPY --from=build-stage /dairycart /dairycart

ENTRYPOINT ["/dairycart"]
