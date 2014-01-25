package main

import (
	"database/sql"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/db"
	_ "github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/bmizerany/pq"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/rcrowley/go-tigertonic"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/web"
	"log"
	"net/http"
	"os"
)

func setupDB(url string) *gorp.DbMap {
	dbHandle, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	dbmap := &gorp.DbMap{Db: dbHandle, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(db.User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(db.Token{}, "tokens")
	dbmap.AddTableWithName(db.Post{}, "posts").SetKeys(true, "Id")

	return dbmap
}

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

	dbMap := setupDB(c.Database)

	mux := setupMux(dbMap)

	server := tigertonic.NewServer(":12346", tigertonic.ApacheLogged(mux))

	err := server.ListenAndServe()

	if nil != err {
		log.Fatalln(err)
	}
}
