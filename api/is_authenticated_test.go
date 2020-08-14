package api_test

import (
	"net/http"
	"testing"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"

	"github.com/Marvin9/uptime-server-microservice/api/setup"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/test"
)

const isAuthenticatedAPI = "/auth/ping"

func TestIsAuthenticatedAPI(t *testing.T) {
	test.FakeDB(test.CREATE)
	defer test.FakeDB(test.DROP)

	db, err := database.ConnectDB()
	if err != nil {
		t.Errorf("Error connecting database.\n%v", err)
	}
	db.AutoMigrate(&models.Users{})
	defer db.Close()

	router := setup.Router()

	pingAuth := test.SimulationData{
		Method:             "GET",
		API:                isAuthenticatedAPI,
		ExpectedStatusCode: http.StatusUnauthorized,
		ErrorMessage:       "Must be unauthorized without login.",
	}

	test.SimulateAPI(t, router, pingAuth)

	jwtToken, _ := generateLogInCookie("abc@gmail.com", "abc")
	cookie := http.Cookie{
		Name:  middlewares.JWTCookieName,
		Value: jwtToken,
	}

	pingAuth.Cookie = cookie
	pingAuth.ExpectedStatusCode = http.StatusOK
	pingAuth.ErrorMessage = "User is authorized with cookie."

	test.SimulateAPI(t, router, pingAuth)
}
