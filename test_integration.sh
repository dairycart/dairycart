docker system prune -f
docker-compose --file docker-compose-test.yml up --abort-on-container-exit --build --remove-orphans --force-recreate
