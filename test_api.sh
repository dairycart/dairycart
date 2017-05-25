docker build -t api_test -f api_test.Dockerfile .
docker run --name api_test --rm api_test