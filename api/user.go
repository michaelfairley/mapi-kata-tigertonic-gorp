package api

import (
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/code.google.com/p/go.crypto/bcrypt"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/coopernurse/gorp"
)

type User struct {
	Id              uint64   `json:"-"`
	Username        string   `json:"username"`
	Password        string   `json:"password,omitempty" db:"-"`
	Realname        string   `json:"real_name"`
	CryptedPassword []byte   `json:"-" db:"password"`
	Followers       []string `json:"followers" db:"-"`
	Following       []string `json:"following" db:"-"`
}

func (u *User) PreInsert(s gorp.SqlExecutor) error {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.CryptedPassword = cryptedPassword

	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.CryptedPassword, []byte(password))
	return err == nil
}
