FROM golang:alpine
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD api .
COPY vendor vendor
COPY migrations /migrations
ENTRYPOINT ["go", "test", "-v"]