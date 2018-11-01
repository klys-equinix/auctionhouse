package dto

import (
	"strings"
)

type CreateAccountDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (account *CreateAccountDto) Validate() (string, bool) {

	if validateEmail(account.Email) {
		return "Email address is required", false
	}

	if validatePassword(account.Password) {
		return "Password is too short", false
	}

	return "Requirement passed", true
}

func validatePassword(password string) bool {
	return len(password) < 6
}

func validateEmail(email string) bool {
	return !strings.Contains(email, "@")
}
