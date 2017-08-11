# build stage
FROM golang:alpine AS build-stage
WORKDIR /go/src/github.com/dairycart/dairycart/admin

ADD server .
RUN go build -o /admin-server
COPY templates/ /templates
COPY dist/ /dist

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=build-stage /admin-server /admin-server
COPY --from=build-stage /templates /templates
COPY --from=build-stage /dist /dist

ENTRYPOINT ["/admin-server"]
