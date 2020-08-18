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

	credentials := test.GenerateFakeCredentials()
	email := credentials.Email
	password := credentials.Password
	db, err := test.RetryConnection()
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

func TestLoginUser(t *testing.T) {
	test.FakeDB(test.CREATE)
	defer test.FakeDB(test.DROP)

	db, err := test.RetryConnection()
	if err != nil {
		t.Errorf("Error connecting database.\n%v", err)
	}
	defer db.Close()
	db.AutoMigrate(&models.Users{})

	credentials := test.GenerateFakeCredentials()
	email := credentials.Email
	password := credentials.Password
	wrongPassword := password + "-"

	statusCode, _, _ := database.LoginUser(email, password)
	if statusCode != http.StatusNotFound {
		t.Errorf("Expected unregistered email to be not found %v, but got %v", http.StatusNotFound, statusCode)
	}

	_, err = database.RegisterUser(email, password)
	if err != nil {
		t.Errorf("Error registering user.\n%v", err)
	}

	statusCode, _, err = database.LoginUser(email, password)
	if statusCode != http.StatusOK {
		t.Errorf("Could not login with correct credentials.\n%v", err)
	}

	statusCode, _, err = database.LoginUser(email, wrongPassword)
	if statusCode != http.StatusUnauthorized {
		t.Errorf("Expected status for wrong password was %v, but got %v", http.StatusUnauthorized, statusCode)
	}
}
