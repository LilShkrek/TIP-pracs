# Практическое занятие №4
## Заикин Д.Ю. ЭФМО-02-25
### 25.11.2025
Маршрутизация с chi (альтернатива — gorilla/mux). Создание небольшого CRUD-сервиса «Список задач».<br>
<ul><strong>Цели:</strong></ul>
<li>Освоить базовую маршрутизацию HTTP-запросов в Go на примере роутера chi</li>
<li>Научиться строить REST-маршруты и обрабатывать методы GET/POST/PUT/DELETE</li>
<li>Реализовать небольшой CRUD-сервис «ToDo» (без БД, хранение в памяти)</li>
<li>Добавить простое middleware (логирование, CORS)</li>
<li>Научиться тестировать API запросами через curl/Postman/HTTPie</li>

### Запуск
Зайти в папку /pz4-todo и выполнить ```go run .```

### Структура проекта
<img width="393" height="433" alt="image" src="https://github.com/user-attachments/assets/3a9df497-7aa5-4c30-a86e-eef215e90666" />

## Фрагменты кода
### Роутер (main.go)
```go
func main() {
	repo, err := task.NewRepoWithFile("tasks.json")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	h := task.NewHandler(repo)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(myMW.Logger)
	r.Use(myMW.SimpleCORS)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/tasks", h.Routes())
		})
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
```

### Middleware логирования
```go
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
```

### Обработчик с пагинацией и валидацией
```go
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)
	doneFilter := parseDoneFilter(r)

	// Получаем все задачи
	allTasks := h.repo.List()

	filteredTasks := filterTasksByDone(allTasks, doneFilter)

	// Вычисляем offset и limit для пагинации (уже для отфильтрованных задач)
	offset := (page - 1) * limit
	if offset > len(filteredTasks) {
		offset = len(filteredTasks)
	}

	end := offset + limit
	if end > len(filteredTasks) {
		end = len(filteredTasks)
	}

	pagedTasks := filteredTasks[offset:end]

	response := map[string]interface{}{
		"tasks": pagedTasks,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      len(filteredTasks),
			"totalPages": (len(filteredTasks) + limit - 1) / limit,
		},
		"filters": map[string]interface{}{
			"done": doneFilter,
		},
	}

	writeJSON(w, http.StatusOK, response)
}

func parseDoneFilter(r *http.Request) *bool {
	doneStr := r.URL.Query().Get("done")
	if doneStr == "" {
		return nil // нет фильтра
	}

	done, err := strconv.ParseBool(doneStr)
	if err != nil {
		return nil // невалидное значение - игнорируем фильтр
	}

	return &done
}

func filterTasksByDone(tasks []*Task, doneFilter *bool) []*Task {
	if doneFilter == nil {
		return tasks // нет фильтра - возвращаем все задачи
	}

	filtered := make([]*Task, 0)
	for _, task := range tasks {
		if task.Done == *doneFilter {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func parsePagination(r *http.Request) (page, limit int) {
	// Значения по умолчанию
	page = 1
	limit = 10

	// Парсим page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Парсим limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > 100 {
				l = 100
			}
			limit = l
		}
	}

	return page, limit
}
```

## Примеры запросов и ответов
### Создание задачи
```http://localhost:8080/api/v1/tasks/2```<br>
<img width="697" height="556" alt="image" src="https://github.com/user-attachments/assets/593ccd2b-8e32-46f1-817d-55e89a9d1546" />

### Получение задачи по id
```http://localhost:8080/api/v1/tasks/2```<br>
<img width="697" height="556" alt="image" src="https://github.com/user-attachments/assets/a30e50c7-26d1-47ce-8519-db60bc53c96f" />

### Получение списка всех задач
```http://localhost:8080/api/v1/tasks```<br>
<img width="774" height="736" alt="image" src="https://github.com/user-attachments/assets/3afc8b11-6993-4e8a-9f01-a2f6e75b808b" />

### Редактирования задачи
```http://localhost:8080/api/v1/tasks/2```<br>
<img width="774" height="564" alt="image" src="https://github.com/user-attachments/assets/f769da1d-55fa-449d-bbaa-1904ecfc6c78" />

### Удаление задачи
```http://localhost:8080/api/v1/tasks/2```<br>

