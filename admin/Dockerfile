# sass stage
FROM jbergknoff/sass AS sass-stage

ADD src .
RUN sass sass/*.sass app.css

# typescript stage
FROM sandrokeil/typescript:latest AS typescript-stage

ADD src .
RUN tsc typescript/*.ts --outFile /app.js

# build stage
FROM golang:alpine AS build-stage
WORKDIR /go/src/github.com/dairycart/dairycart/admin

ADD server .
RUN go build -o /admin-server

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

ADD templates/ /templates
ADD /dist/vendor /dist/vendor
ADD /dist/images /dist/images

COPY --from=sass-stage /app.css /dist/css/app.css
COPY --from=typescript-stage /app.js /dist/js/app.js
COPY --from=build-stage /admin-server /admin-server

EXPOSE 80:1234

ENTRYPOINT ["/admin-server"]
