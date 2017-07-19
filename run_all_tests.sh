docker build -t api_test -f api/test.Dockerfile .
docker run --name api_test --rm api_test
docker-compose --file docker-compose-test.yml up --abort-on-container-exit --build --remove-orphans --force-recreate
