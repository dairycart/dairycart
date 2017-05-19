echo "should print something like:
HTTP/1.1 200 OK
Date: Fri, 19 May 2017 03:06:15 GMT
Content-Type: text/plain; charset=utf-8

=====================HERE WE GO=====================
"

curl --head "localhost:8080/product/skateboard"