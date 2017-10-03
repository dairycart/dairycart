# sass stage
FROM ruby:latest AS sass-stage

RUN gem install sass
ADD src .
RUN sass sass/*.sass app.css

# elm stage
FROM codesimple/elm:0.18 AS elm-stage

ADD src/elm .

RUN elm-make --yes Main.elm --output elm.js

# build stage
FROM golang:alpine AS build-stage
WORKDIR /go/src/github.com/dairycart/dairycart/admin

ADD server .
RUN go build -o /admin-server

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

ADD /assets/vendor /assets/vendor
ADD /assets/images /assets/images

COPY --from=sass-stage /app.css /assets/css/app.css
COPY --from=elm-stage /elm.js /assets/js/elm.js
COPY --from=build-stage /admin-server /admin-server

ENTRYPOINT ["/admin-server"]
