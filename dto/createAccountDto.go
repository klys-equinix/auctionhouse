package dto

import (
	u "golang-poc/utils"
	"strings"
)

type CreateAccountDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (account *CreateAccountDto) Validate() (map[string]interface{}, bool) {

	if validateEmail(account.Email) {
		return u.Message(400, "Email address is required"), false
	}

	if validatePassword(account.Password) {
		return u.Message(400, "Password is too short"), false
	}

	return u.Message(200, "Requirement passed"), true
}

func validatePassword(password string) bool {
	return len(password) < 6
}

func validateEmail(email string) bool {
	return !strings.Contains(email, "@")
}
