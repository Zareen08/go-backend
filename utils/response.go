package utils

import (
    "github.com/labstack/echo/v4"
)

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(c echo.Context, status int, message string, data interface{}) error {
    return c.JSON(status, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c echo.Context, status int, message string, errors interface{}) error {
    return c.JSON(status, Response{
        Success: false,
        Message: message,
        Errors:  errors,
    })
}