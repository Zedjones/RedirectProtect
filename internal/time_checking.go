package internal

import (
	"time"

	"github.com/getlantern/deepcopy"
	"github.com/zedjones/redirectprotect/db"
	"gopkg.in/mgo.v2/bson"
)

var (
	getConnection  = db.GetConnection
	copy           = deepcopy.Copy
	startTimeCheck = StartTimeCheck
	parseDuration  = time.ParseDuration
	sleep          = time.Sleep
)

func AddChecks() error {
	var err error
	connection, err := getConnection()
	if err != nil {
		return err
	}
	redir := &db.Redirect{}
	collection := connection.Collection(db.CollectionName)
	allRedirs := collection.Find(bson.M{})
	for allRedirs.Next(redir) {
		//deepcopy our redirection
		var redirCopy *db.Redirect
		err := copy(&redirCopy, &redir)
		if err != nil {
			return err
		}
		go startTimeCheck(redirCopy, collection)
	}
	return err
}

func StartTimeCheck(redir *db.Redirect, collection db.Collection) error {
	var err error
	ttl, err := parseDuration(redir.TTL)
	if err != nil {
		return err
	}
	if ttl.Seconds() == 0 {
		return err
	}
	completedTime := redir.Created.Add(ttl)
	timeLeft := completedTime.Sub(time.Now())
	sleep(timeLeft)
	err = collection.DeleteDocument(redir)
	return err
}
