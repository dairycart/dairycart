# key stage
FROM alpine:latest AS key-stage
RUN apk --no-cache add --update openssh openssl

RUN ssh-keygen -t rsa -b 2048 -f app.rsa
RUN openssl rsa -in app.rsa -pubout -outform PEM -out app.rsa.pub

# final stage
FROM golang:alpine
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD api .

COPY --from=key-stage app.rsa app.rsa
COPY --from=key-stage app.rsa.pub app.rsa.pub

ENTRYPOINT ["go", "test", "-cover", "-tags", "test"]
