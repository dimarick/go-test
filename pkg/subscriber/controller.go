package subscriber

import (
	"encoding/json"
	appHttp "entities/pkg/http"
	"entities/pkg/schema"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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
		log.Println(`Shutting down subscriber server`)
		db.Close()
	})

	router := mux.NewRouter()

	router.HandleFunc(`/entities/{id}`, func(writer http.ResponseWriter, request *http.Request) {
		user := schema.User{}

		db.Find(&user, `id = ?`, mux.Vars(request)[`id`])

		if user.Id == `` {
			appHttp.ErrorResponse(writer, `Not found`, http.StatusNotFound)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(user)

		if err != nil {
			log.Panic(err, request)
		}

		writer.Write(jsonData)
	}).Methods(http.MethodGet)

	server.Handler = router

	return
}
