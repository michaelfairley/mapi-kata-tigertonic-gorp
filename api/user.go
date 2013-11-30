package api

type User struct {
	Id       uint64 `json:"-"`
	Username string `json:"username"`
	Password string `json:"-"`
	Realname string `json:"real_name"`
}
