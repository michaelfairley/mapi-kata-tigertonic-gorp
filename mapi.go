package main

import (
	"database/sql"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/api"
	_ "github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/bmizerany/pq"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/rcrowley/go-tigertonic"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/repository"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/web"
	"log"
	"net/http"
	"os"
)

func setupDB(url string) *gorp.DbMap {
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(api.User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(api.Token{}, "tokens")
	dbmap.AddTableWithName(api.DbPost{}, "posts").SetKeys(true, "Id")

	return dbmap
}

func setupMux(db *gorp.DbMap) http.Handler {
	mux := tigertonic.NewTrieServeMux()

	userRepository := repository.UserRepository{db}
	tokenRepository := repository.TokenRepository{db}
	postRepository := repository.PostRepository{db}

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

	db := setupDB(c.Database)

	mux := setupMux(db)

	server := tigertonic.NewServer(":12346", tigertonic.ApacheLogged(mux))

	err := server.ListenAndServe()

	if nil != err {
		log.Fatalln(err)
	}
}
