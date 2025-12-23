# Практическое задание №8
## Заикин Д.Ю. ЭФМО-02-25
Работа с MongoDB: подключение, создание коллекции, CRUD-операции

<ul><strong>Цели:</strong></ul>
<li>Понять базовые принципы документной БД MongoDB (документ, коллекция, BSON, _id:ObjectID)</li>
<li>Научиться подключаться к MongoDB из Go с использованием официального драйвера</li>
<li>Создать коллекцию, индексы и реализовать CRUD для одной сущности (notes)</li>
<li>Отработать фильтрацию, пагинацию, обновления (в т.ч. частичные), удаление и обработку ошибок</li>

### Запуск
Предварительно запустить MongoDB:
```bash
docker compose up -d
```

Скопировать .env.example в .env:
```bash
cp .env.example .env
```

Зайти в папку `/pz8-mongo` и запустить сервер:
```bash
go run ./cmd/api
```

Запуск тестов:
```bash
go test ./internal/notes -v
```

### Требования
- Go ≥ 1.21
- Docker (для запуска MongoDB)
- curl или Postman для тестирования API

### Примеры запросов и ожидаемые ответы

**Создать заметку (POST /api/v1/notes)**
```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"First MongoDB Note","content":"Hello from Go and Mongo!"}'
```

Ожидаемый ответ (201 Created):
```json
{
  "id": "651fba2f9a1b2c3d4e5f6g7h",
  "title": "First MongoDB Note",
  "content": "Hello from Go and Mongo!",
  "createdAt": "2025-12-23T12:00:00Z",
  "updatedAt": "2025-12-23T12:00:00Z"
}
```

**Список с поиском (GET /api/v1/notes?q=&limit=&skip=)**
```bash
curl "http://localhost:8080/api/v1/notes?q=MongoDB&limit=5&skip=0"
```

Ожидаемый ответ (200 OK):
```json
[
  {
    "id": "651fba2f9a1b2c3d4e5f6g7h",
    "title": "First MongoDB Note",
    "content": "Hello from Go and Mongo!",
    "createdAt": "2025-12-23T12:00:00Z",
    "updatedAt": "2025-12-23T12:00:00Z"
  }
]
```

**Получить по ID (GET /api/v1/notes/{id})**
```bash
curl "http://localhost:8080/api/v1/notes/651fba2f9a1b2c3d4e5f6g7h"
```

Ожидаемый ответ (200 OK):
```json
{
  "id": "651fba2f9a1b2c3d4e5f6g7h",
  "title": "First MongoDB Note",
  "content": "Hello from Go and Mongo!",
  "createdAt": "2025-12-23T12:00:00Z",
  "updatedAt": "2025-12-23T12:00:00Z"
}
```

**Частичное обновление (PATCH /api/v1/notes/{id})**
```bash
curl -X PATCH http://localhost:8080/api/v1/notes/651fba2f9a1b2c3d4e5f6g7h \
  -H "Content-Type: application/json" \
  -d '{"content":"Updated content from API"}'
```

Ожидаемый ответ (200 OK):
```json
{
  "id": "651fba2f9a1b2c3d4e5f6g7h",
  "title": "First MongoDB Note",
  "content": "Updated content from API",
  "createdAt": "2025-12-23T12:00:00Z",
  "updatedAt": "2025-12-23T12:05:30Z"
}
```

**Удалить (DELETE /api/v1/notes/{id})**
```bash
curl -i -X DELETE http://localhost:8080/api/v1/notes/651fba2f9a1b2c3d4e5f6g7h
```

Ожидаемый ответ (204 No Content):
```
HTTP/1.1 204 No Content
```

**Статистика (GET /api/v1/notes/stats) — задание со звездочкой**
```bash
curl "http://localhost:8080/api/v1/notes/stats"
```

Ожидаемый ответ (200 OK):
```json
{
  "_id": null,
  "totalNotes": 5,
  "avgContentLength": 142.5
}
```

### Типовые проблемы и их решения

**connection refused / i/o timeout** — проверьте, что контейнер с MongoDB запущен и порт 27017 проброшен: `docker compose ps`

