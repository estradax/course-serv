package main

import (
	"log"

	"github.com/estradax/course-serv/internal"
	"github.com/estradax/course-serv/internal/model"
)

func main() {
	err := internal.LoadEnv()
	if err != nil {
		log.Fatalln("Cannot loadEnv: ", err.Error())
	}

	db, err := internal.ConnectDB()
	if err != nil {
		log.Fatalln("Something went wrong: ", err.Error())
	}

	err = db.AutoMigrate(&model.User{}, &model.Course{})
	if err != nil {
		log.Fatalln("Cannot migrate: ", err.Error())
	}
}