### Пагинация и фильтрация
```http://localhost:8080/api/v1/tasks?done=true&page=1&limit=10```<br>
Ответ:<br>
```json
{
    "filters": {
        "done": true
    },
    "pagination": {
        "limit": 10,
        "page": 1,
        "total": 5,
        "totalPages": 1
    },
    "tasks": [
        {
            "id": 2,
            "title": "qweeee",
            "done": true,
            "created_at": "2025-11-24T23:49:13.542153968+03:00",
            "updated_at": "2025-11-25T19:31:04.909742157+03:00"
        },
        {
            "id": 5,
            "title": "qweeee",
            "done": true,
            "created_at": "2025-11-25T19:35:03.760375449+03:00",
            "updated_at": "2025-11-25T19:35:17.35229321+03:00"
        },
        {
            "id": 6,
            "title": "qweeee",
            "done": true,
            "created_at": "2025-11-25T19:35:04.660470614+03:00",
            "updated_at": "2025-11-25T19:35:18.731950952+03:00"
        },
        {
            "id": 3,
            "title": "qweeee",
            "done": true,
            "created_at": "2025-11-25T19:28:36.763186158+03:00",
            "updated_at": "2025-11-25T19:35:14.809690021+03:00"
        },
        {
            "id": 4,
            "title": "qweeee",
            "done": true,
            "created_at": "2025-11-25T19:35:03.072573843+03:00",
            "updated_at": "2025-11-25T19:35:16.073122141+03:00"
        }
    ]
}
```
### Валидация длины title
<img width="862" height="557" alt="image" src="https://github.com/user-attachments/assets/b1873565-03b2-42a9-9a72-4d1aa2d55b35" />

## Обработка ошибок и коды ответа
<ul>
  <li>200 OK - успешный запрос</li>
  <li>201 Created - задача создана</li>
  <li>400 Bad Request - Невалидные данные</li>
  <li>404 Not Found - Задача не найдена</li>
  <li>204 No Content - Успешное удаление</li>
</ul>

<table>
  <tr>
    <td><strong>Маршрут</strong></td>
    <td><strong>Метод</strong></td>
    <td><strong>Запрос</strong></td>
    <td><strong>Ожидаемый ответ</strong></td>
    <td><strong>Фактический ответ</strong></td>
  </tr>
  <tr>
    <td>/api/v1/tasks</td>
    <td>POST</td>
    <td>{"title":"Test"}</td>
    <td>201 Created + задача</td>
    <td>201 Created + задача</td>
  </tr>
  <tr>
    <td>/api/v1/tasks</td>
    <td>POST</td>
    <td>{"title":"A"}</td>
    <td>400 Bad Request</td>
    <td>400 Bad Request</td>
  </tr>
  <tr>
    <td>/api/v1/tasks</td>
    <td>GET</td>
    <td>-</td>
    <td>200 OK + список</td>
    <td>200 OK + список</td>
  </tr>
  <tr>
    <td>/api/v1/tasks?page=1&limit=2</td>
    <td>GET</td>
    <td>-</td>
    <td>200 OK + пагинация</td>
    <td>200 OK + пагинация</td>
  </tr>
  <tr>
    <td>/api/v1/tasks/1</td>
    <td>GET</td>
    <td>-</td>
    <td>200 OK + задача</td>
    <td>200 OK + задача</td>
  </tr>
  <tr>
    <td>/api/v1/tasks/999</td>
    <td>GET</td>
    <td>-</td>
    <td>404 Not Found</td>
    <td>404 Not Found</td>
  </tr>
  <tr>
    <td>/api/v1/tasks/1</td>
    <td>PUT</td>
    <td>{"title":"Updated"}</td>
    <td>200 OK + задача</td>
    <td>200 OK + задача</td>
  </tr>
  <tr>
    <td>/api/v1/tasks/1</td>
    <td>DELETE</td>
    <td>-</td>
    <td>204 No Content</td>
    <td>204 No Content</td>
  </tr>
  <tr>
    <td>/health</td>
    <td>GET</td>
    <td>-</td>
    <td>200 OK + "OK"</td>
    <td>200 OK + "OK"</td>
  </tr>
</table>

