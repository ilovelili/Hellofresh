// Package config parse config from config.json
package config

import (
	"encoding/json"
	"os"
	"path"
)

// GetConfig get config defined in config.json
func GetConfig() (config *Config, err error) {
	configpath := path.Join(os.Getenv("GOPATH"), "src", "hellofresh", "config.json")
	configFile, err := os.Open(configpath)
	defer configFile.Close()

	if err != nil {
		return
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return
	}

	return
}

// DBConfigFields database config
type DBConfigFields struct {
	Host     string `json:"host"`
	Server   string `json:"server"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

// DBConfig database config including different enviroments
type DBConfig struct {
	ProductionDBConfig *DBConfigFields `json:"prod"`
	TestDBConfig       *DBConfigFields `json:"test"`
}

// AuthConfig auth config
type AuthConfig struct {
	Type     string `json:"type"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Config config entry
type Config struct {
	DBConfig   `json:"db"`
	AuthConfig `json:"auth"`
}
