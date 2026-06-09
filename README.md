
# Subscription Service API

REST API сервис для агрегации данных об онлайн-подписках пользователей.

---

##  Функциональность

Сервис предоставляет:

### CRUDL операции с подписками:
- Создание подписки
- Получение списка всех подписок
- Получение подписки по `user_id + service_name`
- Обновление подписки
- Удаление подписки

### Агрегация:
- Подсчёт суммарной стоимости подписок за период
- Фильтрация по:
  - `user_id`
  - `service_name`
  - диапазону дат

---

## Стек технологий

- Go (Gin)
- PostgreSQL
- Swagger (swaggo)
- Docker / Docker Compose

---

## Структура проекта

```

internal/
├── handler       # HTTP handlers 
├── service       # бизнес-логика
├── repository    # работа с БД
├── storage       # подключение к PostgreSQL
├── server        # настройка HTTP сервера и роутинг
├── models        # DTO / модели

````

---

## Конфигурация

Конфигурация задаётся через YAML файл.

Пример `config.yaml`:

```yaml
env: "development"

http_server:
  host: "0.0.0.0"
  port: ":8080"

database:
  host: "postgres"
  port: "5432"
  user: "postgres"
  password: "postgres"
  name: "subscriptions"
  sslmode: "disable"
````

Переменная окружения:

```bash
CONFIG_PATH=./configs/config.yaml
```

---

## Запуск через Docker

```bash
docker-compose up --build
```

---

## Миграции БД

SQL файл инициализации:

```sql
CREATE TABLE subscriptions (
    user_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),
    start_date VARCHAR(7) NOT NULL,
    end_date VARCHAR(7),
    PRIMARY KEY (user_id, service_name)
);
```

---

##  API Endpoints

### Healthcheck

```
GET /health
```

---

### Subscriptions

#### Создать подписку

```
POST /subscriptions
```

#### Получить все подписки

```
GET /subscriptions
```

#### Получить подписку

```
GET /subscriptions/item?user_id=&service_name=
```

#### Обновить подписку

```
PUT /subscriptions/item
```

#### Удалить подписку

```
DELETE /subscriptions/item
```

---

### Агрегация

#### Сумма подписок

```
GET /subscriptions/aggregate?user_id=&service_name=&start_date=&end_date=
```

---

## Пример запроса

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
```

---


