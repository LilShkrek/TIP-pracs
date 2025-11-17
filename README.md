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
```GET    http://localhost:8080/health```<br>
```GET    http://localhost:8080/tasks```<br>
```POST   http://localhost:8080/tasks```<br>
```GET    http://localhost:8080/tasks/{id}```<br>
```PATCH  http://localhost:8080/tasks/{id}```<br>
```DELETE http://localhost:8080/tasks/{id}```<br>

### Примеры запросов
#GET /health
<img width="532" height="495" alt="image" src="https://github.com/user-attachments/assets/dcfafeab-194c-4ad0-803f-1c25c119cd42" /> <br>

#GET /tasks
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/093a940b-141e-45a0-8eb0-c5f5bd9669b2" /> <br>

#POST /tasks
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/7e5152cd-e8c8-40a1-91c6-105a84696c61" /> <br>

#GET /tasks/{id}
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/c9efb69a-0946-44fc-9ab3-c78aa18b8130" /> <br>

#PATCH /tasks/{id}
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/4d9566a3-eced-48cc-9f7c-8cd9b009cfeb" /><br>

#DELETE /tasks/{id}
<img width="545" height="562" alt="image" src="https://github.com/user-attachments/assets/f0cd297f-f291-4771-bca5-1608281449ab" /><br>



## Артефакты и их расположение
<ol>
  <li>
    <strong>Артефакт:</strong> Dockerfile <br>
    <strong>Расположение:</strong> Папка build/, которая традиционно предназначена для скриптов и конфигураций, связанных со сборкой проекта и созданием артефактов развертывания.
  </li>
  <li>
    <strong>Артефакт:</strong> docker-compose.yml <br>
    <strong>Расположение:</strong> Папка deployments/ используется для описания конфигураций оркестраторов и развертывания (docker-compose, Kubernetes), что логически отделяет процесс сборки от процесса запуска в среде.
  </li>
  <li>
    <strong>Артефакт:</strong> deploy.sh <br>
    <strong>Расположение:</strong> Папка scripts/ предназначена для хранения вспомогательных скриптов на bash, Python или других языках, которые автоматизируют рутинные задачи: сборку, тестирование, развертывание, анализ кода. Это позволяет держать корень проекта чистым, а все скрипты — организованными в одном месте.
  </li>
  <li>
    <strong>Артефакт:</strong> config.yaml <br>
    <strong>Расположение:</strong> Папка configs/ существует именно для хранения шаблонов конфигурационных файлов и примеров настройки. Это позволяет пользователю или оператору быстро найти стандартные настройки, скопировать их и адаптировать под свое окружение, не копаясь в коде хендлеров или сервисов.
  </li>
  <li>
    <strong>Артефакт:</strong> index.html <br>
    <strong>Расположение:</strong> Папка web/static/ — стандартное место для хранения всех статических ресурсов (HTML, CSS, JS, изображений, шрифтов), которые должны быть доступны напрямую по URL.
  </li>
</ol>
