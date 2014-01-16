package web

import (
	// "github.com/michaelfairley/mapi-kata-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/repository"
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

func (resource FollowingResource) Unfollow(w http.ResponseWriter, req *http.Request) {
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

	any := resource.Repo.Unfollow(user, other)
	if any {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(404)
	}
}
