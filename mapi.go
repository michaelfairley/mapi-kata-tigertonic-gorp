package main

import (
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/db"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/rcrowley/go-tigertonic"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/web"
	"log"
	"net/http"
	"os"
)

func setupMux(dbMap *gorp.DbMap) http.Handler {
	mux := tigertonic.NewTrieServeMux()

	userRepository := db.UserRepository{dbMap}
	tokenRepository := db.TokenRepository{dbMap}
	postRepository := db.PostRepository{dbMap}

	auther := web.Auther{userRepository, tokenRepository}

	userResource := web.UserResource{userRepository}
	mux.Handle("POST", "/users", tigertonic.Marshaled(userResource.CreateUser))
	mux.Handle("GET", "/users/{username}", tigertonic.Marshaled(userResource.GetUser))

	tokensResource := web.TokensResource{tokenRepository, userRepository}
	mux.Handle("POST", "/tokens", tigertonic.Marshaled(tokensResource.CreateToken))

	postsResource := web.PostsResource{postRepository, userRepository, auther}
	mux.Handle("POST", "/users/{username}/posts", tigertonic.If(auther.Filter, tigertonic.Marshaled(postsResource.CreatePost)))
	mux.Handle("GET", "/posts/{id}", tigertonic.Marshaled(postsResource.GetPost))
	mux.Handle("DELETE", "/posts/{id}", tigertonic.Marshaled(postsResource.DeletePost))
	mux.Handle("GET", "/users/{username}/posts", tigertonic.Marshaled(postsResource.ListPosts))
	mux.Handle("GET", "/users/{username}/timeline", tigertonic.Marshaled(postsResource.ShowTimeline))

	followingResource := web.FollowingResource{userRepository, auther}
	mux.HandleFunc("PUT", "/users/{username}/following/{other}", followingResource.Follow)
	mux.HandleFunc("DELETE", "/users/{username}/following/{other}", followingResource.Unfollow)

	return mux
}

func main() {
	c := &Config{}
	if err := tigertonic.Configure(os.Args[1], c); nil != err {
		log.Fatalln(err)
	}

	dbMap := db.Setup(c.Database)

	mux := setupMux(dbMap)

	server := tigertonic.NewServer(":12346", tigertonic.ApacheLogged(mux))

	err := server.ListenAndServe()

	if nil != err {
		log.Fatalln(err)
	}
}
