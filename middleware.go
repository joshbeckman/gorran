package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
)

type Controller struct {
	session *mgo.Session
}

func NewController() (*Controller, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, fmt.Errorf("no DB connection string provided")
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}
	return &Controller{
		session: session,
	}, nil
}
