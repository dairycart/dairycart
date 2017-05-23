if [ -z "$1" ]
then
    docker system prune -f && docker-compose -f docker-compose-test.yml up --abort-on-container-exit --build --remove-orphans --force-recreate
else
    docker system prune -f && docker-compose -f docker-compose-test.yml up --build --remove-orphans --force-recreate
fi