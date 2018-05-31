# build stage
FROM golang:latest AS build-stage
WORKDIR /go/src/github.com/dairycart/dairycart

ADD . .
RUN go build -o /dairycart github.com/dairycart/dairycart/cmd/server/v1

# final stage
FROM ubuntu:latest

ADD dairyconfigs/test_dairyconfig.toml dairyconfig.toml
COPY --from=build-stage /dairycart /dairycart

ENTRYPOINT ["/dairycart"]
