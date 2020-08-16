package main

import (
	"log"
	"os"

	"github.com/Marvin9/uptime-server-microservice/api/setup"
	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/scheduler"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

func usePort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}
	return ":8000"
}

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting database.\n%v", err)
	}
	database.InitialMigration(db)
	db.Close()

	go scheduler.RestartAllSchedulers()

	r := setup.Router()

	r.Run(usePort())
}
