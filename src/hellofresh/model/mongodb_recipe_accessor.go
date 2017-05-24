package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDBAccessor MongoDB restful accessor
type MongoDBAccessor struct{}

// Description Description
func (accessor *MongoDBAccessor) Description() string {
	return "mongodb restful accessor"
}

// Get get recipe
func (accessor *MongoDBAccessor) Get(db interface{}, id *ID) (*Recipe, error) {
	recipe := Recipe{}
	collection := db.(*mgo.Database).C("recipe")
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(fmt.Sprintf("%s", *id))}).One(&recipe)
	return &recipe, err
}

// Update update recipe
func (accessor *MongoDBAccessor) Update(db interface{}, recipe *Recipe) error {
	collection := db.(*mgo.Database).C("recipe")
	colQuerier := bson.M{"_id": bson.ObjectIdHex(recipe.ID.(string))}
	change := bson.M{"$set": bson.M{"name": recipe.Name, "prep": recipe.Prep, "difficulty": recipe.Difficulty, "vegetarian": recipe.Vegetarian}}
	return collection.Update(colQuerier, change)
}

// Delete delete recipe
func (accessor *MongoDBAccessor) Delete(db interface{}, id *ID) error {
	collection := db.(*mgo.Database).C("recipe")
	return collection.Remove(bson.M{"_id": bson.ObjectIdHex(fmt.Sprintf("%s", *id))})
}

// Create create recipe
func (accessor *MongoDBAccessor) Create(db interface{}, recipe *Recipe) error {
	collection := db.(*mgo.Database).C("recipe")
	return collection.Insert(&Recipe{Name: recipe.Name, Prep: recipe.Prep, Difficulty: recipe.Difficulty, Vegetarian: recipe.Vegetarian})
}

// List get recipe list
func (accessor *MongoDBAccessor) List(db interface{}, start, limit int) ([]*Recipe, error) {
	var recipes []*Recipe
	collection := db.(*mgo.Database).C("recipe")
	err := collection.Find(nil).Skip(start).Limit(limit).All(&recipes)
	return recipes, err
}

// Rate rate recipe
func (accessor *MongoDBAccessor) Rate(db interface{}, id *ID, rate int) error {
	collection := db.(*mgo.Database).C("reciperate")
	return collection.Insert(&RecipeRate{RecipeID: fmt.Sprintf("%s", *id), Rate: rate, User: "Jane Doe" /*dummy or use ip*/, Modified: time.Now()})
}

// Search search recipe by search pattern
func (accessor *MongoDBAccessor) Search(db interface{}, search string) ([]*Recipe, error) {
	var recipes []*Recipe
	collection := db.(*mgo.Database).C("recipe")
	regex := bson.M{"$regex": bson.RegEx{Pattern: search}}
	err := collection.Find(bson.M{"name": regex}).All(&recipes)
	return recipes, err
}