**Аутентификация не прошла** — убедитесь, что в MONGO_URI указаны правильные root/secret и параметр ?authSource=admin

**duplicate key error (HTTP 409)** — сработал уникальный индекс на поле title, используйте другое название заметки

**invalid ObjectID (HTTP 404)** — при неверном формате id возвращается 404, проверьте корректность ID из предыдущего запроса

**Тесты не запускаются** — убедитесь, что MongoDB запущена и доступна на localhost:27017, тесты используют отдельную БД pz8_test

### Ссылка на репозиторий
<a href="https://github.com/LilShkrek/TIP-pracs/tree/prac8">Ссылка на репозиторий с исходным кодом проекта</a>

Обязательная структура репозитория:
- cmd/api/main.go
- internal/db/mongo.go
- internal/notes/model.go
- internal/notes/repo.go
- internal/notes/repo_test.go
- internal/notes/handler.go
- docker-compose.yml
- .env.example
- README.md

# Ответы на контрольные вопросы:

- Чем документная модель MongoDB принципиально отличается от реляционной? Когда она удобнее?

MongoDB использует документы (JSON-подобные объекты) вместо таблиц с фиксированной схемой, позволяя каждому документу иметь разные поля и вложенные структуры. Реляционная БД (PostgreSQL) требует жёсткую схему с таблицами и типами данных. **Удобнее MongoDB когда:** динамическая структура данных (профили с разными атрибутами), хранение вложенных объектов (адреса, массивы тегов), быстрая разработка без миграций схемы, масштабирование через sharding. **Удобнее SQL когда:** нужны строгие ограничения целостности, сложные связи между таблицами (JOIN), стандартные ACID-транзакции.

- Что такое ObjectID и зачем нужен _id? Как корректно парсить/валидировать его в Go?

**ObjectID** — это специальный 12-байтный идентификатор, автоматически создаваемый MongoDB для каждого документа (_id). Содержит временную метку (4 байта), машинный идентификатор (3 байта), процесс ID (2 байта) и счётчик (3 байта), что гарантирует уникальность. **В Go:** используйте тип `primitive.ObjectID` из пакета `go.mongodb.org/mongo-driver/bson/primitive`. Парсинг из hex-строки: `oid, err := primitive.ObjectIDFromHex(hexString)`. Валидация: если ошибка при парсинге, возвращайте 404 (неверный ID).

- Какие операции CRUD предоставляет драйвер MongoDB и какие операторы обновления вы знаете?

**CRUD-операции:** InsertOne/InsertMany (create), FindOne/Find (read), UpdateOne/UpdateMany (update), DeleteOne/DeleteMany (delete). **Операторы обновления:** `$set` (установить значение), `$inc` (увеличить на число), `$push` (добавить в массив), `$pull` (удалить из массива), `$unset` (удалить поле), `$rename` (переименовать поле). В коде используем их через bson.M: `bson.M{"$set": bson.M{"field": value}}`.

- Как устроены индексы в MongoDB? Как создать уникальный индекс и чем он грозит при вставке?

Индексы в MongoDB ускоряют поиск по полям, аналогично SQL. **Типы:** обычный (по одному полю), составной (по нескольким), уникальный (не допускает дубликатов), TTL (удаляет документы по истечении времени). **Создание уникального индекса:** `col.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "title", Value: 1}}, Options: options.Index().SetUnique(true)})`. **Грозит при вставке:** если попытаться вставить документ с дублирующимся значением в уникальном поле, вернётся ошибка `duplicate key error` (HTTP 409 в API).

- Почему важно использовать context.WithTimeout при вызовах к базе? Что произойдет при его срабатывании?

**Важность:** без таймаута приложение может зависнуть на сетевых ошибках, недоступности БД или долгих запросах, что приводит к неотзывчивости сервера. **При срабатывании:** операция будет отменена, драйвер вернёт ошибку context.DeadlineExceeded, что мы обработаем и вернём клиенту HTTP 500. Таймауты также помогают контролировать ресурсы и предотвращать запрос от зависания на бесконечность.
