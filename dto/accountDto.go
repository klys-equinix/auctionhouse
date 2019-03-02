package dto

type AccountDto struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
