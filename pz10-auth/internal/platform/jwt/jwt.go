package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Validator interface {
	Sign(userID int64, email, role string) (string, error)
	Parse(tokenStr string) (map[string]interface{}, error)
	SignRefresh(userID int64) (string, error) // Добавляем метод для refresh токена
}

type HS256 struct {
	secret []byte
	ttl    time.Duration
}

func NewHS256(secret []byte, ttl time.Duration) *HS256 {
	return &HS256{secret: secret, ttl: ttl}
}

func (h *HS256) Sign(userID int64, email, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"iat":   now.Unix(),
		"exp":   now.Add(h.ttl).Unix(),
		"iss":   "pz10-auth",
		"aud":   "pz10-clients",
		"type":  "access", // Добавляем тип токена
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(h.secret)
}
func (h *HS256) SignRefresh(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  now.Unix(),
		"exp":  now.Add(7 * 24 * time.Hour).Unix(), // 7 дней
		"iss":  "pz10-auth",
		"aud":  "pz10-clients",
		"type": "refresh", // Тип refresh
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(h.secret)
}

func (h *HS256) Parse(tokenStr string) (map[string]interface{}, error) { // Измени тип возврата
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return h.secret, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithAudience("pz10-clients"),
		jwt.WithIssuer("pz10-auth"),
	)
	if err != nil || !t.Valid {
		return nil, err
	}

	// Конвертируем jwt.MapClaims в map[string]interface{}
	claims := t.Claims.(jwt.MapClaims)
	result := make(map[string]interface{})
	for k, v := range claims {
		result[k] = v
	}
	return result, nil
}
