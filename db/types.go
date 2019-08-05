package db

import (
	"github.com/go-bongo/bongo"
)

type Redirect struct {
	bongo.DocumentBase `bson:",inline"`
	Path               string
	URL                string
	Password           string
	TTL                string
}
