# delete previous coverage report, if it exists
if [ -f admin-server ]
then
    rm admin-server
fi

(cd server && go build -o ../admin-server)
./admin-server