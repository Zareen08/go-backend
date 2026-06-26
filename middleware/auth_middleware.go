package middleware

import (
    "spotsync-api/service"
    "spotsync-api/utils"
    "strings"
    "github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
    authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
    return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" {
            return utils.ErrorResponse(c, 401, "Authorization header required", nil)
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            return utils.ErrorResponse(c, 401, "Invalid authorization format", nil)
        }

        tokenString := parts[1]
        userID, role, err := m.authService.ValidateToken(tokenString)
        if err != nil {
            return utils.ErrorResponse(c, 401, "Invalid or expired token", nil)
        }

        // Set user info in context
        c.Set("user_id", userID)
        c.Set("role", role)

        return next(c)
    }
}

func (m *AuthMiddleware) RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        role, ok := c.Get("role").(string)
        if !ok || role != "admin" {
            return utils.ErrorResponse(c, 403, "Admin access required", nil)
        }
        return next(c)
    }
}