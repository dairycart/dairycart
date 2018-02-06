FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairycart

ADD . .

CMD go build -buildmode=plugin -o /output/mock_db.so github.com/dairycart/dairycart/storage/database/mock/plugin; go build -buildmode=plugin -o /output/mock_img.so github.com/dairycart/dairycart/storage/images/mock/plugin
