set -e

(cd migrations && go-bindata -nocompress -pkg migrations .)

gnorm gen # --verbose
if [ -z "$1" ]; then
    go test github.com/dairycart/dairycart/storage/database/postgres -cover
else
    go test github.com/dairycart/dairycart/storage/database/postgres -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out
fi