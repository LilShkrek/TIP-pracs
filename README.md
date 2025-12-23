# Практическое задание №9
## Заикин Д.Ю. ЭФМО-02-25
Реализация регистрации и входа пользователей. Хэширование паролей с bcrypt

<ul><strong>Цели:</strong></ul>
<li>Научиться безопасно хранить пароли (bcrypt), валидировать вход и обрабатывать ошибки</li>
<li>Реализовать эндпоинты POST /auth/register и POST /auth/login</li>
<li>Закрепить работу с БД (PostgreSQL + GORM) и валидацией ввода</li>
<li>Подготовить основу для JWT-аутентификации в следующем ПЗ (№10)</li>

### Краткое описание проекта и структуры каталога

**pz9-auth** - минимальный сервис аутентификации с регистрацией и входом пользователя, использующий bcrypt для безопасного хранения паролей в PostgreSQL.

**Структура проекта:**<br>
<img width="464" height="483" alt="image" src="https://github.com/user-attachments/assets/ae5c8399-4d14-4f4c-9b96-20bb876cdb99" />

### Установка и запуск

**1. Создать БД PostgreSQL:**
```bash
psql -U postgres -c "CREATE DATABASE pz9;"
```

**2. Клонировать/скопировать проект и установить зависимости:**
```bash
cd pz9-auth
go mod download
```

**3. Установить переменные окружения:**

```bash
export DB_DSN="postgres://postgres:password@localhost:5432/pz9?sslmode=disable"
export BCRYPT_COST="12"
export APP_ADDR=":8080"
```

**4. Запустить сервер:**
```bash
go run ./cmd/api
```

Сервер запустится на `http://localhost:8080`

### Примеры запросов

**Успешная регистрация (POST /auth/register)**<br>
<img width="499" height="543" alt="image" src="https://github.com/user-attachments/assets/889d6523-dcdd-4497-85c4-367ff6a9d10b" />

**Попытка повторной регистрации (POST /auth/register)**<br>
<img width="1222" height="472" alt="image" src="https://github.com/user-attachments/assets/68485099-7e8b-4095-beb3-ffb451aa1142" />

**Вход с верными данными (POST /auth/login)**<br>
<img width="718" height="532" alt="image" src="https://github.com/user-attachments/assets/23de3688-72db-4d13-b42d-41e3b6315175" />


**Вход с неверным паролем (POST /auth/login)**<br>
<img width="718" height="532" alt="image" src="https://github.com/user-attachments/assets/e1c91e2e-8f77-4d47-89fb-38ebe2d72202" />

**Вход с неверным email (POST /auth/login)**<br>
<img width="718" height="532" alt="image" src="https://github.com/user-attachments/assets/e4610faa-b7e9-4a3d-9974-d026c1f2bc9a" />

### Ключевые фрагменты кода

**Хэширование пароля при регистрации (internal/http/handlers/auth.go):**
```go
// bcrypt hash
hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), h.BcryptCost)
if err != nil {
    writeErr(w, http.StatusInternalServerError, "hash_failed")
    return
}

u := core.User{Email: in.Email, PasswordHash: string(hash)}
if err := h.Users.Create(r.Context(), &u); err != nil {
    if err == repo.ErrEmailTaken {
        writeErr(w, http.StatusConflict, "email_taken")
        return
    }
    writeErr(w, http.StatusInternalServerError, "db_error")
    return
}
```

**Проверка пароля при логине (internal/http/handlers/auth.go):**
```go
u, err := h.Users.ByEmail(context.Background(), in.Email)
if err != nil {
    // Не раскрываем, что именно не так
    writeErr(w, http.StatusUnauthorized, "invalid_credentials")
    return
}

if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
    writeErr(w, http.StatusUnauthorized, "invalid_credentials")
    return
}
```

**Модель User (internal/core/user.go):**
```go
type User struct {
    ID           int64     `gorm:"primaryKey" json:"id"`
    Email        string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
    PasswordHash string    `gorm:"size:255;not null" json:"-"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}
```

**Миграция и CRUD (internal/repo/user_repo.go):**
```go
// Создание таблицы (выполняется автоматически при запуске)
func (r *UserRepo) AutoMigrate() error {
    return r.db.AutoMigrate(&core.User{})
}

// Создание пользователя
func (r *UserRepo) Create(ctx context.Context, u *core.User) error {
    if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
        if errors.Is(err, gorm.ErrDuplicatedKey) {
            return ErrEmailTaken
        }
        return err
    }
    return nil
}

// Поиск по email
func (r *UserRepo) ByEmail(ctx context.Context, email string) (core.User, error) {
    var u core.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return core.User{}, ErrUserNotFound
    }
    return u, err
}
```

### SQL / Миграции

GORM автоматически создаёт таблицу `users` при первом запуске (`AutoMigrate()`). Эквивалент SQL:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Переменные окружения

| Переменная | Значение | Описание |
|-----------|----------|---------|
| DB_DSN | `postgres://postgres:password@localhost:5432/pz9?sslmode=disable` | Строка подключения к PostgreSQL |
| BCRYPT_COST | `12` | Сложность хэширования bcrypt (2^12 = 4096 итераций) |
| APP_ADDR | `:8080` | Адрес для запуска сервера |

