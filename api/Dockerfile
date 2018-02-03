# build stage
FROM golang:latest AS build-stage
WORKDIR /go/src/github.com/dairycart/dairycart

# this Dockerfile should be executed with the upper directory's context
ADD api api
ADD vendor vendor
ADD storage storage

RUN go build -o /dairycart github.com/dairycart/dairycart/api
COPY api/migrations/ /migrations

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=build-stage /dairycart /dairycart
COPY --from=build-stage /migrations /migrations

ENV DB_TO_USE "postgres"
ENV DAIRYSECRET "do-not-use-secrets-like-this-plz"
ENV MIGRATE_EXAMPLE_DATA "no thank you. :)"
ENV DAIRYCART_MIGRATIONS_DIR 'file://migrations'

ENTRYPOINT ["/dairycart"]
