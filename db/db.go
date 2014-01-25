package db

import (
	"database/sql"
	_ "github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/bmizerany/pq"
	"github.com/michaelfairley/mapi-kata-tigertonic-gorp/github.com/coopernurse/gorp"
)

func Setup(url string) *gorp.DbMap {
	dbHandle, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	dbmap := &gorp.DbMap{Db: dbHandle, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Token{}, "tokens")
	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	return dbmap
}
