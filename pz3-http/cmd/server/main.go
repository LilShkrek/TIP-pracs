package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/pz3-http/internal/api"
	"example.com/pz3-http/internal/storage"
)

func main() {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		api.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Коллекция
	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	// Элемент
	mux.HandleFunc("GET /tasks/", h.GetTask)
	mux.HandleFunc("PATCH /tasks/", h.UpdateTask)
	mux.HandleFunc("DELETE /tasks/", h.DeleteTask)

	// Подключаем логирование
	loggedHandler := api.Logging(mux)
	CorsHandler := api.EnableCORS(loggedHandler)

	// Создаем HTTP-сервер с настройками
	server := &http.Server{
		Addr:         ":8080",
		Handler:      CorsHandler,
		ReadTimeout:  15 * time.Second, // Максимальное время чтения запроса
		WriteTimeout: 15 * time.Second, // Максимальное время записи ответа
		IdleTimeout:  60 * time.Second, // Максимальное время простоя соединения
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Server starting on %s", server.Addr)
		log.Printf("Available endpoints:")
		log.Printf("  GET    http://localhost%s/health", server.Addr)
		log.Printf("  GET    http://localhost%s/tasks", server.Addr)
		log.Printf("  POST   http://localhost%s/tasks", server.Addr)
		log.Printf("  GET    http://localhost%s/tasks/{id}", server.Addr)
		log.Printf("  PATCH  http://localhost%s/tasks/{id}", server.Addr)
		log.Printf("  DELETE http://localhost%s/tasks/{id}", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Канал для получения сигналов ОС
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Ждем сигнал остановки (Ctrl+C)
	<-stop
	log.Println("Shutdown signal received")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Останавливаем сервер
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
