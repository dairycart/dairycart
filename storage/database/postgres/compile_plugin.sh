docker build -t plugins --file plugin.Dockerfile .
docker run --volume=$(pwd):/output --rm -t plugins