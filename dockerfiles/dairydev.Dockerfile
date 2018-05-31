FROM golang:1.10.2-stretch

RUN set -ex; \
    go get -u -v \
    gnorm.org/gnorm \
    github.com/rakyll/statik