package repository

import (
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-tigertonic-gorp/utils"
)

type TokenRepository struct {
	Db *gorp.DbMap
}

func (repo TokenRepository) Insert(token *api.Token) {
	err := repo.Db.Insert(token)
	utils.CheckErr(err)
}
