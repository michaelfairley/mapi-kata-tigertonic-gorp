package web

import (
	"errors"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/db"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/rcrowley/go-tigertonic"
	"net/http"
	"regexp"
)

var (
	unauthorized = tigertonic.Unauthorized{errors.New("")}
	forbidden    = tigertonic.Forbidden{errors.New("")}
	tokenRegex   = regexp.MustCompile("^Token (.*)$")
)

type Auther struct {
	UserRepo  db.UserRepository
	TokenRepo db.TokenRepository
}

func (auther Auther) Auth(headers http.Header) *db.User {
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
