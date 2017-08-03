docker login --username=$DOCKER_USER --password=$DOCKER_PASS
(cd api; docker build --no-cache --tag dairycart/api-server:pre-alpha .)
docker push dairycart/api-server:pre-alpha