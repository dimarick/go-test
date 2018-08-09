package main

import (
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
}
