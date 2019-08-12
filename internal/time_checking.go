package internal

import (
	"time"

	"github.com/getlantern/deepcopy"
	"github.com/go-bongo/bongo"
	"github.com/zedjones/redirectprotect/db"
	"gopkg.in/mgo.v2/bson"
)

func AddChecks() error {
	var err error
	connection, err := db.GetConnection()
	if err != nil {
		return err
	}
	redir := &db.Redirect{}
	collection := connection.Collection(db.CollectionName)
	allRedirs := collection.Find(bson.M{})
	for allRedirs.Next(redir) {
		//deepcopy our redirection
		var redirCopy *db.Redirect
		err := deepcopy.Copy(&redirCopy, &redir)
		if err != nil {
			return err
		}
		go StartTimeCheck(redirCopy, collection)
	}
	return err
}

func StartTimeCheck(redir *db.Redirect, collection *bongo.Collection) error {
	var err error
	ttl, err := time.ParseDuration(redir.TTL)
	if err != nil {
		return err
	}
	if ttl.Seconds() == 0 {
		return err
	}
	completedTime := redir.Created.Add(ttl)
	timeLeft := completedTime.Sub(time.Now())
	time.Sleep(timeLeft)
	err = collection.DeleteDocument(redir)
	return err
}
