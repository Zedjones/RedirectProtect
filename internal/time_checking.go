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
	now            = time.Now
)

//AddChecks start a coroutine for each item in the database
//to check if it has a timeout and delete the item after its
//time is up
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

//StartTimeCheck waits until the timeout for the provided db.Redirect
//is done and then deletes the item from the database
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
	timeLeft := completedTime.Sub(now())
	sleep(timeLeft)
	return collection.DeleteDocument(redir)
}
