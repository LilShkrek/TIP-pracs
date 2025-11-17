# Практическое занятие №3
## Заикин Д.Ю. ЭФМО-02-25
Реализация простого HTTP-сервера на стандартной библиотеке net/http. Обработка запросов GET/POST<br>
<ul><strong>Цели:</strong></ul>
<li>Освоить базовую работу со стандартной библиотекой net/http без сторонних фреймворков</li>
<li>Научиться поднимать HTTP-сервер, настраивать маршрутизацию через http.ServeMux</li>
<li>Научиться обрабатывать параметры запроса (query, path), тело запроса (JSON/form-data) и формировать корректные ответы (код статуса, заголовки, JSON).</li>
<li>Научиться базовому логированию запросов и обработке ошибок.</li>

### Запуск
Зайти в папку /pz3-http и выполнить ```go run ./cmd/server```

### Структура проекта
<img width="444" height="383" alt="image" src="https://github.com/user-attachments/assets/5943cbe5-7c73-49ae-9f82-1bd56de1bba6" />

### Список эндпоинтов
[requests.md](https://github.com/LilShkrek/TIP-pracs/blob/prac3/requests.md)

### Примеры запросов
### GET /health<br>
<img width="532" height="495" alt="image" src="https://github.com/user-attachments/assets/dcfafeab-194c-4ad0-803f-1c25c119cd42" />

### GET /tasks<br>
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/093a940b-141e-45a0-8eb0-c5f5bd9669b2" />

### POST /tasks<br>
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/7e5152cd-e8c8-40a1-91c6-105a84696c61" />

### GET /tasks/{id}<br>
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/c9efb69a-0946-44fc-9ab3-c78aa18b8130" />

### PATCH /tasks/{id}<br>
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/4d9566a3-eced-48cc-9f7c-8cd9b009cfeb" />

### DELETE /tasks/{id}<br>
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/f0cd297f-f291-4771-bca5-1608281449ab" />

### Валидация title (пример с длиной меньше 3)<br>
<img width="516" height="494" alt="image" src="https://github.com/user-attachments/assets/fa48a679-209c-4436-ba86-348bf445360b" />

### Валидация title (пример с длиной больше 140)<br>
<img width="1423" height="522" alt="image" src="https://github.com/user-attachments/assets/44ec85d4-d55f-4b54-9061-6e58f96fc4b7" />

### Конфигурация порта в переменной окружения PORT<br>
<img width="828" height="185" alt="image" src="https://github.com/user-attachments/assets/bf19ec15-df98-48e0-9442-0351ea278fed" />

### Graceful shutdown<br>
<img width="751" height="223" alt="image" src="https://github.com/user-attachments/assets/96e217e8-fcd7-4a18-903d-274f5c1437d6" />

### Юнит-тесты<br>
<img width="780" height="44" alt="image" src="https://github.com/user-attachments/assets/66a6e63c-a405-4cb8-a828-021aa1e6d4eb" />


## Использование Makefile
```make help``` - вывод всех доступных команд<br>
```make run``` - запуск сервера<br>
```make build``` - сборка бинарника<br>
```make test``` - запуск тестов<br>
