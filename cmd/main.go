package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	todoAPI "github.com/howieter/rest_api_todoL/internal"
	"github.com/sirupsen/logrus"
)

func main() {
	err := todoAPI.CreateDB()
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("DB created successfuly")
	}

	db, err := todoAPI.ConnectDB()
	if err != nil {
		log.Fatalf("error while connecting to DB: %v", err)
	}
	log.Println("DB connected successfuly")
	defer db.Close(context.Background())

	err = todoAPI.MakeQuery(db, "db/0001_create_table_tasks.up.sql")
	if err != nil {
		log.Printf("%v", err)
	}
	err = todoAPI.MakeQuery(db, "db/0002_insert_fields.up.sql")
	if err != nil {
		log.Printf("%v", err)
	}

	webApp := fiber.New()
	todoAPI.SetupRoutes(webApp, db)

	logrus.Fatal(webApp.Listen(":80"))
}
