docker build -t plugins --file plugin.Dockerfile .
docker run --volume=$GOPATH/src/github.com/dairycart/dairycart/api/example_files/plugins:/output --rm -t plugins
