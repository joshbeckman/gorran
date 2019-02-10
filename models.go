package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx/types"
	"time"
)

// Models
type Account struct {
	Id                 string `db:"_id" json:"_id"`
	Username           string
	Image              sql.NullString
	Vanity             string
	TunesCategories    sql.NullString `db:"itunesCategories"`
	Email              string
	PodcastTitle       sql.NullString `db:"podcastTitle"`
	PodcastDescription sql.NullString `db:"podcastDescription"`
}

type Article struct {
	Id          string `db:"_id" json:"_id"`
	Title       string
	Url         sql.NullString
	Mp3URL      string  `db:"mp3URL"`
	Mp3Length   float64 `db:"mp3Length"`
	Description string
	AccountId   string `db:"accountId"`
	Created     time.Time
	Links       types.JSONText
}

type ArticleLink struct {
	Id   string `db:"_id" json:"_id"`
	Href string
	Text string
}
