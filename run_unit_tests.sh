docker system prune -f
docker build -t api_test -f test.Dockerfile .
docker run --name api_test --rm api_test