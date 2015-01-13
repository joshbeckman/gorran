package main

import (
    "time"
    "labix.org/v2/mgo/bson"
)

// Models
type Account struct {
    Id bson.ObjectId `bson:"_id" json:"_id"`
    Username string
    Vanity string
    TunesCategories string
    Email string
}

type Article struct {
    Id bson.ObjectId `bson:"_id" json:"_id"`
    Title string
    Url string
    Mp3URL string `bson:"mp3URL"`
    Mp3Length float64 `bson:"mp3Length"`
    Description string
    accountId string
    Created time.Time
}