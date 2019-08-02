package db

import (
	"github.com/go-bongo/bongo"
	"github.com/spf13/viper"
)

const defaultPath = "config/db.json"

var config = &bongo.Config{
	ConnectionString: "172.19.0.2",
}

func loadConfig(configPath string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")
}
