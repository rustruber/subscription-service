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