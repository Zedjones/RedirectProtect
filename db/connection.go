package db

import (
	"errors"
	"os"

	"github.com/go-bongo/bongo"
	"github.com/joho/godotenv"
)

const (
	defaultPath    = ".env"
	CollectionName = "redirections"
	connString     = "CONNECTION_STRING"
)

var config *bongo.Config
var connection *bongo.Connection

func loadConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	loadedConnString := os.Getenv(connString)
	if loadedConnString == "" {
		return errors.New("Database configuration missing connection string or database name")
	}
	config = &bongo.Config{ConnectionString: loadedConnString}
	return nil
}

func GetConnection() (Connection, error) {
	if connection == nil {
		var err error
		if config == nil {
			if err = loadConfig(); err != nil {
				return nil, err
			}
		}
		connection, err = bongo.Connect(config)
		return BongoConnection{connection}, err
	}
	return BongoConnection{connection}, nil
}
