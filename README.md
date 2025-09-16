# Практическое занятие №1
## Заикин Д.Ю. ЭФМО-02-25
Установка и настройка окружения Go.<br>
<strong>Цель:</strong> Развернуть рабочее окружение Go на Windows, создать минимальный HTTP-сервис на net/http, подключить и использовать внешнюю зависимость, собрать и проверить приложение.
<ul>Задание:</ul>
<li>Установить Go и Git, проверить версии.</li>
<li>Инициализировать модуль Go в новом проекте.</li>
<li>Реализовать HTTP-сервер с маршрутами /hello (текст) и /user (JSON).</li>
<li>Подключить внешнюю библиотеку (генерация UUID) и использовать её в /user.</li>
<li>Запустить и проверить ответы curl/браузером.</li>
<li>Собрать бинарник .exe и подготовить README и отчёт.</li>

### Запуск
Зайти в папку /helloapi и выполнить ```go run ./cmd/server```

### Примеры запросов
<img width="1226" height="248" alt="image" src="https://github.com/user-attachments/assets/af7c9be6-ad7a-40d3-817d-053b5ecca26a" />
<img width="1226" height="248" alt="image" src="https://github.com/user-attachments/assets/9666096a-e07b-453e-9812-f30229857f09" />
<img width="1226" height="248" alt="image" src="https://github.com/user-attachments/assets/b2ff188c-f0bd-4adf-83e9-41f1f3d925ff" />

### Запуск с другим портом (дефолтный порт - 8080)
<img width="1080" height="56" alt="image" src="https://github.com/user-attachments/assets/bddd1237-49bf-440e-a008-4f1c11772798" />

### Список эндпоинтов
```http://localhost:8080/hello```<br>
```http://localhost:8080/user```<br>
```http://localhost:8080/health```<br>
### Структура проекта
<img width="162" height="132" alt="image" src="https://github.com/user-attachments/assets/450702c9-b63c-41aa-b31a-64ff9ce3a3db" />

### Версия GoLang
<img width="801" height="49" alt="image" src="https://github.com/user-attachments/assets/41ca9628-cefc-4041-b596-b1d96745aa99" />

### Вывод команд go fmt./... и go vet./...
<img width="812" height="73" alt="image" src="https://github.com/user-attachments/assets/907ddc79-0e38-42b2-b7b9-5bb62f4cd03d" />
