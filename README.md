## Subscriptions Service

Сервис для управления подписками (создание, обновление, удаление, суммирование) с использованием Go, PostgreSQL и Docker.

### Функционал

REST API для управления подписками:

POST /create_sub — создать подписку

GET /get_all_subs — получить все подписки

PUT /update_sub — обновить подписку

DELETE /delete_sub — удалить подписку

POST /sum_subs — суммировать стоимость подписок по диапазону


### Установка и запуск

1. Клонируйте репозиторий
```plaitext
git clone <url_repo>
cd subs
```

2. Docker-compose запуск
```plaintext
docker docker compose up --build
```

Swagger документация доступна по:

http://localhost:8080/swagger/index.html

После запуска:

API: http://localhost:8080

Swagger: http://localhost:8080/swagger/index.html

База PostgreSQL на порту 5444 (хост) или 5432 (сеть Docker)


### Технологии
Go 1.24

PostgreSQL 16

Docker / Docker Compose

Swagger для документации
