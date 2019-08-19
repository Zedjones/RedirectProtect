package db

import (
	"github.com/go-bongo/bongo"
)

//Connection represents a Bongo connection
type Connection interface {
	Connect() error
	Collection(name string) Collection
}

//BongoConnection is just a wrapper for *bongo.Connection for mocking
type BongoConnection struct {
	*bongo.Connection
}

//Collection is an interface that specifies all methods from
//*bongo.Collection that we need/need to mock
type Collection interface {
	Save(doc bongo.Document) error
	Find(query interface{}) *bongo.ResultSet
	FindOne(query interface{}, doc interface{}) error
	DeleteDocument(doc bongo.Document) error
}

//Collection wraps *bongo.Connection.Collection()
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
