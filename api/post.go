package api

import (
	"fmt"
)

type DbPost struct {
	Id     uint64
	UserId uint64 `db:"user_id"`
	Text   string
}

type Post struct {
	Id     uint64 `json:"id"`
	Text   string `json:"text"`
	Author string `json:"author"`
}

type Posts struct {
	Posts []*Post `json:"posts"`
	Next  Next    `json:"next"`
}

type Next struct {
	Username string
	Type     string
	After    uint64
}

func (next Next) MarshalJSON() ([]byte, error) {
	url := fmt.Sprintf("\"/users/%s/%s?after=%d\"", next.Username, next.Type, next.After)
	return []byte(url), nil
}
