package account

type AccountDto struct {
	Email string `json:"email"`
	Token string `json:"token";sql:"-"`
}
