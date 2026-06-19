# Subscription
REST-сервис для агрегации данных об онлайн подписках пользователей

## План реализации
1. [x] Проектирование структуры
2. [x] Domain слой
3. [x] Use Case слой
4. [x] Repository слой
5. [x] Delivery слой
6. [x] Миграции PostgreSQL
7. [x] Конфигурация (.env)
8. [x] Swagger документация
9. [x] Docker Compose
10. [x] Логирование
11. [x] CRUDL
12. [x] Выставить HTTP-ручку для подсчета суммарной стоимости всех подписок за выбранный период с фильтрацией по `id` пользователя и названию подписки

## Документация
- [Документация API](http://localhost:8080/swagger/index.html)
- [История версий](./CHANGELOG.md)
- [План развития](./ROADMAP.md)


### Загрузка

## Docker Hub
```bash
docker pull setter2000/subscription-service
```

### Запуск сервиса из docker-compose.yml
Создайте файл `docker-compose.yml` и вставьте код ниже
```bash
services:
  postgres:
    image: postgres:15-alpine
    container_name: subscription-db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: subscription
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  app:
    image: setter2000/subscription-service:0.3.1
    container_name: subscription-service
    ports:
      - "8080:8080"
    environment:
      SERVER_PORT: 8080
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: subscription
      DB_SSL_MODE: disable
      LOG_LEVEL: info
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

volumes:
  postgres_data:

```

### Запустить docker compose
```bash
docker compose up --build
```


--- 

## GitHub
```bash
git clone git@github.com:rustruber/subscription-service.git
```

### 1. Полный запуск с миграциями
```bash
docker compose up -d --build
```

### 2. Или отдельно миграции
```bash
docker compose up migrate
```

### 3. Проверить, что таблицы создались
```bash
docker exec -it subscription-db psql -U postgres -d subscription -c "\dt"
```

### 4. Если нужно откатить
```bash
make migrate-down
```

## Пример запросов

Подсчёт суммарной стоимости всех подписок за выбранный период с фильтрацией по `id` пользователя и названию подписки

```bash
GET /subscriptions/total-cost?start_date=01-2025&end_date=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus
````

Общая сумма всех оплат
```bash
GET  /subscriptions/total-cost?start_date=01-2026&end_date=12-2026
```

### Доступные команды
```bash
make help
```

## Логи в контейнере

### Логи в реальном времени
```bash
docker compose logs -f
```

### Логи всех сервисов
```bash
docker compose logs
```
### Логи конкретного сервиса
```bash
# Логи твоего сервиса
docker compose logs app
# Логи БД
docker compose logs postgres
# Логи миграций
docker compose logs migrate
```

### Логи последние 50 строк
```bash
docker compose logs --tail 50
```