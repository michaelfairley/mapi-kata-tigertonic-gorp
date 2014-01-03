package web

import (
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/github.com/coopernurse/gorp"
	"net/http"
	"net/url"
)

type UserResource struct {
	Db *gorp.DbMap
}

func (resource *UserResource) CreateUser(url *url.URL, in_headers http.Header, user *api.User) (int, http.Header, interface{}, error) {
	existing, err := resource.Db.SelectInt("SELECT count(*) FROM users WHERE username = $1", user.Username)
	if err != nil {
		panic(err)
	}
	if existing > 0 {
		return validationError(map[string][]string{"username": []string{"is taken"}})
	}
	if len(user.Password) < 8 {
		return validationError(map[string][]string{"password": []string{"is too short"}})
	}

	resource.Db.Insert(user)

	headers := http.Header{
		"Location": []string{"http://localhost:12346/users/" + user.Username},
	}
	return 303, headers, nil, nil
}

func (resource *UserResource) GetUser(url *url.URL, in_headers http.Header, _ interface{}) (int, http.Header, *api.User, error) {
	user := api.User{}

	err := resource.Db.SelectOne(&user, "SELECT username, realname FROM users WHERE username = $1", url.Query().Get("username"))
	if err != nil {
		panic(err)
	}

	return 200, http.Header{}, &user, nil
}
