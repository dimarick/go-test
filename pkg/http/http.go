package http

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
)

func ErrorResponse(writer http.ResponseWriter, message string, httpCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpCode)

	jsonData, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: message,
	})

	writer.Write(jsonData)

	log.Println(`Http error: ` + strconv.Itoa(httpCode) + `, ` + message)
}

func SetErrorHandler(router http.Handler, errorHandler func(interface{}, http.ResponseWriter, *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	handler := func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorHandler(err, writer, request)
			}
		}()

		router.ServeHTTP(writer, request)
	}

	return handler
}

func DefaultErrorHandler(err interface{}, writer http.ResponseWriter, request *http.Request) {
	ErrorResponse(writer, `Failed to handle request. Please try again later`, http.StatusInternalServerError)
	log.Println("Unexpected error while handling request", err, string(debug.Stack()))
}
