package main

import (
	"fmt"
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
