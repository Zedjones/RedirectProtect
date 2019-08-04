package db

import (
	"time"

	"github.com/go-bongo/bongo"
)

type Redirect struct {
	bongo.DocumentBase `bson:",inline"`
	Path               string
	URL                string
	Password           string
	TTL                *time.Duration
}
