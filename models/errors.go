package models

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse is used to display errors as API responses
type ErrResponse struct {
	StatusCode int    `json:"status"`            // user-level status message
	ErrorText  string `json:"message,omitempty"` // application-level error message, for debugging
}

// Render sets HTTP Status code from the ErrResponse struct
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}
