package handlers

import (
	"github.com/go-chi/render"
	"github.com/nimbo-stratuz/bikeshare-directions/models"
)

// ErrBadRequest creates an ErrResponse for 400 Bad Request
func ErrBadRequest(message string) render.Renderer {
	return Err(400, message)
}

// ErrServerError creates an ErrResponse for Server Errors
func ErrServerError() render.Renderer {
	return Err(500, "Internal Server Error")
}

// Err creates an ErrResponse
func Err(status int, message string) render.Renderer {
	return &models.ErrResponse{
		StatusCode: 500,
		ErrorText:  message,
	}
}
