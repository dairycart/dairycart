# delete previous coverage report, if it exists
if [ -f coverage.out ]
then
    rm coverage.out
fi

# run tests
DAIRYSECRET="do-not-use-secrets-like-this-plz" go test github.com/verygoodsoftwarenotvirus/dairycart/api -coverprofile=coverage.out -tags test
go tool cover -html=coverage.out

# delete the new coverage report so I don't accidentally commit it to the repo somehow
rm coverage.out