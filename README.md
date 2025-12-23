# Практическое задание №10
## Заикин Д.Ю. ЭФМО-02-25
JWT-аутентификация: создание и проверка токенов. Middleware для авторизации

<ul><strong>Цели:</strong></ul>
<li>Научиться генерировать и валидировать JWT токены (HS256)</li>
<li>Реализовать middleware аутентификации (AuthN) для проверки подписи токена</li>
<li>Реализовать middleware авторизации (AuthZ) с RBAC (Role-Based Access Control)</li>
<li>Разделить access и refresh токены с разными сроками действия</li>
<li>Реализовать защищённые маршруты и обработку ошибок 401/403</li>
<li>Подготовить основу для масштабируемой аутентификации в production</li>

### Краткое описание проекта и структуры каталога

**pz10-auth** - полнофункциональный сервис JWT-аутентификации с разделением access и refresh токенов, middleware проверки подписи токена и авторизацией по ролям (RBAC). Сервис демонстрирует stateless аутентификацию, при которой информация о пользователе зашита в самом токене (sub, email, role) и защищена HS256-подписью. Благодаря этому сервер может проверить подлинность токена без обращения к БД, что упрощает масштабирование и развёртывание на кластер серверов.

**Структура проекта:**
<img width="578" height="574" alt="image" src="https://github.com/user-attachments/assets/88a0f7b8-1e4c-4160-9834-1b93bbb8ec63" />

### Установка и запуск

**1. Установить переменную окружения JWT_SECRET:**

**Linux/macOS (bash/zsh):**
```bash
export APP_PORT=8080
export JWT_SECRET="your-super-secret-key"
export JWT_TTL=24h
```

**2. Установить зависимости:**
```bash
cd pz10-auth
go mod download
```

**3. Запустить сервер:**
```bash
go run ./cmd/server
```

Сервер запустится на `http://localhost:8080`

### Примеры запросов

**Вход (POST /api/v1/login)**<br>
<img width="1413" height="597" alt="image" src="https://github.com/user-attachments/assets/a3269cc4-604b-4136-8c07-a1f2f59dff40" />

**Получение профиля (GET /api/v1/me)**<br>
<img width="983" height="576" alt="image" src="https://github.com/user-attachments/assets/3e1caeda-7bb1-46e7-b9dc-0d92d3b36ed1" />

**Доступ без токена (GET /api/v1/me)**<br>
<img width="983" height="576" alt="image" src="https://github.com/user-attachments/assets/ec5653f3-c357-41e4-b08f-3b7e85bfc9fe" />

**Запрос админ-статистики (GET /api/v1/admin/stats) для admin**<br>
<img width="989" height="495" alt="image" src="https://github.com/user-attachments/assets/9cc978f4-26e9-4b0e-b986-e620a644c8ac" />

**Попытка доступа как user**<br>
<img width="989" height="495" alt="image" src="https://github.com/user-attachments/assets/9929b933-10c4-4d58-af1a-af1c913416dd" />

**Обновление access токена**<br>
<img width="1419" height="610" alt="image" src="https://github.com/user-attachments/assets/95474b54-842a-4c6d-a9af-b3846c298a7c" />

**Получение пользователя по id**<br>
<img width="1419" height="610" alt="image" src="https://github.com/user-attachments/assets/c4532bf2-1aad-4f34-9737-618bafea8f21" />

### Тестовые пользователи

| Email | Пароль | Роль |
|-------|--------|------|
| admin@example.com | secret123 | admin |
| user@example.com | secret123 | user |

### Ключевые фрагменты кода

**Генерация JWT (internal/platform/jwt/jwt.go):**
```go
func (h *HS256) Sign(userID int64, email, role string) (string, error) {
    now := time.Now()
    claims := jwt.MapClaims{
       "sub":   userID,
       "email": email,
       "role":  role,
       "iat":   now.Unix(),
       "exp":   now.Add(h.ttl).Unix(),
       "iss":   "pz10-auth",
       "aud":   "pz10-clients",
       "type":  "access",
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return t.SignedString(h.secret)
}
```

**Middleware аутентификации (internal/http/middleware/authn.go):**
```go
func AuthN(v jwt.Validator) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          h := r.Header.Get("Authorization")
          if h == "" || !strings.HasPrefix(h, "Bearer ") {
             http.Error(w, "unauthorized", http.StatusUnauthorized)
             return
          }
          raw := strings.TrimPrefix(h, "Bearer ")
          claims, err := v.Parse(raw)
          if err != nil {
             http.Error(w, "unauthorized", http.StatusUnauthorized)
             return
          }
          ctx := context.WithValue(r.Context(), CtxClaimsKey, claims)
          next.ServeHTTP(w, r.WithContext(ctx))
       })
    }
}
```

**Middleware авторизации (internal/http/middleware/authz.go):**
```go
func AuthZRoles(allowed ...string) func(http.Handler) http.Handler {
    set := map[string]struct{}{}
    for _, a := range allowed {
       set[a] = struct{}{}
    }
    return func(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          claims := r.Context().Value(CtxClaimsKey).(map[string]any)
          role, _ := claims["role"].(string)
          if _, ok := set[role]; !ok {
             http.Error(w, "forbidden", http.StatusForbidden)
             return
          }
          next.ServeHTTP(w, r)
       })
    }
}
```