### Краткие выводы

**Почему нельзя хранить пароли в открытом виде:**
- Утечка БД даёт полный доступ ко всем аккаунтам
- Пользователи используют один пароль на разных сервисах
- Нарушение базовой безопасности и доверия

**Почему bcrypt, а не SHA-256:**
- SHA-256 слишком быстрый - перебор паролей возможен миллионами в секунду
- bcrypt имеет встроенную соль - разные хэши для одного пароля
- bcrypt имеет регулируемую сложность (cost) - можно увеличить защиту без изменения кода
- bcrypt специально разработан для паролей, SHA-256 - для целостности данных

### Ссылка на репозиторий
<a href="https://github.com/LilShkrek/TIP-pracs/tree/prac9">Ссылка на репозиторий с исходным кодом проекта</a>

# Ответы на контрольные вопросы:

### 1. В чём разница между хранением пароля и хранением его хэша? Зачем соль? Почему bcrypt, а не SHA-256?

**Хранение пароля в открытом виде:** уязвимость, если БД украдена - все пароли скомпрометированы. **Хранение хэша:** одностороннее преобразование - пароль невозможно восстановить из хэша (или крайне затратно).

**Соль (salt):** случайная строка, добавляемая перед хэшированием. Без соли два одинаковых пароля дают одинаковый хэш, что позволяет атакующему использовать rainbow tables (предвычисленные таблицы хэшей). С солью каждый пароль имеет уникальный хэш даже при совпадении самого пароля. bcrypt встраивает соль в итоговый хэш автоматически.

**Почему bcrypt, а не SHA-256:** SHA-256 очень быстрый (~миллиарды хэшей в секунду на GPU), что позволяет перебирать пароли огромными скоростями. bcrypt специально разработан для паролей - он медленный (регулируемый через параметр cost, обычно 2^12 = 4096 итераций), что делает brute-force атаки экономически нецелесообразными. Кроме того, bcrypt имеет встроенную соль и правильно обрабатывает проблемы безопасности, тогда как SHA-256 требует ручного управления солью и уязвим для атак.

### 2. Что произойдёт при снижении/повышении cost у bcrypt? Как подобрать значение?

**Снижение cost (например, с 12 на 10):** хэширование ускоряется примерно в 4 раза (cost = log₂(итерации)), но безопасность падает - перебор паролей становится быстрее. **Повышение cost (например, с 12 на 14):** хэширование замедляется в 4 раза, безопасность растёт, но регистрация/логин будет медленнее.

**Подбор значения:** обычно используют 10–14. Для учебных целей и dev-машин - 12 (примерно 250ms на хэширование). Для production на мощных серверах можно 13–14 (время в пределах 500ms–1s). Нужен баланс: достаточно медленно, чтобы защитить от атак, но достаточно быстро, чтобы не делать UX неприятным.

### 3. Какие статусы и ответы должны возвращать POST /auth/register и POST /auth/login в типичных сценариях?

**POST /auth/register:**
- **201 Created** - успешная регистрация, возвращаем `{status: "ok", user: {id, email}}`
- **400 Bad Request** - email пустой, пароль <8 символов или невалидный JSON
- **409 Conflict** - email уже зарегистрирован, ошибка `email_taken`
- **500 Internal Server Error** - ошибка БД или хэширования

**POST /auth/login:**
- **200 OK** - успешный вход, возвращаем `{status: "ok", user: {id, email}}`
- **400 Bad Request** - email или пароль пустой, невалидный JSON
- **401 Unauthorized** - email не найден или пароль неверный (одно сообщение `invalid_credentials` для обоих случаев!)
- **500 Internal Server Error** - ошибка БД

### 4. Какие риски несут подробные сообщения об ошибках при логине?

**Риск enumeration attack (перебор email-адресов):** если сервер отвечает «email найден, но пароль неверный», то атакующий может перебирать email-адреса и узнавать, какие реально зарегистрированы. Затем он сосредотачивает атаку на известные email-адреса.

**Правильный подход:** всегда возвращаем одно сообщение `invalid_credentials` без уточнений. Это скрывает информацию о наличии/отсутствии email в системе.

**Дополнительные риски:** логирование паролей (даже случайное), отправка пароля в ответе или логах, отсутствие rate limiting (позволяет быстрый brute-force).

### 5. Почему в этом ПЗ не выдаём токен, и что изменится в ПЗ10 (JWT)?

**В ПЗ9:** сосредотачиваемся на безопасном хранении пароля и валидации при логине. После успешного логина просто возвращаем информацию пользователя (без токена) - это демонстрирует понимание основ аутентификации.

**В ПЗ10 (JWT):** после успешной проверки пароля генерируем JWT (JSON Web Token) - подписанный токен, который клиент отправляет в каждом последующем запросе в заголовке `Authorization: Bearer <token>`. JWT содержит информацию о пользователе (id, email) и имеет срок действия (exp). На сервере мы проверяем подпись и не требуем повторной проверки пароля - это **аутентификация** (ПЗ9) плюс **авторизация** (ПЗ10).

Таким образом: ПЗ9 = логин + проверка пароля, ПЗ10 = логин + выдача токена для последующих запросов.
