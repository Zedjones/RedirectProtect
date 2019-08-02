package db

import (
	"errors"

	"github.com/go-bongo/bongo"
	"github.com/spf13/viper"
)

const (
	defaultPath = "config/db.json"
	connString  = "connection_string"
	dbKey       = "database"
)

var config *bongo.Config

func loadConfig(configPath string) error {
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath(defaultPath)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if !(viper.IsSet(connString) && viper.IsSet(dbKey)) {
		return errors.New("Database configuration missing connection string or database name")
	}
	config = &bongo.Config{ConnectionString: viper.GetString(connString), Database: viper.GetString(dbKey)}
	return nil
}

func getConnection() (*bongo.Connection, error) {
	if config == nil {
		if err := loadConfig(""); err != nil {
			return nil, err
		}
	}
	return bongo.Connect(config)
}
