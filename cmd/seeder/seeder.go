package main

import (
	"log"

	"github.com/estradax/course-serv/internal"
	"github.com/estradax/course-serv/internal/model"
	"github.com/go-faker/faker/v4"
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

	var courses = [...]model.Course{
		{Title: "ReactJS", Description: faker.Paragraph(), Price: 50000},
		{Title: "VueJS", Description: faker.Paragraph(), Price: 50000},
		{Title: "NextJS", Description: faker.Paragraph(), Price: 50000},
		{Title: "Golang", Description: faker.Paragraph(), Price: 50000},
		{Title: "Bootstrap", Description: faker.Paragraph(), Price: 50000},
		{Title: "HTML", Description: faker.Paragraph(), Price: 50000},
		{Title: "CSS", Description: faker.Paragraph(), Price: 50000},
		{Title: "JavaScript", Description: faker.Paragraph(), Price: 50000},
	}

	for _, course := range courses {
		result := db.Create(&course)
		if result.Error != nil {
			log.Fatalln("Cannot create course: ", result.Error.Error())
		}	
	}
}
