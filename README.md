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
### Доступные команды
```bash
make help
```
