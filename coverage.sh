if [ -f coverage.out ]
then
    rm coverage.out
fi
go test github.com/verygoodsoftwarenotvirus/dairycart/api -coverprofile=coverage.out
go tool cover -html=coverage.out
rm coverage.out