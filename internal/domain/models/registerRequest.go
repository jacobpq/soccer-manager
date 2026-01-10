package models

import (
	"errors"
	"net/mail"
	"strings"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TeamName string `json:"team_name"`
	Country  string `json:"country"`
}

func (r *RegisterRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email_required")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid_email")
	}

	if len(r.Password) < 6 {
		return errors.New("password_short")
	}

	if strings.TrimSpace(r.TeamName) == "" {
		return errors.New("team_name_required")
	}

	if strings.TrimSpace(r.Country) == "" {
		return errors.New("country_required")
	}

	return nil
}
