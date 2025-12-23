package core

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/pz10-auth/internal/http/middleware"
	"github.com/go-chi/chi/v5"
)

// Интерфейс для информации о пользователе
type UserInfo interface {
	GetID() int64
	GetEmail() string
	GetRole() string
}

type userRepo interface {
	CheckPassword(email, pass string) (UserInfo, error)
}

type jwtSigner interface {
	Sign(userID int64, email, role string) (string, error)
	SignRefresh(userID int64) (string, error)
	Parse(tokenStr string) (map[string]interface{}, error)
}

type Service struct {
	repo userRepo
	jwt  jwtSigner
}

func NewService(r userRepo, j jwtSigner) *Service {
	return &Service{repo: r, jwt: j}
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string
		Password string
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		httpError(w, 400, "invalid_credentials")
		return
	}

	u, err := s.repo.CheckPassword(in.Email, in.Password)
	if err != nil {
		httpError(w, 401, "unauthorized")
		return
	}

	accessToken, err := s.jwt.Sign(u.GetID(), u.GetEmail(), u.GetRole())
	if err != nil {
		httpError(w, 500, "token_error")
		return
	}

	refreshToken, err := s.jwt.SignRefresh(u.GetID())
	if err != nil {
		httpError(w, 500, "refresh_token_error")
		return
	}

	jsonOK(w, map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *Service) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.RefreshToken == "" {
		httpError(w, 400, "refresh_token_required")
		return
	}

	// Парсим refresh токен
	claims, err := s.jwt.Parse(in.RefreshToken)
	if err != nil {
		httpError(w, 401, "invalid_refresh_token")
		return
	}

	// Проверяем что это refresh токен
	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		httpError(w, 401, "not_a_refresh_token")
		return
	}

	// Получаем ID пользователя
	userID, ok := claims["sub"].(float64)
	if !ok {
		httpError(w, 401, "invalid_token")
		return
	}

	// Получаем данные пользователя (упрощённо)
	var email, role string
	if userID == 1 {
		email, role = "admin@example.com", "admin"
	} else if userID == 2 {
		email, role = "user@example.com", "user"
	} else {
		httpError(w, 401, "user_not_found")
		return
	}

	// Генерируем новую пару токенов
	accessToken, err := s.jwt.Sign(int64(userID), email, role)
	if err != nil {
		httpError(w, 500, "token_generation_error")
		return
	}

	refreshToken, err := s.jwt.SignRefresh(int64(userID))
	if err != nil {
		httpError(w, 500, "refresh_token_generation_error")
		return
	}

	jsonOK(w, map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *Service) MeHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.CtxClaimsKey).(map[string]any)
	jsonOK(w, map[string]any{
		"id":    claims["sub"],
		"email": claims["email"],
		"role":  claims["role"],
	})
}

func (s *Service) AdminStats(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]any{"users": 2, "version": "1.0"})
}

func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.CtxClaimsKey).(map[string]any)
	userRole, _ := claims["role"].(string)
	userID, _ := claims["sub"].(float64)

	// Получаем ID из URL параметра
	requestedID := chi.URLParam(r, "id")

	// ABAC правило: user может получать только свой профиль
	if userRole == "user" {
		// Правильное сравнение ID
		currentUserID := strconv.Itoa(int(userID))
		if requestedID != currentUserID {
			httpError(w, 403, "access_denied")
			return
		}
	}

	// Возвращаем моковые данные
	jsonOK(w, map[string]any{
		"id":    requestedID,
		"email": "user" + requestedID + "@example.com",
		"role":  "user",
	})
}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
