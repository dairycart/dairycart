# need to do this to ensure the database is always fresh
docker system prune -f
# need this because this actually runs the tests, duh
docker-compose --file integration_tests.yml up --abort-on-container-exit --build --remove-orphans --force-recreate
