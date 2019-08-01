package db

import "github.com/go-bongo/bongo"

type Redirect struct {
	bongo.DocumentBase `bson:",inline"`
	URL                string
	Password           string
}
