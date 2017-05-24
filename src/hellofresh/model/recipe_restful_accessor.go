package model

import (
	"errors"
	"strings"
)

// RecipeRestFulAccessor db accessor interface
type RecipeRestFulAccessor interface {
	Description() string
	List(db interface{}, start, limit int) ([]*Recipe, error)
	Create(db interface{}, recipe *Recipe) error
	Get(db interface{}, id *ID) (*Recipe, error)
	Update(db interface{}, recipe *Recipe) error
	Delete(db interface{}, id *ID) error
	Rate(db interface{}, id *ID, rate int) error
	Search(db interface{}, search string) ([]*Recipe, error)
}

// GetAccessor get accessor by client
func GetAccessor(client string) (accessor RecipeRestFulAccessor, err error) {
	switch strings.ToLower(client) {
	case "postgres":
		accessor = &PostGresAccessor{}
		return
	case "mongodb":
		accessor = &MongoDBAccessor{}
		return
	default:
		err = errors.New("Not supportted")
		return
	}
}