## Выводы
### Что получилось:
<ul>
  <li>Успешно реализован полнофункциональный CRUD API</li>
  <li>Работает пагинация</li>
  <li>Добавлена валидация длины заголовка задач</li>
  <li>Middleware логирования и CORS работают корректно</li>
  <li>Реализована версионирование API</li>
</ul>

### Что можно улучшить
<ul>
  <li>Добавить аутентификацию поьзователей</li>
  <li>Интегрировать реальную базу данных вместо in-memory хранилища</li>
</ul>

## Контрольные вопросы
<ol>
  <li>
    Чем отличается роутер стандартной библиотеки (http.ServeMux) от chi/gorilla/mux? Какие преимущества chi?<br>
    <ul><strong>http.ServeMux (стандартный):</strong>
      <li>Не умеет обрабатывать параметры</li>
      <li>Нет встроенной поддержки middleware</li>
      <li>Простой, но ограниченный</li>
    </ul>
    <ul><strong>chi/gorilla/mux:</strong>
      <li>Умеет обрабатывать параметры: /tasks/123 → id=123</li>
      <li>Можно группировать маршруты</li>
      <li>Преимущество chi - он быстрее и проще</li>
    </ul>
  </li>
  <li>
    Чем отличаются GET, POST, PUT, PATCH, DELETE и как их корректно применять в REST?
    <ul><strong>GET</strong> - "посмотреть"
      <li>Получить данные (список задач, одну задачу)</li>
    </ul>
    <ul><strong>POST</strong> - "создать"
      <li>Создать новую задачу</li>
      <li>Отправляем данные в теле запроса</li>
    </ul>
    <ul><strong>PUT</strong> - "полностью заменить"
      <li>Обновить ВСЕ поля задачи  </li>
      <li>Если поле не указали - оно станет пустым</li>
    </ul>
    <ul><strong>PATCH</strong> - "частично обновить"
      <li>Обновить только некоторые поля </li>
      <li>Остальные поля не трогаем</li>
    </ul>
    <ul><strong>DELETE</strong> - "удалить"
      <li>Удалить задачу</li>
    </ul>
  </li>
  <li>
    Что такое middleware и какие типичные задачи им решают? Примеры.<br>
    Middleware - промежуточный слой между клиентом и сервером<br>
    <ul> Что делают:
      <li>Логирование - что и когда было сделано</li>
      <li>CORS - разрешает доступ с других сайтов</li>
      <li>Авторизация - проверяет пароль/токен</li>
    </ul>
    Пример: Middleware проверяет, есть ли у пользователя токен доступа. Если токен есть — пропускает дальше, если нет — возвращает ошибку.
  </li>
  <li>
    Как правильно возвращать ошибки и коды статусов в JSON-API? Приведите примеры.
    <img width="265" height="67" alt="image" src="https://github.com/user-attachments/assets/fa19a9a1-caa4-4e46-acd7-363756e32851" /><br>
    С кодом 404<br>
    <ul><strong>Основные коды</strong>
      <li>200 - ок</li>
      <li>201 - создано (после POST)</li>
      <li>400 - некорректные данные с клиента</li>
      <li>500 - ошибка на сервере</li>
    </ul>
    Простой способ<br>
    <img width="675" height="116" alt="image" src="https://github.com/user-attachments/assets/79cf3dca-bac5-4af3-a47e-9c981833c5ed" /><br>
  </li>
  <li>
    Какие подходы к хранению данных можно использовать в учебном CRUD-проекте и что поменяется в коде при переходе с памяти на БД?<br>
    <strong>В памяти:</strong><br>
    <img width="341" height="133" alt="image" src="https://github.com/user-attachments/assets/771595b5-4a54-4a95-8f23-fc26d620447f" /><br>
    <strong>Плюсы</strong> - быстро и просто<br>
    <strong>Минусы</strong> - данные теряются при перезапуске<br>
    <ul><strong>База данных</strong>
      <li>Repo методы (из repo.go) будут делать SQL запросы вместо работы с map</li>
      <li>Нужно подключение к БД в main.go</li>
    </ul>
    <ul><strong>Переход с памяти за БД</strong>
      <li>Меняется только repo.go</li>
      <li>API остаётся тем же</li>
    </ul>
  </li>
</ol>
