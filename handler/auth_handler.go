package handler

import (
    "spotsync-api/dto"
    "spotsync-api/service"
    "spotsync-api/utils"
    "github.com/labstack/echo/v4"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
    var req dto.RegisterRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body", err.Error())
    }

    if err := c.Validate(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed", err.Error())
    }

    user, err := h.authService.Register(&req)
    if err != nil {
        if err == utils.ErrEmailAlreadyExists {
            return utils.ErrorResponse(c, 400, "Email already registered", nil)
        }
        return utils.ErrorResponse(c, 500, "Registration failed", err.Error())
    }

    return utils.SuccessResponse(c, 201, "User registered successfully", user)
}

func (h *AuthHandler) Login(c echo.Context) error {
    var req dto.LoginRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body", err.Error())
    }

    if err := c.Validate(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed", err.Error())
    }

    response, err := h.authService.Login(&req)
    if err != nil {
        if err == utils.ErrInvalidCredentials {
            return utils.ErrorResponse(c, 401, "Invalid credentials", nil)
        }
        return utils.ErrorResponse(c, 500, "Login failed", err.Error())
    }

    return utils.SuccessResponse(c, 200, "Login successful", response)
}