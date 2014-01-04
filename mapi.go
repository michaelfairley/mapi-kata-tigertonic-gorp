package main

import (
	"database/sql"
	"github.com/michaelfairley/mapi-tigertonic-gorp/api"
	_ "github.com/michaelfairley/mapi-tigertonic-gorp/github.com/bmizerany/pq"
	"github.com/michaelfairley/mapi-tigertonic-gorp/github.com/coopernurse/gorp"
	"github.com/michaelfairley/mapi-tigertonic-gorp/github.com/rcrowley/go-tigertonic"
	"github.com/michaelfairley/mapi-tigertonic-gorp/repository"
	"github.com/michaelfairley/mapi-tigertonic-gorp/web"
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

	return dbmap
}

func setupMux(db *gorp.DbMap) http.Handler {
	mux := tigertonic.NewTrieServeMux()

	userResource := web.UserResource{repository.UserRepository{db}}
	mux.Handle("POST", "/users", tigertonic.Marshaled(userResource.CreateUser))
	mux.Handle("GET", "/users/{username}", tigertonic.Marshaled(userResource.GetUser))

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
