package web

import (
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/repository"
	"net/http"
	"net/url"
)

type TokensResource struct {
	Repo     repository.TokenRepository
	UserRepo repository.UserRepository
}

func (resource *TokensResource) CreateToken(url *url.URL, inHeaders http.Header, auth *api.Auth) (int, http.Header, interface{}, error) {
	user := resource.UserRepo.FindByUsername(auth.Username)
	if user == nil {
		return 401, nil, nil, nil
	}
	if !user.CheckPassword(auth.Password) {
		return 401, nil, nil, nil
	}

	token := api.NewTokenForUser(*user)
	resource.Repo.Insert(&token)

	return 200, nil, token, nil
}
