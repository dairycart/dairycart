set -e

(cd migrations && go-bindata -nocompress -pkg migrations .)

gnorm gen # --verbose
if [ -z "$1" ]; then
    go test github.com/dairycart/postgres -cover
else
    go test github.com/dairycart/postgres -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out
fi