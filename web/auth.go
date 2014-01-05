package web

import (
	"errors"
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

func (auther Auther) Auth(request *http.Request) (http.Header, error) {
	authentication := request.Header["Authentication"]
	if len(authentication) != 1 {
		return nil, unauthorized
	}
	if !tokenRegex.MatchString(authentication[0]) {
		return nil, unauthorized
	}
	tokenValue := authentication[0][6:]
	token := auther.TokenRepo.FindByValue(tokenValue)
	if token == nil {
		return nil, unauthorized
	}

	user := auther.UserRepo.FindByUsername(request.URL.Query().Get("username"))

	if user == nil {
		return nil, forbidden
	}

	if user.Id != token.UserId {
		return nil, forbidden
	}

	return nil, nil
}
