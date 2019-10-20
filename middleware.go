package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji/web"
	"net/http"
	"os"
)

type Controller struct {
	dbURI string
}

func NewController() (*Controller, error) {
	uri := os.Getenv("DATABASE_URL")
	if uri == "" {
		return nil, fmt.Errorf("no DB connection string provided")
	}
	return &Controller{
		dbURI: uri,
	}, nil
}

func findRequestAccountByVanity(c web.C, r *http.Request, session *sqlx.DB) (Account, error) {
	result := Account{}
	_, pass, ok := r.BasicAuth()
	if !ok {
		accountErr := session.QueryRowx("SELECT _id, username, image, vanity, \"itunesCategories\", email, \"podcastTitle\", \"podcastDescription\" FROM accounts WHERE vanity = $1 AND \"podcastPassword\" is NULL", c.URLParams["vanity"]).StructScan(&result)
		return result, accountErr
	}
	accountErr := session.QueryRowx("SELECT _id, username, image, vanity, \"itunesCategories\", email, \"podcastTitle\", \"podcastDescription\" FROM accounts WHERE vanity = $1 AND \"podcastPassword\" = $2", c.URLParams["vanity"], pass).StructScan(&result)
	return result, accountErr
}
