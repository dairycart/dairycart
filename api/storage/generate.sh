set -e

find mock -type f ! -iname "main.go" -delete
find postgres -type f ! -iname "helpers_test.go" ! -iname "main.go" -delete
find models -type f ! -iname "helper_types.go" -delete
rm storage.go
gnorm gen --config="gnorm.toml" # --verbose

if [ -z "$1" ]; then
    (cd postgres && go test -cover)
else
    (cd postgres && go test -cover -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out)
fi