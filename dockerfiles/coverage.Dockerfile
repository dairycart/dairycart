FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairycart

# this Dockerfile should be executed with the upper directory's context
ADD . .
RUN mkdir /output

CMD ["go", "test", "-coverprofile=/output/coverage.out", "github.com/dairycart/dairycart/api/v1", "-parallel=1"]
