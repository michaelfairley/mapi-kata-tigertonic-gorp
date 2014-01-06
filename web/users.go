package web

import (
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/repository"
	"net/http"
	"net/url"
)

type UserResource struct {
	Repository repository.UserRepository
}

func (resource *UserResource) CreateUser(url *url.URL, inHeaders http.Header, user *api.User) (int, http.Header, interface{}, error) {
	existing := resource.Repository.ContainsUserWithUsername(user.Username)

	if existing {
		return validationError(map[string][]string{"username": []string{"is taken"}})
	}
	if len(user.Password) < 8 {
		return validationError(map[string][]string{"password": []string{"is too short"}})
	}

	resource.Repository.Insert(user)

	headers := http.Header{
		"Location": []string{"http://localhost:12346/users/" + user.Username},
	}
	return 303, headers, nil, nil
}

func (resource *UserResource) GetUser(url *url.URL, inHeaders http.Header, _ interface{}) (int, http.Header, *api.User, error) {
	user := resource.Repository.FindByUsername(url.Query().Get("username"))

	return 200, http.Header{}, user, nil
}
