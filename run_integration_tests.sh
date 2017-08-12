docker system prune -f
docker-compose --file api/integration_tests.yml up --abort-on-container-exit --build --remove-orphans --force-recreate
