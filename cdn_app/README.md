
Миграции делаем с помощью пакета goose (https://github.com/pressly/goose)

Предварительно должны быть созданы базы данных: cdn, cdn_test
CREATE DATABASE cdn;
CREATE DATABASE cdn_test;

Находясь в папке cdn_app выполнить команду:
goose -dir migrations postgres "postgres://postgres:postgres@localhost/cdn?sslmode=disable" up
goose -dir migrations postgres "postgres://postgres:postgres@localhost/cdn_test?sslmode=disable" up

Примеры запросов
curl http://localhost:8084/api/files/user/user1
["file1.mp4","file2.mp4"]

curl http://localhost:8084/api/files/user/user2
["file2.mp4"]

curl http://localhost:8084/api/files/server/srv1
["file1.mp4","file2.mp4"]

curl http://localhost:8084/api/files/server/srv3
["file1.mp4"]

http://localhost:8084/api/files/area/msk.ru
["file1.mp4","file2.mp4"]

http://localhost:8084/api/files/area/spb.ru
["file1.mp4"]


Тест написан в cdn_app/tests/api_test.go пройдет только если создана база cdn_test и dsn такой
dsn := "host=localhost user=postgres password=postgres dbname=cdn_test port=5432 TimeZone=Europe/Moscow"
