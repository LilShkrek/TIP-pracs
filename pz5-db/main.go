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
