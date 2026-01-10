package models

import (
	"errors"
	"net/mail"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email_required")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid_email")
	}

	if len(r.Password) < 6 {
		return errors.New("password_short")
	}

	return nil
}
