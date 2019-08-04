package db

import (
	"errors"

	"github.com/go-bongo/bongo"
	"github.com/spf13/viper"
)

const (
	defaultPath    = ".env"
	CollectionName = "redirections"
	connString     = "CONNECTION_STRING"
)

var config *bongo.Config
var connection *bongo.Connection

func loadConfig(configPath string) error {
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath(defaultPath)
	viper.SetConfigType("dotenv")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if !(viper.IsSet(connString)) {
		return errors.New("Database configuration missing connection string or database name")
	}
	config = &bongo.Config{ConnectionString: viper.GetString(connString)}
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
