package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func extractToken(r *http.Request) (string, error) {
	// Пробуем получить токен из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Ожидаемый формат: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
		return "", fmt.Errorf("invalid Authorization header format")
	}

	// Если токен не найден в заголовке, пробуем получить его из куки
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", fmt.Errorf("token not found in cookies or headers")
	}

	return cookie.Value, nil
}

func validateToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("Invalid or expired token")
	}
	return token, nil
}

func extractUserID(token *jwt.Token) (int64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("Invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("Invalid token payload")
	}

	return int64(userID), nil
}
