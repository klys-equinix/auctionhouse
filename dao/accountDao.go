package dao

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang-poc/dto"
	u "golang-poc/utils"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Account struct {
	CommonModelFields
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

func (account *Account) ValidateUnique() (string, bool) {
	temp := &Account{}

	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error

	if account.failedToGetRecord(err) {
		return "Connection error. Please retry", false
	}

	if temp.Email != "" {
		return "Email address already in use by another user.", false
	}

	return "Requirement passed", true
}

func Create(createAccountDto *dto.CreateAccountDto) (*dto.AccountDto, error) {
	account := NewAccount(createAccountDto)

	if resp, ok := account.ValidateUnique(); !ok {
		return nil, errors.New(resp)
	}

	GetDB().Create(account)

	account.generateJwtToken()

	accountDto := &dto.AccountDto{Email: account.Email, Token: account.Token}
	return accountDto, nil
}

func Login(email, password string) map[string]interface{} {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(400, "Email address not found")
		}
		return u.Message(400, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(400, "Invalid login credentials. Please try again")
	}

	account.Password = ""

	account.generateJwtToken()

	resp := u.Message(400, "Logged In")
	resp["account"] = &dto.AccountDto{Email: account.Email, Token: account.Token}
	return resp
}

func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" {
		return nil
	}

	acc.Password = ""
	return acc
}

func (account *Account) generateJwtToken() {
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString
}

func (account *Account) failedToGetRecord(err error) bool {
	return err != nil && err != gorm.ErrRecordNotFound
}

func NewAccount(accountDto *dto.CreateAccountDto) *Account {
	account := &Account{}
	account.Email = accountDto.Email
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(accountDto.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	return account
}
