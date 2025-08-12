package ysgo

import "fmt"

type YSError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *YSError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("YS API Error %d: %s - %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("YS API Error %d: %s", e.Code, e.Message)
}

func NewYSError(code int, message string, detail ...string) *YSError {
	err := &YSError{
		Code:    code,
		Message: message,
	}
	if len(detail) > 0 {
		err.Detail = detail[0]
	}
	return err
}

var (
	ErrInvalidCredentials = NewYSError(401, "Invalid credentials")
	ErrUnauthorized       = NewYSError(403, "Unauthorized access")
	ErrNotFound           = NewYSError(404, "Resource not found")
	ErrRateLimited        = NewYSError(429, "Rate limited")
	ErrServerError        = NewYSError(500, "Internal server error")
	ErrBadRequest         = NewYSError(400, "Bad request")
)
