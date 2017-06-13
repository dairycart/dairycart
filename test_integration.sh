docker system prune -f
docker-compose -f docker-compose-test.yml up --abort-on-container-exit --build --remove-orphans --force-recreate
