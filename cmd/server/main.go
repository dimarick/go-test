package main

import (
	"context"
	appHttp "entities/pkg/http"
	"entities/pkg/publisher"
	"entities/pkg/subscriber"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logfile, err := os.OpenFile(`error.log`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Panic(err)
	}

	defer logfile.Close()

	log.SetOutput(io.MultiWriter(os.Stderr, logfile))

	publisherServer := &http.Server{Addr: `:8081`}
	publisher.Setup(publisherServer)
	publisherServer.Handler = http.HandlerFunc(appHttp.SetErrorHandler(publisherServer.Handler, appHttp.DefaultErrorHandler))

	go func() {
		err = publisherServer.ListenAndServe()

		if err != nil {
			log.Println(err)
		}
	}()

	subscriberServer := &http.Server{Addr: `:8082`}
	subscriber.Setup(subscriberServer)
	subscriberServer.Handler = http.HandlerFunc(appHttp.SetErrorHandler(subscriberServer.Handler, appHttp.DefaultErrorHandler))

	go func() {
		err = subscriberServer.ListenAndServe()

		if err != nil {
			log.Println(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)

	sig := <-sigChan

	switch sig {
	case syscall.SIGINT:
		log.Println(`SIGINT received`)
	case syscall.SIGHUP:
		log.Println(`SIGHUP received`)
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	publisherServer.Shutdown(ctx)
	subscriberServer.Shutdown(ctx)

	time.Sleep(1 * time.Second)
}
