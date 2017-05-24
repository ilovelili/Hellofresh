package main

import (
	"encoding/json"
	"hellofresh/config"
	"hellofresh/dal"
	"hellofresh/util"
	"log"
	"net/http"
	"strconv"
	"time"

	"hellofresh/model"

	"github.com/gorilla/mux"
)

// App the app container
type App struct {
	Router *mux.Router
	DB     interface{}
	Config *config.Config
}

// Enviroment enviroment
type Enviroment int

const (
	// Prod production
	Prod Enviroment = iota + 1
	// Test test
	Test
)

// Initialize init the app
func (app *App) Initialize(enviroment Enviroment) {
	config, err := config.GetConfig()
	if err != nil {
		util.PanicOnError(err)
	}
	app.Config = config

	// open dal
	if enviroment == Prod {
		if app.DB, err = dal.Open(config.ProductionDBConfig); err != nil {
			util.PanicOnError(err)
		}
	} else {
		if app.DB, err = dal.Open(config.TestDBConfig); err != nil {
			util.PanicOnError(err)
		}
	}

	// set up new router
	app.Router = mux.NewRouter()
	// init routes
	app.initializeRoutes()
}

// Run ListenAndServe
func (app *App) Run(addr string) {
	// set timeout to 15 seconds
	srv := &http.Server{
		Handler:      app.Router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// initializeRoutes init routes
func (app *App) initializeRoutes() {
	// alive check
	// GET / | non-protected
	app.Router.HandleFunc("/", app.aliveCheck).Methods("GET")

	// create recipe
	// POST /recipes | basic auth
	app.Router.HandleFunc("/recipes", util.Use(app.createRecipe, util.BasicAuth)).Methods("POST")

	// get single recipe
	// GET /recipes/{id} | non-protected
	app.Router.HandleFunc("/recipes/{id}", app.getRecipe).Methods("GET")

	// get recipe list
	// GET /recipes/{start:[0-9]+}/{limit:[0-9]+} | non-protected
	app.Router.HandleFunc("/recipes/{start:[0-9]+}/{limit:[0-9]+}", app.getRecipes).Methods("GET")

	// update recipe
	// PUT /recipes/{id} | basic auth
	app.Router.HandleFunc("/recipes/{id}", util.Use(app.updateRecipe, util.BasicAuth)).Methods("PUT")

	// delete recipe
	// DELETE /recipes/{id} | basic auth
	app.Router.HandleFunc("/recipes/{id}", util.Use(app.deleteRecipe, util.BasicAuth)).Methods("DELETE")

	// rate recipe
	// PUT /recipes/{id}/rate/{rate:[1-5]} | basic auth
	app.Router.HandleFunc("/recipes/{id}/rate/{rate:[1-5]}", util.Use(app.rateRecipe, util.BasicAuth)).Methods("PUT")

	// search recipe by name
	// GET /recipes/search/{name} | non-protected
	app.Router.HandleFunc("/recipes/search/{search:.+}", app.searchRecipes).Methods("GET")
}

// main app entry
func main() {
	app := &App{}
	app.Initialize(Prod)
	app.Run(":8080")
}

// aliveCheck GET /
func (app *App) aliveCheck(w http.ResponseWriter, r *http.Request) {
	util.ResponseWithJSON(w, http.StatusOK, "alive")
}

// getRecipes GET /recipes/{start}/{limit}
func (app *App) getRecipes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// pagination
	// default starts from 0 and take 10 records
	start, err := strconv.Atoi(vars["start"])
	if err != nil {
		start = 0
	}
	limit, err := strconv.Atoi(vars["limit"])
	if err != nil {
		limit = 10
	}

	recipe := &model.Recipe{}
	recipes, err := recipe.GetRecipes(app.DB, start, limit)
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.ResponseWithJSON(w, http.StatusOK, recipes)
}

// createRecipe POST /recipes
func (app *App) createRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe model.Recipe
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recipe); err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := recipe.CreateRecipe(app.DB); err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.ResponseWithJSON(w, http.StatusCreated, recipe)
}

// getRecipe GET /recipes/{id}
func (app *App) getRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := (model.ID)(vars["id"])
	recipe, err := id.GetRecipe(app.DB)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseWithJSON(w, http.StatusOK, recipe)
}

// updateRecipe PUT /recipes/{id}
func (app *App) updateRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipe := &model.Recipe{ID: vars["id"]}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recipe); err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()
	if err := recipe.UpdateRecipe(app.DB); err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.ResponseWithJSON(w, http.StatusOK, recipe)
}

// deleteRecipe DELETE /recipes/{id}
func (app *App) deleteRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := (model.ID)(vars["id"])
	if err := id.DeleteRecipe(app.DB); err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.ResponseWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// rateRecipe PUT /recipes/{id}/rate
func (app *App) rateRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := (model.ID)(vars["id"])
	rate, err := strconv.Atoi(vars["rate"])
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := id.RateRecipe(app.DB, rate); err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.ResponseWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// searchRecipes GET /recipes/search/{name}
func (app *App) searchRecipes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	search := vars["search"]
	if search == "" {
		util.ResponseWithError(w, http.StatusInternalServerError, "No search pattern")
		return
	}

	recipe := &model.Recipe{}
	recipes, err := recipe.SearchRecipes(app.DB, search)
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.ResponseWithJSON(w, http.StatusOK, recipes)
}
