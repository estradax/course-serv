package main

import (
	"log"

	"github.com/estradax/course-serv/internal"
	"github.com/estradax/course-serv/internal/model"
)

func main() {
	db, err := internal.ConnectDB()
	if err != nil {
		log.Fatalln("Something went wrong: ", err.Error())
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalln("Cannot migrate: ", err.Error())
	}
}