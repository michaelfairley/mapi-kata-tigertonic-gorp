package repository

import (
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/github.com/coopernurse/gorp"
)

type UserRepository struct {
	Db *gorp.DbMap
}

func (repo UserRepository) ContainsUserWithUsername(username string) bool {
	existing, err := repo.Db.SelectInt("SELECT count(*) FROM users WHERE username = $1", username)
	if err != nil {
		panic(err)
	}

	return existing > 0
}

func (repo UserRepository) Insert(user *api.User) {
	repo.Db.Insert(user)
}

func (repo UserRepository) FindByUsername(username string) *api.User {
	user := &api.User{}

	err := repo.Db.SelectOne(user, "SELECT username, realname FROM users WHERE username = $1", username)
	if err != nil {
		panic(err)
	}

	return user
}
