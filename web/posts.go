package web

import (
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	"github.com/michaelfairley/mapi-tigertonic-gorp/repository"
	"net/http"
	"net/url"
	"strconv"
)

type PostsResource struct {
	Repo     repository.PostRepository
	UserRepo repository.UserRepository
	Auther   Auther
}

func (resource *PostsResource) CreatePost(url *url.URL, inHeaders http.Header, post *api.Post) (int, http.Header, interface{}, error) {
	user := resource.UserRepo.FindByUsername(url.Query().Get("username"))
	dbPost := &api.DbPost{UserId: user.Id, Text: post.Text}

	resource.Repo.Insert(dbPost)

	headers := http.Header{
		"Location": []string{"http://localhost:12346/posts/" + strconv.FormatUint(dbPost.Id, 10)},
	}

	return 303, headers, nil, nil
}

func (resource *PostsResource) GetPost(url *url.URL, inHeaders http.Header, _ interface{}) (int, http.Header, *api.Post, error) {
	id, err := strconv.ParseUint(url.Query().Get("id"), 10, 64)
	if err != nil {
		return 404, nil, nil, nil
	}

	dbPost := resource.Repo.Find(id)
	if dbPost == nil {
		return 404, nil, nil, nil
	}

	user := resource.UserRepo.Find(dbPost.UserId)
	post := &api.Post{Author: user.Username, Text: dbPost.Text, Id: dbPost.Id}

	return 200, http.Header{}, post, nil
}

func (resource *PostsResource) DeletePost(url *url.URL, inHeaders http.Header, _ interface{}) (int, http.Header, interface{}, error) {
	id, err := strconv.ParseUint(url.Query().Get("id"), 10, 64)
	if err != nil {
		return 404, nil, nil, nil
	}

	dbPost := resource.Repo.Find(id)
	if dbPost == nil {
		return 404, nil, nil, nil
	}

	author := resource.Auther.Auth(inHeaders)

	if author == nil {
		return 401, nil, nil, nil
	}
	if author.Id != dbPost.UserId {
		return 403, nil, nil, nil
	}

	resource.Repo.Delete(dbPost)

	return 204, http.Header{}, nil, nil
}

func (resource *PostsResource) ListPosts(url *url.URL, inHeaders http.Header, _ interface{}) (int, http.Header, *api.Posts, error) {
	user := resource.UserRepo.FindByUsername(url.Query().Get("username"))

	dbPosts := resource.Repo.FindByUserId(user.Id)

	posts := make([]*api.Post, len(dbPosts))

	for i := range dbPosts {
		posts[i] = &api.Post{Author: user.Username, Text: dbPosts[i].Text, Id: dbPosts[i].Id}
	}

	return 200, http.Header{}, &api.Posts{Posts: posts}, nil
}
