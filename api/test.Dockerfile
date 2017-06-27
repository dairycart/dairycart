FROM golang:alpine
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD api .
ENTRYPOINT ["go", "test", "-cover", "-tags", "test"]
