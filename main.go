package main

import (
	"github.com/Marvin9/uptime-server-microservice/api/setup"
	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/scheduler"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, _ := database.ConnectDB()
	database.InitialMigration(db)
	db.Close()

	go scheduler.RestartAllSchedulers()

	r := setup.Router()

	r.Run(":8000")
}
