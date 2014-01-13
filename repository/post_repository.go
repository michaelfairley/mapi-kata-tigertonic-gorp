package repository

import (
	"database/sql"
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-tigertonic-gorp/utils"
	"strconv"
	"strings"
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
	return repo.FindByUserIds([]uint64{id}, after)
}

func (repo PostRepository) FindByUserIds(ids []uint64, after *uint64) []*api.DbPost {
	var posts []*api.DbPost
	var err error

	// Parameter binding for arrays would be nice.
	strIds := make([]string, len(ids))
	for i := range ids {
		strIds[i] = strconv.FormatUint(ids[i], 10)
	}
	in := strings.Join(strIds, ", ")

	if after == nil {
		_, err = repo.Db.Select(&posts, "SELECT * FROM posts WHERE user_id IN ("+in+") ORDER BY id DESC LIMIT 50")
	} else {
		_, err = repo.Db.Select(&posts, "SELECT * FROM posts WHERE user_id IN ("+in+") AND id < $1 ORDER BY id DESC LIMIT 50", after)

	}
	utils.CheckErr(err)

	return posts
}
