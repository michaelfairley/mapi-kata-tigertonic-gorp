package repository

import (
	"database/sql"
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-tigertonic-gorp/utils"
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
	err := repo.Db.Insert(user)
	utils.CheckErr(err)
}

func (repo UserRepository) FindByUsername(username string) *api.User {
	user := &api.User{}

	err := repo.Db.SelectOne(user, "SELECT * FROM users WHERE username = $1", username)
	if err == sql.ErrNoRows {
		return nil
	}
	utils.CheckErr(err)

	return user
}

func (repo UserRepository) Find(id uint64) *api.User {
	user := &api.User{}

	err := repo.Db.SelectOne(user, "SELECT * FROM users WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil
	}
	utils.CheckErr(err)

	return user
}

func (repo UserRepository) FindFollowers(user *api.User) []*api.User {
	var users []*api.User

	_, err := repo.Db.Select(&users, "SELECT id, username FROM users JOIN followings ON followings.follower_id = users.id WHERE followings.followee_id = $1", user.Id)
	utils.CheckErr(err)

	return users
}

func (repo UserRepository) FindFollowing(user *api.User) []*api.User {
	var users []*api.User

	_, err := repo.Db.Select(&users, "SELECT id, username FROM users JOIN followings ON followings.followee_id = users.id WHERE followings.follower_id = $1", user.Id)
	utils.CheckErr(err)

	return users
}

func (repo UserRepository) Follow(follower, followee *api.User) {
	repo.Db.Exec("INSERT INTO followings (follower_id, followee_id) VALUES ($1, $2)", follower.Id, followee.Id)
}

func (repo UserRepository) Unfollow(follower, followee *api.User) bool {
	res, err := repo.Db.Exec("DELETE FROM followings WHERE follower_id = $1 AND followee_id = $2", follower.Id, followee.Id)
	utils.CheckErr(err)
	rows, err := res.RowsAffected()
	utils.CheckErr(err)
	return rows > 0
}
