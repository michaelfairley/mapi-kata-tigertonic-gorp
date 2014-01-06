package web

import (
	"errors"
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/github.com/rcrowley/go-tigertonic"
	"github.com/michaelfairley/mapi-tigertonic-gorp/repository"
	"net/http"
	"regexp"
)

var (
	unauthorized = tigertonic.Unauthorized{errors.New("")}
	forbidden    = tigertonic.Forbidden{errors.New("")}
	tokenRegex   = regexp.MustCompile("^Token (.*)$")
)

type Auther struct {
	UserRepo  repository.UserRepository
	TokenRepo repository.TokenRepository
}

func (auther Auther) Auth(headers http.Header) *api.User {
	authentication := headers["Authentication"]
	if len(authentication) != 1 {
		return nil
	}
	if !tokenRegex.MatchString(authentication[0]) {
		return nil
	}
	tokenValue := authentication[0][6:]
	token := auther.TokenRepo.FindByValue(tokenValue)
	if token == nil {
		return nil
	}

	user := auther.UserRepo.Find(token.UserId)

	return user
}

func (auther Auther) Filter(request *http.Request) (http.Header, error) {
	tokenUser := auther.Auth(request.Header)
	if tokenUser == nil {
		return nil, unauthorized
	}

	pathUser := auther.UserRepo.FindByUsername(request.URL.Query().Get("username"))

	if pathUser == nil {
		return nil, forbidden
	}

	if pathUser.Id != tokenUser.Id {
		return nil, forbidden
	}

	return nil, nil
}
