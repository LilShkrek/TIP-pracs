package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)          // GET /tasks
	r.Post("/", h.create)       // POST /tasks
	r.Get("/{id}", h.get)       // GET /tasks/{id}
	r.Put("/{id}", h.update)    // PUT /tasks/{id}
	r.Delete("/{id}", h.delete) // DELETE /tasks/{id}
	return r
}

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

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	t, err := h.repo.Get(id)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

type createReq struct {
	Title string `json:"title"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		httpError(w, http.StatusBadRequest, "title is required")
		return
	} else if len(req.Title) > 100 { // Проверка длины заголовка
		httpError(w, http.StatusBadRequest, "title is too long")
		return
	} else if len(req.Title) < 3 {
		httpError(w, http.StatusBadRequest, "title is too short")
		return
	}

	t := h.repo.Create(req.Title)
	writeJSON(w, http.StatusCreated, t)
}

type updateReq struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return
	}
	t, err := h.repo.Update(id, req.Title, req.Done)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	if err := h.repo.Delete(id); err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// helpers

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		httpError(w, http.StatusBadRequest, "invalid id")
		return 0, true
	}
	return id, false
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
