# docker run --name admin_panel --rm admin_panel --port 80:3000
docker system prune -f
docker-compose --file docker-compose.yml up --abort-on-container-exit --build --remove-orphans --force-recreate