package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/jmoiron/sqlx/types"
	"regexp"
	"strings"
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
	PodcastPassword    sql.NullString `db:"podcastPassword"`
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

func (article *Article) enclosureURL(acct Account) string {
	wavExt, _ := regexp.MatchString(".wav$", article.Mp3URL)
	extension := ".mp3"
	if wavExt {
	    extension = ".wav"
	}
	if !acct.PodcastPassword.Valid {
		return strings.Join([]string{"https://www.narro.co/article/", article.Id, extension}, "")
	}
	if acct.PodcastPassword.Valid {
		if acct.PodcastPassword.String == "" {
			return strings.Join([]string{"https://www.narro.co/article/", article.Id, extension}, "")
		}
	}
	h := hmac.New(sha256.New, []byte(acct.PodcastPassword.String))
	h.Write([]byte(article.Id))
	token := strings.Join([]string{"?token=", hex.EncodeToString(h.Sum(nil))}, "")
	return strings.Join([]string{"https://www.narro.co/article/", article.Id, extension, token}, "")
}

type ArticleLink struct {
	Id   string `db:"_id" json:"_id"`
	Href string
	Text string
}
