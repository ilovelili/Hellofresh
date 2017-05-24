package main_test

import (
	"encoding/json"
	"fmt"
	. "hellofresh"
	"hellofresh/model"
	"hellofresh/util"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// current ID
var currentRecipeID string = ""

var _ = Describe("Intergration test using mongodb as data storage", func() {
	It("should return alive when checking alive", func() {
		req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			Expect(res.StatusCode).To(Equal(200))
		}
	})

	It("should create recipe correctly when auth passed", func() {
		testRecipe := model.Recipe{
			Name:       "Test",
			Prep:       time.Now(),
			Difficulty: 1,
			Vegetarian: true,
		}

		expected := model.Recipe{}
		params, _ := json.Marshal(testRecipe)
		paramstr := string(params)

		req, _ := http.NewRequest("POST", "http://localhost:8080/recipes", strings.NewReader(paramstr))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		// hardcode basic auth hellofresh / hellofresh
		req.Header.Set("Authorization", "Basic aGVsbG9mcmVzaDpoZWxsb2ZyZXNo")
		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 201 - created
			Expect(res.StatusCode).To(Equal(201))

			bodyBytes, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(bodyBytes, &expected)
			Expect(expected.Name).To(Equal("Test"))
			Expect(int(expected.Difficulty)).To(Equal(1))
			Expect(expected.Vegetarian).To(Equal(true))
		}
	})

	It("should return 401 when create recipe if auth not passed", func() {
		testRecipe := model.Recipe{
			Name:       "Test",
			Prep:       time.Now(),
			Difficulty: 1,
			Vegetarian: true,
		}

		params, _ := json.Marshal(testRecipe)
		paramstr := string(params)

		// no auth in header
		req, _ := http.NewRequest("POST", "http://localhost:8080/recipes", strings.NewReader(paramstr))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 401
			Expect(res.StatusCode).To(Equal(401))
		}
	})

	It("should be able to get recipes", func() {
		req, _ := http.NewRequest("GET", "http://localhost:8080/recipes/0/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)
		expected := []model.Recipe{}

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 200
			Expect(res.StatusCode).To(Equal(200))
			bodyBytes, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(bodyBytes, &expected)
			currentRecipeID = expected[0].ID.(string)

			Expect(expected).To(HaveLen(1))
		}
	})

	It("should be able to get single recipe by id", func() {
		req, _ := http.NewRequest("GET", "http://localhost:8080/recipes/"+currentRecipeID, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)
		expected := model.Recipe{}

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 200
			Expect(res.StatusCode).To(Equal(200))
			bodyBytes, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(bodyBytes, &expected)

			Expect(fmt.Sprintf("%s", expected.ID)).To(Equal(currentRecipeID))
		}
	})

	It("should be able to update single recipe by id", func() {
		testRecipe := model.Recipe{
			Name:       "Test_Updated",
			Prep:       time.Now(),
			Difficulty: 1,
			Vegetarian: true,
		}

		params, _ := json.Marshal(testRecipe)
		paramstr := string(params)
		expected := model.Recipe{}

		req, _ := http.NewRequest("PUT", "http://localhost:8080/recipes/"+currentRecipeID, strings.NewReader(paramstr))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Basic aGVsbG9mcmVzaDpoZWxsb2ZyZXNo")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 200
			Expect(res.StatusCode).To(Equal(200))

			bodyBytes, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(bodyBytes, &expected)
			Expect(expected.Name).To(Equal("Test_Updated"))
			Expect(int(expected.Difficulty)).To(Equal(1))
			Expect(expected.Vegetarian).To(Equal(true))
		}
	})

	It("should return 401 when update single recipe if auth not passed", func() {
		testRecipe := model.Recipe{
			Name:       "Test_Updated",
			Prep:       time.Now(),
			Difficulty: 1,
			Vegetarian: true,
		}

		params, _ := json.Marshal(testRecipe)
		paramstr := string(params)

		req, _ := http.NewRequest("PUT", "http://localhost:8080/recipes/"+currentRecipeID, strings.NewReader(paramstr))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 401
			Expect(res.StatusCode).To(Equal(401))
		}
	})

	It("should be able to rate if auth passed", func() {
		// /recipes/{id}/rate/{rate:[1-5]}
		req, _ := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/recipes/%s/rate/%s", currentRecipeID, "5"), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Basic aGVsbG9mcmVzaDpoZWxsb2ZyZXNo")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 200
			Expect(res.StatusCode).To(Equal(200))
		}
	})

	It("should not be able to rate if auth not passed", func() {
		// /recipes/{id}/rate/{rate:[1-5]}
		req, _ := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/recipes/%s/rate/%s", currentRecipeID, "5"), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 401
			Expect(res.StatusCode).To(Equal(401))
		}
	})

	It("should be able to search recipes", func() {
		req, _ := http.NewRequest("GET", "http://localhost:8080/recipes/search/est", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		expected := []model.Recipe{}

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 200
			Expect(res.StatusCode).To(Equal(200))
			bodyBytes, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(bodyBytes, &expected)
			Expect(len(expected)).Should(BeNumerically(">=", 1))
		}
	})

	It("should not be able to delete recipe if auth not passed", func() {
		req, _ := http.NewRequest("DELETE", "http://localhost:8080/recipes/"+currentRecipeID, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 401
			Expect(res.StatusCode).To(Equal(401))
		}
	})

	It("should be able to delete recipe if auth passed", func() {
		req, _ := http.NewRequest("DELETE", "http://localhost:8080/recipes/"+currentRecipeID, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Basic aGVsbG9mcmVzaDpoZWxsb2ZyZXNo")

		client := &http.Client{Timeout: time.Duration(2 * time.Second)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			// Assert server response is 200
			Expect(res.StatusCode).To(Equal(200))
		}
	})
})

var _ = Describe("Restful Accessor Test", func() {
	It("should generate postgres accessor if client is postgres", func() {
		client := "postgres"
		accessor, _ := model.GetAccessor(client)
		Expect(accessor.Description()).To(Equal("postgres restful accessor"))
	})

	It("should generate mongodb accessor if client is mongodb", func() {
		client := "mongodb"
		accessor, _ := model.GetAccessor(client)
		Expect(accessor.Description()).To(Equal("mongodb restful accessor"))
	})
})

func clearCollection(app *App) {
	collection := app.DB.(*mgo.Database).C("recipe")
	if _, err := collection.RemoveAll(nil); err != nil {
		util.PanicOnError(err)
	}
}
