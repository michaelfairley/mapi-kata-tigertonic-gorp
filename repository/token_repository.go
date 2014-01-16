package repository

import (
	"database/sql"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/utils"
)

type TokenRepository struct {
	Db *gorp.DbMap
}

func (repo TokenRepository) Insert(token *api.Token) {
	err := repo.Db.Insert(token)
	utils.CheckErr(err)
}

func (repo TokenRepository) FindByValue(value string) *api.Token {
	token := &api.Token{}

	err := repo.Db.SelectOne(token, "SELECT * FROM tokens WHERE value = $1", value)
	if err == sql.ErrNoRows {
		return nil
	}
	utils.CheckErr(err)

	return token

}
