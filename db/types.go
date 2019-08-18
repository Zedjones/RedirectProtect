package db

import (
	"github.com/go-bongo/bongo"
)

type Connection interface {
	Connect() error
	Collection(name string) Collection
}

type BongoConnection struct {
	*bongo.Connection
}

type Collection interface {
	Save(doc bongo.Document) error
	Find(query interface{}) *bongo.ResultSet
	FindOne(query interface{}, doc interface{}) error
	DeleteDocument(doc bongo.Document) error
}

func (c BongoConnection) Collection(name string) Collection {
	return &bongo.Collection{
		Connection: c.Connection,
		Name:       name,
	}
}

//Redirect represents a redirect object in our database
type Redirect struct {
	bongo.DocumentBase `bson:",inline"`
	Path               string
	URL                string
	Password           string
	TTL                string
}
