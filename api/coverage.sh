if [ ! -e coverage.out ] ; then
    touch coverage.out
fi

go test -coverprofile=coverage.out
go tool cover -html=coverage.out