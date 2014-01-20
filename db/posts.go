package db

import (
	"database/sql"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/utils"
	"strconv"
	"strings"
)

type Post struct {
	Id     uint64
	UserId uint64 `db:"user_id"`
	Text   string
}

type PostRepository struct {
	Db *gorp.DbMap
}

func (repo PostRepository) Insert(post *Post) {
	repo.Db.Insert(post)
}

func (repo PostRepository) Find(id uint64) *Post {
	post := &Post{}

	err := repo.Db.SelectOne(post, "SELECT * FROM posts WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil
	}
	utils.CheckErr(err)

	return post
}

func (repo PostRepository) Delete(post *Post) {
	_, err := repo.Db.Delete(post)
	utils.CheckErr(err)
}

func (repo PostRepository) FindByUserId(id uint64, after *uint64) []*Post {
	return repo.FindByUserIds([]uint64{id}, after)
}

func (repo PostRepository) FindByUserIds(ids []uint64, after *uint64) []*Post {
	var posts []*Post
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
