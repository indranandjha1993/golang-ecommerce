package utils

import (
	"github.com/labstack/echo"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func WriteSuccessResponse(c echo.Context, message string, data interface{}, statusCode int) error {
	response := SuccessResponse{Message: message, Data: data}
	return c.JSON(statusCode, response)
}

func WriteErrorResponse(c echo.Context, message string, err error, statusCode int) error {
	response := ErrorResponse{Message: message, Error: err}
	return c.JSON(statusCode, response)
}

func WriteErrorsResponse(c echo.Context, message string, err ErrorWarper, statusCode int) error {
	type response struct {
		Message string            `json:"message"`
		Errors  map[string]string `json:"errors"`
	}
	return c.JSON(statusCode, response{Message: message, Errors: err.Errors})
}
