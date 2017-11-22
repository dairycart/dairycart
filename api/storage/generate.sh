set -e
gnorm gen --config="gnorm.toml" # --verbose

if [ -z "$1" ]; then
    go test github.com/dairycart/dairycart/api/storage/postgres -cover
    go build -o=fart github.com/dairycart/dairycart/api/storage/mock && rm fart
else
    go build -o=fart github.com/dairycart/dairycart/api/storage/mock && rm fart
    go test github.com/dairycart/dairycart/api/storage/postgres -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out
fi