# sass stage
FROM ruby:latest AS sass-stage

RUN gem install sass
ADD src .
RUN sass sass/*.sass app.css

# node stage
FROM node:latest AS node-stage

ADD src/javascript build
WORKDIR build

RUN npm install
RUN npm run docker-build

# build stage
FROM golang:alpine AS build-stage
WORKDIR /go/src/github.com/dairycart/dairycart/admin

ADD server .
RUN go build -o /admin-server

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

ADD server/html /html
ADD server/html /html
ADD /assets/vendor /assets/vendor
ADD /assets/images /assets/images

COPY --from=sass-stage /app.css /assets/css/app.css
COPY --from=node-stage /build/output/js/app.js /assets/js/app.js
COPY --from=node-stage /build/output/js/app.js.map /assets/js/app.js.map
COPY --from=build-stage /admin-server /admin-server

ENTRYPOINT ["/admin-server"]
