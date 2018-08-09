package main

import (
	"encoding/json"
	"entities/pkg/schema"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	dbUrl := os.Getenv(`DATABASE_URL`)
	db, err := gorm.Open("postgres", dbUrl)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc(`/expensive-entities`, func(writer http.ResponseWriter, request *http.Request) {
		user := schema.User{}
		body, err := ioutil.ReadAll(request.Body)

		defer request.Body.Close()

		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &user)

		if err != nil {
			panic(err)
		}

		err = db.Create(&user).Error

		if err != nil {
			panic(err)
		}

		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(user)

		if err != nil {
			panic(err)
		}

		writer.Write(jsonData)
	}).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(`/entities/{id}`, func(writer http.ResponseWriter, request *http.Request) {
		user := schema.User{}

		db.Find(&user, `id = ?`, mux.Vars(request)[`id`])

		if user.Id == `` {
			writer.WriteHeader(http.StatusNotFound)
		}

		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(user)

		if err != nil {
			panic(err)
		}

		writer.Write(jsonData)
	}).Methods(http.MethodGet)

	err = http.ListenAndServe(`:8081`, router)

	if err != nil {
		panic(err)
	}
}
