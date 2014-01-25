package api

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password,omitempty" db:"-"`
	Realname  string   `json:"real_name"`
	Followers []string `json:"followers" db:"-"`
	Following []string `json:"following" db:"-"`
}
