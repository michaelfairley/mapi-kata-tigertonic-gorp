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

func (resource *UserResource) CreateUser(url *url.URL, in_headers http.Header, user *api.User) (int, http.Header, *string, error) {
	resource.Db.Insert(user)

	headers := http.Header{
		"Location": []string{"http://localhost:12346/users/" + user.Username},
	}
	return 303, headers, nil, nil
}

func (resource *UserResource) GetUser(url *url.URL, in_headers http.Header, _ interface{}) (int, http.Header, *api.User, error) {
	user := api.User{}

	err := resource.Db.SelectOne(&user, "SELECT * FROM users WHERE username = $1", url.Query().Get("username"))
	if err != nil {
		panic(err)
	}

	return 200, http.Header{}, &user, nil
}
