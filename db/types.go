package db

import (
	"github.com/go-bongo/bongo"
)

//Redirect represents a redirect object in our database
type Redirect struct {
	bongo.DocumentBase `bson:",inline"`
	Path               string
	URL                string
	Password           string
	TTL                string
}
