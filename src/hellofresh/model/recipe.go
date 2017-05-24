// Package model Data Model
package model

import (
	"hellofresh/config"
	"hellofresh/util"
	"time"
)

// Difficulty three level of difficulties
type Difficulty int

// ID recipe id in string
type ID string

const (
	// Easy easy
	Easy Difficulty = iota + 1
	// Normal normal
	Normal
	// Hard hard
	Hard
)

// Recipe recipe entity
type Recipe struct {
	// ID can be string or bson.ObjectId
	ID         interface{} `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string      `json:"name"`
	Prep       time.Time   `json:"prep"`
	Difficulty Difficulty  `json:"difficulty"`
	Vegetarian bool        `json:"vegetarian"`
}

// accessor database accessor
var accessor RecipeRestFulAccessor

// init init accessor
func init() {
	config, err := config.GetConfig()
	if err != nil {
		util.PanicOnError(err)
	}
	accessor, err = GetAccessor(config.ProductionDBConfig.Host)
	if err != nil {
		util.PanicOnError(err)
	}
}

// GetRecipe get single recipe
func (id *ID) GetRecipe(db interface{}) (*Recipe, error) {
	return accessor.Get(db, id)
}

// UpdateRecipe update single recipe
func (recipe *Recipe) UpdateRecipe(db interface{}) error {
	return accessor.Update(db, recipe)
}

// DeleteRecipe delete single recipe
func (id *ID) DeleteRecipe(db interface{}) error {
	return accessor.Delete(db, id)
}

// CreateRecipe create single recipe
func (recipe *Recipe) CreateRecipe(db interface{}) error {
	return accessor.Create(db, recipe)
}

// GetRecipes get recipe list
func (recipe *Recipe) GetRecipes(db interface{}, start, limit int) ([]*Recipe, error) {
	return accessor.List(db, start, limit)
}

// RateRecipe rate recipe
func (id *ID) RateRecipe(db interface{}, rate int) error {
	return accessor.Rate(db, id, rate)
}

// SearchRecipes search recipes
func (recipe *Recipe) SearchRecipes(db interface{}, search string) ([]*Recipe, error) {
	return accessor.Search(db, search)
}
