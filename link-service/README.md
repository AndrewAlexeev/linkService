# Link Storage Service

>Сервис позволяет сохранять ссылки, получать их по короткому идентификатору и вести статистику обращений.

## 📌 Оглавление
1. [Описание](#описание)
2. [Установка](#установка)
3. [Конфигурация](#конфигурация)
4. [Использование](#использование)
5. [Технологии](#технологии)
6. [Авторы](#авторы)

1. ## Описание
Сервис предоставляет следующую функциональность:
- ✅ Создание коротких ссылок
- ✅ Получение оригинального URL по короткому коду (с автоинкрементом счетчика)
- ✅ Список всех ссылок с пагинацией
- ✅ Удаление ссылок
- ✅ Статистика по ссылкам
- ✅ Кеширование в Redis
2. ## Установка
```bash
git clone https://github.com/AndrewAlexeev/linkService
cd link-service
docker compose build
docker compose up
```

3. ## Конфигурация

### Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|-----------|----------|----------------------|
| `PORT` | Порт сервера | `8080` |
| `DB_HOST` | Хост PostgreSQL | `postgres` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `password` |
| `DB_NAME` | Имя БД | `link_storage` |
| `REDIS_ADDR` | Адрес Redis | `redis:6379` |
| `REDIS_PASSWORD` | Пароль Redis | `""` |
| `CACHE_TTL` | Время жизни кэша (сек) | `3600` |

Можно переопределить в файле  `docker-compose.yml`.

4. ## Использование
### 4.1 Создание короткой ссылки

```bash
curl -X POST http://localhost:80/links \
  -H "Content-Type: application/json" \
  -d '{"url" : "documents.mvideo.ru"}'
```

```json
Успешный ответ (201 Created):

{
  "short_code": "cGjK1xxjO9"
}

Ошибка (400 Bad Request):

{
  "error": "url is required"
}
```

### 4.2 Получение ссылки по короткому коду

```bash
curl -X GET http://localhost:8080/links/cGjK1xxjO9

```

```json
Успешный ответ (200 Ok):


{
  "url": "documents.mvideo.ru",
  "visits": 1
}

Ошибка (404 Not Found):

{
  "error": "error info: Not found url in db"
}
```

### 4.3 Получение статистики по короткому коду

```bash
curl -X GET http://localhost:80/links/cGjK1xxjO9/stats

```

```json
Успешный ответ (200 Ok):

{
  "short_code": "cGjK1xxjO9",
  "url": "documents.mvideo.ru",
  "visits": 1,
  "created_at": "2026-05-18T22:05:39.33773Z"
}

Ошибка (404 Not Found):

{
  "error": "error info: Not found url in db"
}
```

### 4.4 Удаление по короткому коду

```bash
curl -X DELETE http://localhost:80/links/cGjK1xпxjO9
```

```json
Успешный ответ (204 Not Content):

{
  "short_code": "cGjK1xxjO9",
  "url": "documents.mvideo.ru",
  "visits": 1,
  "created_at": "2026-05-18T22:05:39.33773Z"
}
```

### 4.5 Получение списка ссылок

curl -X GET "http://localhost:8080/links?limit=10&offset=0"

```json

Успешный ответ (200 Ok):

[
  {
    "ShortCode": "cGjK1xxjO9",
    "Url": "documents.mvideo.ru",
    "Visits": 1,
    "CreatedAt": "2026-05-18T22:05:39.33773Z"
  },
  {
    "ShortCode": "eEgbozNfli",
    "Url": "doc.mvieo.ru/2",
    "Visits": 13,
    "CreatedAt": "2026-05-18T21:48:06.582785Z"
  },
  {
    "ShortCode": "5fBC8qd5HP",
    "Url": "doc.mvieo.ru/new/2",
    "Visits": 0,
    "CreatedAt": "2026-05-18T21:48:05.533423Z"
  }
]
```

Графический интерфейс redis доступен по адресу: http://localhost:8083

5. ## Технологии
- [Golang]
- [postgres]
- [docker]
- [redis]


6. ## Авторы
- Андрей Алексеев
