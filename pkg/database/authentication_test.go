package database_test

import (
	"net/http"
	"testing"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/Marvin9/uptime-server-microservice/test"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

func TestRegisterUser(t *testing.T) {
	test.FakeDB(test.CREATE)
	defer test.FakeDB(test.DROP)
	email := "test"
	password := "testt"
	db, err := database.ConnectDB()
	if err != nil {
		t.Errorf("Error connecting database.\n%v", err)
	}
	defer db.Close()
	db.AutoMigrate(&models.Users{})

	_, err = database.RegisterUser(email, password)
	if err != nil {
		t.Errorf("Error found while registering user.\n%v", err)
	}

	var registeredUser models.Users
	notRegistered := db.Where("email = ?", email).First(&registeredUser).RecordNotFound()
	if notRegistered {
		t.Errorf("Function returned success but user was not registered in database.")
	}

	if registeredUser.Password == password {
		t.Errorf("Password was stored in database without encryption.")
	}

	statusCode, _ := database.RegisterUser(email, password)
	if statusCode != http.StatusConflict {
		t.Errorf("Email was already registered, expected %v as status code, got %v", http.StatusConflict, statusCode)
	}
}