**ABAC в GetUserHandler (internal/core/service.go):**
```go
func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(middleware.CtxClaimsKey).(map[string]any)
    userRole, _ := claims["role"].(string)
    userID, _ := claims["sub"].(float64)

    requestedID := chi.URLParam(r, "id")

    // ABAC: user может видеть только свой профиль
    if userRole == "user" {
       currentUserID := strconv.Itoa(int(userID))
       if requestedID != currentUserID {
          http.Error(w, "access_denied", http.StatusForbidden)
          return
       }
    }

    jsonOK(w, map[string]any{
       "id":    requestedID,
       "email": "user" + requestedID + "@example.com",
       "role":  "user",
    })
}
```

**Refresh-токен (internal/platform/jwt/jwt.go):**
```go
func (h *HS256) SignRefresh(userID int64) (string, error) {
    now := time.Now()
    claims := jwt.MapClaims{
       "sub":  userID,
       "iat":  now.Unix(),
       "exp":  now.Add(7 * 24 * time.Hour).Unix(),
       "iss":  "pz10-auth",
       "aud":  "pz10-clients",
       "type": "refresh",
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return t.SignedString(h.secret)
}
```

### Маршруты и статусы

| Метод | Маршрут | Описание | Статусы |
|-------|---------|---------|---------|
| POST | /api/v1/login | Вход, выдача access/refresh токенов | 200, 401 |
| POST | /api/v1/refresh | Обновление access токена | 200, 401 |
| GET | /api/v1/me | Профиль текущего пользователя | 200, 401 |
| GET | /api/v1/admin/stats | Статистика (только admin) | 200, 401, 403 |
| GET | /api/v1/users/{id} | Профиль пользователя (ABAC) | 200, 401, 403 |

### Переменные окружения

| Переменная | Значение | Описание |
|-----------|----------|---------|
| APP_PORT | 8080 | Порт для запуска сервера |
| JWT_SECRET | (обязательно) | Секретный ключ для подписи JWT (минимум 32 символа) |
| JWT_TTL | 24h | Срок действия access-токена |

### Ссылка на репозиторий
<a href="https://github.com/LilShkrek/TIP-pracs/tree/prac10">Ссылка на репозиторий с исходным кодом проекта</a>

Обязательная структура репозитория:
- cmd/server/main.go
- internal/core/user.go
- internal/core/service.go
- internal/http/router.go
- internal/http/middleware/authn.go
- internal/http/middleware/authz.go
- internal/repo/user_mem.go
- internal/platform/jwt/jwt.go
- internal/platform/config/config.go
- go.mod
- README.md

---

# Ответы на контрольные вопросы:

### 1. Что такое клеймы JWT и чем отличаются registered, public, private? Почему важно exp?

**Клеймы** — утверждения о пользователе в payload токена (`sub`, `email`, `role`, `iat`, `exp`, `iss`, `aud`). **Registered claims** — стандартные (RFC 7519): `iss`, `sub`, `aud`, `exp`, `nbf`, `iat`, `jti`. **Public claims** — кастомные общего назначения (`email`, `role`). **Private claims** — специфичные для приложения (`type: "access"`). **Почему exp важен:** это deadline токена. Если украден → действует ограниченное время. Без exp токен валиден вечно.

### 2. Чем stateless-аутентификация на JWT отличается от сессионных cookie на сервере? Плюсы/минусы.

**JWT (stateless):** информация в токене, сервер не хранит состояние. Любой сервер проверит токен. **Плюсы:** масштабирование, меньше нагрузка на сервер, микросервисы. **Минусы:** сложнее отозвать (blacklist), payload открыт, размер больше. **Sessions (stateful):** ID в cookie, данные на сервере. **Плюсы:** легко отозвать, меньше по сети. **Минусы:** синхронизация сессий, БД bottleneck. **Вывод:** JWT для REST/микросервисов, sessions для веб-приложений.

### 3. Как устроена цепочка middleware и почему AuthZ должна идти после AuthN?

**Цепь:** Request → AuthN → AuthZ → Handler → Response. **AuthN** проверяет подпись, извлекает claims в контекст, возвращает 401 если невалиден. **AuthZ** читает claims из контекста, проверяет роль, возвращает 403 если недостаточно. **Почему после:** 1) AuthN готовит claims для AuthZ. 2) Нет смысла проверять роль без аутентификации. 3) 401 (не аутентифицирован) vs 403 (недостаточно прав).

### 4. RBAC vs ABAC: когда что выбирать? Примеры.

**RBAC:** контроль только по роли (admin, user). Простая логика. Пример: "/admin/stats только для admin". **ABAC:** по множеству атрибутов (role, owner, IP, время). Сложнее. Пример: "user видит только свой профиль". **Выбор:** RBAC для простых сценариев (несколько ролей, все admin могут всё). ABAC для гранулярного контроля (SaaS, multi-tenant). **В ПЗ10:** RBAC в middleware `AuthZRoles()`, ABAC в handler `GetUserHandler` (проверка владельца).

### 5. Как безопасно хранить пароль и почему bcrypt/argon2, а не SHA-256?

**Не SHA-256:** 1) Очень быстрый → перебор паролей миллиардами в секунду. 2) Два пароля → один хэш → rainbow tables. **bcrypt:** 1) Медленный (250ms/хэш) → brute-force нецелесообразен. 2) Встроённая соль → каждый пароль уникален. 3
