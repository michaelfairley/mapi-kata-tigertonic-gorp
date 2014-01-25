package web

import (
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/code.google.com/p/go.crypto/bcrypt"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/db"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/utils"
	"net/http"
	"net/url"
)

type UserResource struct {
	Repository db.UserRepository
}

func webUserToDb(web *api.User) *db.User {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(web.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.CheckErr(err)
	}

	db := db.User{
		Realname: web.Realname,
		Password: cryptedPassword,
		Username: web.Username,
	}

	return &db
}

func dbUserToWeb(db *db.User) *api.User {
	api := api.User{
		Username: db.Username,
		Realname: db.Realname,
	}

	return &api
}

func (resource *UserResource) CreateUser(url *url.URL, inHeaders http.Header, user *api.User) (int, http.Header, interface{}, error) {
	existing := resource.Repository.ContainsUserWithUsername(user.Username)

	if existing {
		return validationError(map[string][]string{"username": []string{"is taken"}})
	}
	if len(user.Password) < 8 {
		return validationError(map[string][]string{"password": []string{"is too short"}})
	}

	resource.Repository.Insert(webUserToDb(user))

	headers := http.Header{
		"Location": []string{"http://localhost:12346/users/" + user.Username},
	}
	return 303, headers, nil, nil
}

func (resource *UserResource) GetUser(url *url.URL, inHeaders http.Header, _ interface{}) (int, http.Header, *api.User, error) {
	user := resource.Repository.FindByUsername(url.Query().Get("username"))
	if user == nil {
		return 404, nil, nil, nil
	}

	followers := resource.Repository.FindFollowers(user)
	following := resource.Repository.FindFollowing(user)

	followerNames := make([]string, len(followers))
	for i := range followers {
		followerNames[i] = followers[i].Username
	}

	followingNames := make([]string, len(following))
	for i := range following {
		followingNames[i] = following[i].Username

	}

	webUser := dbUserToWeb(user)

	webUser.Followers = followerNames
	webUser.Following = followingNames

	return 200, http.Header{}, webUser, nil
}
