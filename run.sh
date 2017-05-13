if [ -f dairycart ]
then
    rm dairycart
fi
go build && clear && ./dairycart
