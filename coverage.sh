if [ -f coverage.out ]
then
    rm coverage.out
fi

ssh-keygen -t rsa -b 4096 -f api/app.rsa -N ''
openssl rsa -in  api/app.rsa -pubout -outform PEM -out api/app.rsa.pub

go test github.com/verygoodsoftwarenotvirus/dairycart/api -coverprofile=coverage.out -tags test
go tool cover -html=coverage.out
rm coverage.out

rm api/app.rsa*