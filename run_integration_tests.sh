docker system prune -f && docker-compose -f integration-tests-compose.yml up --abort-on-container-exit --build --remove-orphans --force-recreate
