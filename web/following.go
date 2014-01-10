package web

import (
	// "github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/repository"
	"net/http"
	// "net/url"
)

type FollowingResource struct {
	Repo   repository.UserRepository
	Auther Auther
}

func (resource FollowingResource) Follow(w http.ResponseWriter, req *http.Request) {
	user := resource.Repo.FindByUsername(req.URL.Query().Get("username"))
	other := resource.Repo.FindByUsername(req.URL.Query().Get("other"))

	currentUser := resource.Auther.Auth(req.Header)
	if currentUser == nil {
		w.WriteHeader(401)
		return
	}
	if currentUser.Id != user.Id {
		w.WriteHeader(403)
		return
	}

	resource.Repo.Follow(user, other)

	w.WriteHeader(201)
}
