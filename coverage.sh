#!/usr/bin/env bash

# exit if anything fails
set -e

function removeCoverageOut () {
    if [ -f coverage.out ]
    then
        rm coverage.out
    fi
}

removeCoverageOut
docker build -t dairycoverage --file coverage.Dockerfile .
docker run --volume=$GOPATH/src/github.com/dairycart/dairycart:/output --rm -t dairycoverage
go tool cover -html=coverage.out
removeCoverageOut
