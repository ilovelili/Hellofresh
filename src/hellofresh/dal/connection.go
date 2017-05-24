// Package dal Data Access Layer
package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"hellofresh/config"
	"strings"

	mgo "gopkg.in/mgo.v2"

	_ "github.com/lib/pq"
)

// Open open database connection by config
func Open(config *config.DBConfigFields) (interface{}, error) {
	host := strings.ToLower(config.Host)
	switch host {
	case "postgres":
		connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Server, config.UserName, config.Password, config.DBName)
		database, err := sql.Open("postgres", connectionString)
		return database, err

	case "mongodb":
		session, err := mgo.Dial(config.Server)
		if err != nil {
			return nil, err
		}
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		return session.DB(config.DBName), nil

	default:
		return nil, errors.New("Unsupportted host")
	}
}
