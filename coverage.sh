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
DAIRYSECRET="do-not-use-secrets-like-this-plz" go test github.com/dairycart/dairycart/api -coverprofile=coverage.out -parallel=1
go tool cover -html=coverage.out
removeCoverageOut