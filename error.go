package ysgo

import (
	"errors"
	"fmt"
	"strings"
)

type YSError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *YSError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("YS API error %d: %s - %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("YS API error %d: %s", e.Code, e.Message)
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

func parseYSError(statusCode int, body []byte, fallback string) error {
	message := strings.TrimSpace(string(body))
	if message == "" {
		message = fallback
	}

	if strings.HasPrefix(message, "ERR") {
		parts := strings.SplitN(message, "\n", 2)
		if len(parts) == 2 {
			message = strings.TrimSpace(parts[1])
		}
	}

	message = strings.TrimSuffix(message, ".")
	message = strings.TrimSpace(message)

	switch {
	case strings.Contains(message, "需提供:glmm"):
		return ErrInvalidCredentials
	case strings.Contains(message, "用户信息丢失"):
		return ErrSessionNotInitialized
	case strings.Contains(message, "管理员登陆后才可以编辑目录"):
		return ErrAdminRequired
	case strings.Contains(message, "需要输入登陆密码") || strings.Contains(message, "限制访客登陆"):
		return ErrSpacePasswordRequired
	case strings.Contains(message, "登陆密码不正确"):
		return ErrInvalidSpacePassword
	default:
		return &YSError{Code: statusCode, Message: fallback, Detail: message}
	}
}

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUnauthorized          = errors.New("unauthorized access")
	ErrNotFound              = errors.New("resource not found")
	ErrRateLimited           = errors.New("rate limited")
	ErrServerError           = errors.New("internal server error")
	ErrBadRequest            = errors.New("bad request")
	ErrSessionNotInitialized = errors.New("session not initialized")
	ErrAdminRequired         = errors.New("admin required")
	ErrSpacePasswordRequired = errors.New("space password required")
	ErrInvalidSpacePassword  = errors.New("invalid space password")
)

const httpStatusAlreadyReported = 208
