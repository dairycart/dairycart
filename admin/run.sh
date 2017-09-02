# docker system prune -f
# docker build -t admin_panel .
# docker run -i -t admin_panel -p 80:80 --rm

# docker-compose --file docker-compose.yml up --abort-on-container-exit --build --remove-orphans --force-recreate
docker-compose --file docker-compose.yml up --abort-on-container-exit --build --remove-orphans