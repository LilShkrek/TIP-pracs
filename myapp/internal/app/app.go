package app

import (
	"fmt"
	"net/http"

	"github.com/LilShkrek/TIP-pracs/myapp/internal/app/handlers"
	"github.com/LilShkrek/TIP-pracs/myapp/utils"
)

type pingResp struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func Run() {
	mux := http.NewServeMux()

	// Корневой маршрут
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "Hello, Go project structure!")
	})

	// Пример JSON-ручки: /ping
	//mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	//	utils.LogRequest(r)
	//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	_ = json.NewEncoder(w).Encode(pingResp{
	//		Status: "ok",
	//		Time:   time.Now().UTC().Format(time.RFC3339),
	//	})
	//})

	mux.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		utils.WriteErr(w, http.StatusBadRequest, "bad_request_example")
	})

	mux.HandleFunc("/ping", handlers.Ping)

	handler := withRequestID(mux)

	utils.LogInfo("Server is starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		utils.LogError("server error: " + err.Error())
	}
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = utils.NewID16()
		}
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}
