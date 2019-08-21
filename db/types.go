package db

import (
	"github.com/go-bongo/bongo"
)

//Connection represents a Bongo connection
type Connection interface {
	Connect() error
	Collection(name string) Collection
}

//Collection is an interface that specifies all methods from
//*bongo.Collection that we need/need to mock
type Collection interface {
	Save(doc bongo.Document) error
	Find(query interface{}) ResultSet
	FindOne(query interface{}, doc interface{}) error
	DeleteDocument(doc bongo.Document) error
}

//BongoConnection is just a wrapper for *bongo.Connection for mocking
type BongoConnection struct {
	*bongo.Connection
}

//Collection wraps *bongo.Connection.Collection()
func (c BongoConnection) Collection(name string) Collection {
	return BongoCollection{
		&bongo.Collection{
			Connection: c.Connection,
			Name:       name,
		},
	}
}

type BongoCollection struct {
	*bongo.Collection
}

func (c BongoCollection) Find(query interface{}) ResultSet {
	return c.Collection.Find(query)
}

type ResultSet interface {
	Next(doc interface{}) bool
}

//Redirect represents a redirect object in our database
type Redirect struct {
	bongo.DocumentBase `bson:",inline"`
	Path               string
	URL                string
	Password           string
	TTL                string
}
