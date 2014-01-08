package repository

import (
	"database/sql"
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-tigertonic-gorp/utils"
)

type PostRepository struct {
	Db *gorp.DbMap
}

func (repo PostRepository) Insert(post *api.DbPost) {
	repo.Db.Insert(post)
}

func (repo PostRepository) Find(id uint64) *api.DbPost {
	post := &api.DbPost{}

	err := repo.Db.SelectOne(post, "SELECT * FROM posts WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil
	}
	utils.CheckErr(err)

	return post
}

func (repo PostRepository) Delete(post *api.DbPost) {
	_, err := repo.Db.Delete(post)
	utils.CheckErr(err)
}

func (repo PostRepository) FindByUserId(id uint64, after *uint64) []*api.DbPost {
	var posts []*api.DbPost
	var err error

	if after == nil {
		_, err = repo.Db.Select(&posts, "SELECT * FROM posts WHERE user_id = $1 ORDER BY id DESC LIMIT 50", id)
	} else {
		_, err = repo.Db.Select(&posts, "SELECT * FROM posts WHERE user_id = $1 AND id < $2 ORDER BY id DESC LIMIT 50", id, after)

	}
	utils.CheckErr(err)

	return posts
}
