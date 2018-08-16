package main

import (
	"entities/pkg/queue"
	"entities/pkg/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

func main() {
	db, err := gorm.Open("postgres", os.Getenv(`DATABASE_URL`))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.AutoMigrate(schema.User{})

	conn := queue.Connect()
	channel, err := conn.Channel()
	channel.ExchangeDeclare(`main`, `direct`, true, false, false, false, nil)
	channel.QueueDeclare(`main.insert_user`, true, false, false, false, nil)
	channel.QueueBind(`main.insert_user`, `main.insert_user`, `main`, false, nil)
}
