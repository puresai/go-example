# docker-compose

```
docker pull redis:5.0

docker run -itd --name redis5 -p 6389:6379 redis:5.0

docker build -t mysql-for-user ./docker/mysql/

docker run  -itd --name mysql-for-user -p 3316:3306 -e MYSQL_ROOT_PASSWORD=111111 mysql-for-user

```