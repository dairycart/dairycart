# sass stage
FROM ruby:latest AS sass-stage

RUN gem install sass
ADD src .
RUN sass sass/*.sass app.css

# typescript stage
# FROM node:latest AS typescript-stage
FROM sandrokeil/typescript:latest as typescript-stage

ADD src .
# RUN npm install -g typescript
RUN tsc typescript/*.ts --outFile app.js

# build stage
FROM golang:alpine AS build-stage
WORKDIR /go/src/github.com/dairycart/dairycart/admin

ADD server .
RUN go build -o /admin-server

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY templates/ /templates
COPY /dist/vendor /dist/vendor
COPY /dist/images /dist/images

COPY --from=sass-stage /app.css /dist/css/app.css
COPY --from=typescript-stage /app.js /dist/js/app.js

COPY --from=build-stage /admin-server /admin-server

ENTRYPOINT ["/admin-server"]
