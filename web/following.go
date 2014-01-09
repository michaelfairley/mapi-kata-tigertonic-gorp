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

	resource.Repo.Follow(user, other)

	w.WriteHeader(201)
}
