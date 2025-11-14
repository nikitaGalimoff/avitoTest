package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware представляет простую аутентификацию через заголовки
// В реальном приложении здесь была бы полноценная JWT аутентификация
type AuthMiddleware struct {
	adminToken string
	userToken  string
}

// NewAuthMiddleware создает новый middleware для аутентификации
func NewAuthMiddleware(adminToken, userToken string) *AuthMiddleware {
	return &AuthMiddleware{
		adminToken: adminToken,
		userToken:  userToken,
	}
}

// AdminOnly проверяет, что запрос содержит админский токен
func (a *AuthMiddleware) AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return WriteError(c, nil, 401)
		}

		// Убираем префикс "Bearer " если есть
		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimSpace(token)

		if token != a.adminToken {
			return WriteError(c, nil, 401)
		}

		return next(c)
	}
}

// AdminOrUser проверяет, что запрос содержит админский или пользовательский токен
func (a *AuthMiddleware) AdminOrUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return WriteError(c, nil, 401)
		}

		// Убираем префикс "Bearer " если есть
		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimSpace(token)

		if token != a.adminToken && token != a.userToken {
			return WriteError(c, nil, 401)
		}

		return next(c)
	}
}
