package response

import "net/http"

type Context interface {
	JSON(code int, obj any)
	Status(code int)
}

type APIResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewError(code int, message string, details any) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

var (
	ErrorNotFound          = NewError(http.StatusNotFound, "NOT_FOUND", "The requested resource was not found.")
	ErrInvalidInput        = NewError(http.StatusBadRequest, "INVALID_INPUT", "The input is invalid.")
	ErrUnauthorized        = NewError(http.StatusUnauthorized, "UNAUTHORIZED", "You are not authorized to access this resource.")
	ErrOAuth               = NewError(http.StatusBadRequest, "OAUTH_ERROR", "OAuth error.")
	ErrInternalServerError = NewError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error.")
)

func RespondSuccess(c Context, data any) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

func RespondCreated(c Context, data any) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

func RespondNoContent(c Context) {
	c.Status(http.StatusNoContent)
}

func RespondError(c Context, err *APIError) {
	c.JSON(err.Code, APIResponse{
		Success: false,
		Error:   err,
	})
}
