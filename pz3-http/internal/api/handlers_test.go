package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/pz3-http/internal/storage"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	// Создаем временный handler для health
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Вместо сравнения строк, парсим JSON и проверяем структуру
	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("handler returned unexpected status: got %v want %v", response["status"], "ok")
	}
}

func TestCreateTaskHandler(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)

	taskData := map[string]string{"title": "Test task"}
	jsonData, _ := json.Marshal(taskData)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(h.CreateTask)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Проверяем что задача создалась
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response["title"] != "Test task" {
		t.Errorf("handler returned unexpected task title: got %v want %v", response["title"], "Test task")
	}
}

func TestListTasksHandler(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)

	// Сначала создаем задачу
	store.Create("Test task")

	req := httptest.NewRequest("GET", "/tasks", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(h.ListTasks)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var tasks []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &tasks)

	if len(tasks) != 1 {
		t.Errorf("handler returned unexpected number of tasks: got %v want %v", len(tasks), 1)
	}
}
