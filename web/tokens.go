package web

import (
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/repository"
	"net/http"
	"net/url"
)

type TokensResource struct {
	Repo     repository.TokenRepository
	UserRepo repository.UserRepository
}

func (resource *TokensResource) CreateToken(url *url.URL, in_headers http.Header, auth *api.Auth) (int, http.Header, interface{}, error) {
	user := resource.UserRepo.FindByUsername(auth.Username)
	if user == nil {
		return 401, nil, nil, nil
	}

	token := api.NewTokenForUser(*user)
	resource.Repo.Insert(&token)

	return 200, nil, token, nil
}
