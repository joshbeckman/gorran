package main

import (
	"labix.org/v2/mgo/bson"
	"time"
)

// Models
type Account struct {
	Id              bson.ObjectId `bson:"_id" json:"_id"`
	Username        string
	Image           string
	Vanity          string
	TunesCategories string `bson:"itunesCategories"`
	Email           string
}

type Article struct {
	Id          bson.ObjectId `bson:"_id" json:"_id"`
	Title       string
	Url         string
	Mp3URL      string  `bson:"mp3URL"`
	Mp3Length   float64 `bson:"mp3Length"`
	Description string
	AccountId   string
	Created     time.Time
	Links       []*ArticleLink
}

type ArticleLink struct {
	Id   bson.ObjectId `bson:"_id" json:"_id"`
	Href string
	Text string
}
