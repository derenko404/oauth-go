package response

import "net/http"

type Context interface {
	JSON(code int, obj any)
	Status(code int)
}

type APISuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data,omitempty"`
}

type APIError struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"NOT_FOUND"`
	Details string `json:"details,omitempty" example:"The requested resource was not found."`
}

type APIErrorResponse struct {
	Error   *APIError `json:"error"`
	Success bool      `json:"success" example:"false"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewError(code int, message string, details string) *APIError {
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
	c.JSON(http.StatusOK, APISuccessResponse{
		Success: true,
		Data:    data,
	})
}

func RespondCreated(c Context, data any) {
	c.JSON(http.StatusCreated, APISuccessResponse{
		Success: true,
		Data:    data,
	})
}

func RespondNoContent(c Context) {
	c.Status(http.StatusNoContent)
}

func RespondError(c Context, err *APIError) {
	c.JSON(err.Code, APIErrorResponse{
		Success: false,
		Error:   err,
	})
}
