package publisher

import (
	"encoding/json"
	"entities/pkg/queue"
	"entities/pkg/schema"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Setup(server *http.Server) {

	dbUrl := os.Getenv(`DATABASE_URL`)
	db, err := gorm.Open("postgres", dbUrl)

	if err != nil {
		log.Panic(err)
	}

	defer server.RegisterOnShutdown(func() {
		log.Println(`Shutting down publisher server`)
		db.Close()
	})

	router := mux.NewRouter()
	amqpConn := queue.Connect()

	router.HandleFunc(`/expensive-entities`, func(writer http.ResponseWriter, request *http.Request) {
		user := schema.User{}
		body, err := ioutil.ReadAll(request.Body)

		if err != nil {
			log.Panic(err, request)
		}

		defer request.Body.Close()

		err = json.Unmarshal(body, &user)

		if err != nil {
			log.Panic(err, request)
		}

		amqpConn.Publish(`insert_user`, user)

		err = db.Create(&user).Error

		if err != nil {
			log.Panic(err, request)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(user)

		if err != nil {
			log.Panic(err, request)
		}

		writer.Write(jsonData)
	}).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	server.Handler = router

	return
}
