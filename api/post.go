package api

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
