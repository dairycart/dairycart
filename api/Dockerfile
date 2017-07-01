# build stage
FROM golang:alpine AS build-stage
WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/dairycart

ADD . .
RUN go build -o /dairycart
COPY migrations/ /migrations

# key stage
FROM alpine:latest AS key-stage
RUN apk --no-cache add --update openssh openssl

RUN ssh-keygen -t rsa -b 4096 -f app.rsa
RUN openssl rsa -in app.rsa -pubout -outform PEM -out app.rsa.pub

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=build-stage /dairycart /dairycart
COPY --from=build-stage /migrations /migrations

COPY --from=key-stage app.rsa app.rsa
COPY --from=key-stage app.rsa.pub app.rsa.pub

ENV MIGRATE_EXAMPLE_DATA $MIGRATE_EXAMPLE_DATA
ENV DAIRYCART_MIGRATIONS_DIR 'file://migrations'

ENTRYPOINT ["/dairycart"]
