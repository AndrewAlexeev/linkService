# Link Storage Service

>Сервис позволяет сохранять ссылки, получать их по короткому идентификатору и вести статистику обращений.

## 📌 Оглавление
1. [Описание](#1-описание)
2. [Установка](#2-установка)
3. [Конфигурация](#3-конфигурация)
4. [Использование](#4-использование)
5. [Технологии](#5-технологии)
6. [Авторы](#6-авторы)

## 1. Описание

Сервис предоставляет следующую функциональность:
- ✅ Создание коротких ссылок
- ✅ Получение оригинального URL по короткому коду (с автоинкрементом счетчика)
- ✅ Список всех ссылок с пагинацией
- ✅ Удаление ссылок
- ✅ Статистика по ссылкам
- ✅ Кеширование в Redis

## 2. Установка

```bash
git clone https://github.com/AndrewAlexeev/linkService
cd link-service
docker compose build
docker compose up
```

## 3. Конфигурация

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

Можно переопределить в файле `docker-compose.yml`.

## 4. Использование

### 4.1 Создание короткой ссылки

```bash
curl -X POST http://localhost:80/links \
  -H "Content-Type: application/json" \
  -d '{"url": "documents.mvideo.ru"}'
```

Успешный ответ (201 Created):

```json
{
  "short_code": "cGjK1xxjO9"
}
```
Ошибка (400 Bad Request):

```json
{
  "error": "url is required"
}
```

### 4.2 Получение ссылки по короткому коду

```bash
curl -X GET http://localhost:80/links/cGjK1xxjO9
```

Успешный ответ (200 OK):

```json
{
  "url": "documents.mvideo.ru",
  "visits": 1
}
```
Ошибка (404 Not Found):

```json
{
  "error": "Not found url in db"
}
```

### 4.3 Получение статистики по короткому коду

```bash
curl -X GET http://localhost:80/links/cGjK1xxjO9/stats
```
Успешный ответ (200 OK):

```json
{
  "short_code": "cGjK1xxjO9",
  "url": "documents.mvideo.ru",
  "visits": 1,
  "created_at": "2026-05-18T22:05:39.33773Z"
}
```
Ошибка (404 Not Found):
```json
{
  "error": "Not found url in db"
}
```

### 4.4 Удаление по короткому коду

```bash
curl -X DELETE http://localhost:80/links/cGjK1xxjO9
```
Успешный ответ (204 Not Content)

(пустое тело ответа)

### 4.5 Получение списка ссылок

```bash
curl -X GET "http://localhost:80/links?limit=10&offset=0"
```

Успешный ответ (200 OK):

```json
[
  {
    "short_code": "cGjK1xxjO9",
    "url": "documents.mvideo.ru",
    "visits": 1,
    "created_at": "2026-05-18T22:05:39.33773Z"
  },
  {
    "short_code": "eEgbozNfli",
    "url": "doc.mvieo.ru/2",
    "visits": 13,
    "created_at": "2026-05-18T21:48:06.582785Z"
  },
  {
    "short_code": "5fBC8qd5HP",
    "url": "doc.mvieo.ru/3",
    "visits": 0,
    "created_at": "2026-05-18T21:48:05.533423Z"
  }
]
```

Графический интерфейс redis доступен по адресу: http://localhost:8083

## 5. Технологии

- [Golang]
- [postgres]
- [docker]
- [redis]


## 6. Авторы

- Андрей Алексеев