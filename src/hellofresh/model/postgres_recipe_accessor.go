package model

import (
	"database/sql"
	"fmt"
	"hellofresh/config"
	"hellofresh/dal"
	"hellofresh/util"
	"log"
	"strings"
	"time"
)

// PostGresAccessor PostGres restful accessor
type PostGresAccessor struct{}

func init() {
	config, err := config.GetConfig()
	if err != nil {
		util.PanicOnError(err)
	}

	// do nothing if database is not postgres
	if strings.ToLower(config.ProductionDBConfig.Host) != "postgres" {
		return
	}

	// open
	db, err := dal.Open(config.ProductionDBConfig)
	if err != nil {
		util.PanicOnError(err)
	}

	ensureTableExists(db.(*sql.DB))
}

// ensureTableExists make sure recipes table exists when use postgre
func ensureTableExists(db *sql.DB) {
	if _, err := db.Exec(recipeTableCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(recipeRateTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

const recipeTableCreationQuery = `CREATE TABLE IF NOT EXISTS recipes
(
	id SERIAL,
	name TEXT NOT NULL,
	prep TIMESTAMP NOT NULL DEFAULT now(),
	difficulty INT NOT NULL,
	vegetarian BOOLEAN NOT NULL,
	CONSTRAINT recipes_pkey PRIMARY KEY (id)
)`

const recipeRateTableCreationQuery = `CREATE TABLE IF NOT EXISTS reciperates
(
	id SERIAL,
	recipeId TEXT NOT NULL,
	rate INT NOT NULL,
	rateuser VARCHAR(100) NOT NULL,
	modified TIMESTAMP NOT NULL,
	CONSTRAINT reciperates_pkey PRIMARY KEY (id)
)`

// Description Description
func (accessor *PostGresAccessor) Description() string {
	return "postgres restful accessor"
}

// Get get single recipe
func (accessor *PostGresAccessor) Get(db interface{}, id *ID) (*Recipe, error) {
	recipe := Recipe{}
	err := db.(*sql.DB).QueryRow("SELECT name, prep, difficulty, vegetarian FROM recipes WHERE id=$1", fmt.Sprintf("%s", *id)).Scan(&recipe.Name, &recipe.Prep, &recipe.Difficulty, &recipe.Vegetarian)
	return &recipe, err
}

// Update update single recipe
func (accessor *PostGresAccessor) Update(db interface{}, recipe *Recipe) error {
	_, err := db.(*sql.DB).Exec("UPDATE recipes SET name=$1, prep=$2, difficulty=$3, vegetarian=$4 WHERE id=$5", recipe.Name, recipe.Prep, recipe.Difficulty, recipe.Vegetarian, recipe.ID)
	return err
}

// Delete delete single recipe
func (accessor *PostGresAccessor) Delete(db interface{}, id *ID) error {
	_, err := db.(*sql.DB).Exec("DELETE FROM recipes WHERE id=$1", fmt.Sprintf("%s", *id))
	return err
}

// Create create single recipe
func (accessor *PostGresAccessor) Create(db interface{}, recipe *Recipe) error {
	return db.(*sql.DB).QueryRow("INSERT INTO recipes(name, prep, difficulty, vegetarian) VALUES($1, $2, $3, $4) RETURNING id", recipe.Name, recipe.Prep, recipe.Difficulty, recipe.Vegetarian).Scan(&recipe.ID)
}

// List get recipe list
func (accessor *PostGresAccessor) List(db interface{}, start, limit int) ([]*Recipe, error) {
	recipes := []*Recipe{}
	rows, err := db.(*sql.DB).Query("SELECT id, name, prep, difficulty, vegetarian FROM recipes LIMIT $1 OFFSET $2", limit, start)
	defer rows.Close()
	if err != nil {
		return recipes, err
	}

	for rows.Next() {
		recipe := Recipe{}
		if err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.Prep, &recipe.Difficulty, &recipe.Vegetarian); err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}

// Rate rate recipe
func (accessor *PostGresAccessor) Rate(db interface{}, id *ID, rate int) error {
	_, err := db.(*sql.DB).Exec("INSERT INTO reciperates(recipeId, rate, rateuser, modified) VALUES($1, $2, $3, $4) RETURNING id", fmt.Sprintf("%s", *id), rate, "Jane Doe" /*dummy or use ip*/, time.Now())
	return err
}

// Search search recipes
func (accessor *PostGresAccessor) Search(db interface{}, search string) ([]*Recipe, error) {
	recipes := []*Recipe{}
	rows, err := db.(*sql.DB).Query("SELECT id, name, prep, difficulty, vegetarian FROM recipes where name LIKE '%' || $1 || '%'", search)
	defer rows.Close()
	if err != nil {
		return recipes, err
	}

	for rows.Next() {
		recipe := Recipe{}
		if err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.Prep, &recipe.Difficulty, &recipe.Vegetarian); err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}
