package db

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/utils"
)

type TokenRepository struct {
	Db *gorp.DbMap
}

type Token struct {
	Value  string
	UserId uint64 `db:"user_id"`
}

func (repo TokenRepository) Insert(token *Token) {
	err := repo.Db.Insert(token)
	utils.CheckErr(err)
}

func (repo TokenRepository) FindByValue(value string) *Token {
	token := &Token{}

	err := repo.Db.SelectOne(token, "SELECT * FROM tokens WHERE value = $1", value)
	if err == sql.ErrNoRows {
		return nil
	}
	utils.CheckErr(err)

	return token

}

func GenerateTokenForUser(user User) Token {
	return Token{generateToken(), user.Id}
}

func generateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	b[8] = (b[8] | 0x80) & 0xBF
	b[6] = (b[6] | 0x40) & 0x4F

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
