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

func loadConfig(configPath string) error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if os.Getenv(connString) == "" {
		return errors.New("Database configuration missing connection string or database name")
	}
	config = &bongo.Config{ConnectionString: os.Getenv(connString)}
	return nil
}

func GetConnection() (*bongo.Connection, error) {
	if connection == nil {
		var err error
		if config == nil {
			if err = loadConfig(""); err != nil {
				return nil, err
			}
		}
		connection, err = bongo.Connect(config)
		return connection, err
	}
	return connection, nil
}
