package middleware

import (
	"net/http"
	"strings"
)

var authKeyToken string

// InitMiddleware инициализирует middleware с заданным секретным ключом
// Функция должна быть вызвана при старте приложения для установки ключа
func InitMiddleware(key string) {
	authKeyToken = key
}

// AuthMiddleware проверяет аутентификацию по Bearer токену
//
// Требует наличия заголовка Authorization в формате:
//
//	Authorization: Bearer <token>
//
// При отсутствии или неверном токене возвращает:
//
//	401 Unauthorized или 400 Bad Request
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Authentication token required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authentication header format is incorrect", http.StatusBadRequest)
			return
		}

		if parts[1] != authKeyToken {
			http.Error(w, "Authentication token is invalid", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
