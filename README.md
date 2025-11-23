PullRequest Service. Сервис, который назначает ревьюеров на PR из команды автора, позволяет выполнять переназначение ревьюверов и получать список PR’ов, назначенных конкретному пользователю, а также управлять командами и активностью пользователей. После merge PR изменение состава ревьюверов запрещено.

Архитектура в проекте представлена следующим образом:
```bash
pr-service/
├── cmd/app/                # Точка входа   
├── internal/
│   ├── api/handlers/       # HTTP обработчики
│   ├── app/               # (Service Provider)
│   ├── dto/                # DTO (handler)
│   ├── config/            # Конфигурация (DB, HTTP)
│   ├── infrastructure/    # Инфраструктура (DB connection)
│   ├── model/             # Модели данных
│   ├── repository/        # Слой репозиториев
│   └── service/           # Слой бизнес-логики
├── pkg                # Доп.утилиты 
├── migrations/           # Файлы для миграции БД
```

Запуск приложения:
```bash
git clone https://github.com/vengeancegod/pr-service.git
cd pr-service

docker compose up -d --build
```
Создание миграций:
```bash
make migrate-up
```

На написание тестов, к сожалению, не хватило времени, поэтому API были протестированы только вручную.
Сценарии тестирования:
Создание команды:
```bash
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "payments",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": true
      },
      {
        "user_id": "u2", 
        "username": "Bob",
        "is_active": true
      },
      {
        "user_id": "u3",
        "username": "Charlie",
        "is_active": true
      },
      {
        "user_id": "u4",
        "username": "David",
        "is_active": true
      },
      {
        "user_id": "u5",
        "username": "Eva",
        "is_active": true
  }'] } "is_active": false",
  ```
Получение команды:
```bash
curl "http://localhost:8080/team/get?team_name=payments"
```

Создание PR'а:
```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add payment processing",
    "author_id": "u1"
  }'
```

Получение всех PR, где юзер назначен ревьюером:
```bash
curl "http://localhost:8080/users/getReview?user_id=u3"
```

Смена активности ревьюера:
```bash
curl -X POST http://localhost:8080/users/setIsActive \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "u3",
    "is_active": false
  }'
```

Смена ревьюера на PR'e:
```bash
curl -X POST http://localhost:8080/pullRequest/replace   -H "Content-Type: application/json"   -d '{
    "pull_request_id": "pr-1001",
    "old_user_id": "u3"
  }'
```

Мерж PR'a:
```bash
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001"
  }'
```

Попытка переназначить смерженный PR:
```bash
curl -X POST http://localhost:8080/pullRequest/replace   -H "Content-Type: application/json"   -d '{
    "pull_request_id": "pr-1001",
    "old_user_id": "u2"
  }'
```

При разработке сервиса, попытался впервые использовать кодген с помощью Ogen. Генерация файлов прошла успешно, однако возникли большие сложности с точки зрения архитектуры, поэтому решил обойтись без кодогенерации.
Разработка велась с помощью net/http в соответствии со специффикацией openapi.yml.
