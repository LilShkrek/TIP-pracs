# Практическое задание №5
## Заикин Д.Ю. ЭФМО-02-25
Подключение к PostgreSQL через database/sql. Выполнение простых запросов (INSERT, SELECT)
<ul><strong>Цели:</strong></ul>
<li>Установить и настроить PostgreSQL локально</li>
<li>Подключиться к БД из Go с помощью database/sql и драйвера PostgreSQL</li>
<li>Выполнить параметризованные запросы INSERT и SELECT</li>
<li>Корректно работать с context, пулом соединений и обработкой ошибок</li>

### Описание окружения
Версия go lang: go version go1.24.7 linux/amd64
Версия PostgreSQ: psql (PostgreSQL) 16.11 (Ubuntu 16.11-0ubuntu0.24.04.1)
Версия OC: Ubuntu 24.04.3 LTS

### Запуск
Зайти в папку /pz5-db и выполнить ```go run .```

### Структура проекта
<img width="449" height="248" alt="image" src="https://github.com/user-attachments/assets/733f8e95-999d-4d9d-aebc-0dc51cfa64f3" />

### Скриншоты работы
Создание БД<br>
<img width="326" height="32" alt="image" src="https://github.com/user-attachments/assets/a2fdbedb-7ce7-441b-bd45-0e64bcd1ef29" />

Список таблиц БД<br>
<img width="363" height="141" alt="image" src="https://github.com/user-attachments/assets/05a30c15-e1e9-4ee8-98f6-6cdc8e384c50" />

Вывод команды ```go run .```<br>
<img width="695" height="679" alt="image" src="https://github.com/user-attachments/assets/23d08a2a-4d1b-40aa-8ae5-d0df0761173c" />

Вывод команды ```SELECT * FROM tasks;```<br>
<img width="695" height="679" alt="image" src="https://github.com/user-attachments/assets/7e035aea-084e-4371-b9b2-5330c8d065c4" />

Код db.go:
```go
package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// настройки пула — достаточно для локалки
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// проверка соединения с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	log.Println("Connected to PostgreSQL")
	return db, nil
}
```
<br>

Код repository.go:
```go
package main

import (
	"context"
	"database/sql"
	"time"
)

// Task — модель для сканирования результатов SELECT
type Task struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo { return &Repo{DB: db} }

// CreateTask — параметризованный INSERT с возвратом id
func (r *Repo) CreateTask(ctx context.Context, title string) (int, error) {
	var id int
	const q = `INSERT INTO tasks (title) VALUES ($1) RETURNING id;`
	err := r.DB.QueryRowContext(ctx, q, title).Scan(&id)
	return id, err
}

// ListTasks — базовый SELECT всех задач (демо для занятия)
func (r *Repo) ListTasks(ctx context.Context) ([]Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks ORDER BY id;`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// ListDone --- возвращает только выполненные (done=true) или невыполненные (done=false) задачи
func (r *Repo) ListDone(ctx context.Context, done bool) ([]Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks WHERE done = $1 ORDER BY id;`

	rows, err := r.DB.QueryContext(ctx, q, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// FindByID --- находит задачу по ID, возвращает nil если не найдена
func (r *Repo) FindByID(ctx context.Context, id int) (*Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks WHERE id = $1;`

	var t Task
	err := r.DB.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// CreateMany --- массовая вставка через транзакцию (все или ничего)
func (r *Repo) CreateMany(ctx context.Context, titles []string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const q = `INSERT INTO tasks (title) VALUES ($1);`

	for _, title := range titles {
		if _, err := tx.ExecContext(ctx, q, title); err != nil {
			return err
		}
	}

	return tx.Commit()
}
```
<br>

Код main.go:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// .env не обязателен; если файла нет — ошибка игнорируется
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// fallback — прямой DSN в коде (только для учебного стенда!)
		dsn = "postgres://postgres:password@localhost:5432/todo?sslmode=disable"
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// 1) Вставим пару задач
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	// Массовая вставка
	err = repo.CreateMany(ctx, []string{"Задача A", "Задача B", "Задача C"})
	if err != nil {
		log.Printf("CreateMany error: %v", err)
	} else {
		fmt.Println("Массовую вставку выполнили!")
	}

	// 2) Прочитаем список задач
	ctxList, cancelList := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err := repo.ListTasks(ctxList)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	// Вывод выполненных задач
	doneTasks, err := repo.ListDone(ctx, true)
	if err != nil {
		log.Printf("ListDone error: %v", err)
	} else {
		fmt.Println("=== Выполненные задачи ===")
		for _, t := range doneTasks {
			fmt.Printf("#%d: %s\n", t.ID, t.Title)
		}
	}

	// Поиск по id
	task, err := repo.FindByID(ctx, 1)
	if err != nil {
		log.Printf("FindByID error: %v", err)
	} else if task != nil {
		fmt.Printf("Найдена задача #%d: %s (done=%v)\n", task.ID, task.Title, task.Done)
	} else {
		fmt.Println("Задача не найдена")
	}

}
```
<br>
