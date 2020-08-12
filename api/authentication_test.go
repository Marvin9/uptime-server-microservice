package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Marvin9/uptime-server-microservice/api/setup"
	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/Marvin9/uptime-server-microservice/test"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	registerAPI = "/auth/register"
	loginAPI    = "/auth/login"
)

func TestAuthenticationAPI(t *testing.T) {
	test.FakeDB(test.CREATE)
	defer test.FakeDB(test.DROP)

	db, err := database.ConnectDB()
	if err != nil {
		t.Errorf("Error connecting database.\n%v\n", err)
	}
	defer db.Close()
	db.AutoMigrate(&models.Users{})

	router := setup.Router()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", registerAPI, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		test.ResponseError(t, fmt.Sprintf("%v should not accept empty body.", registerAPI), http.StatusBadRequest, w.Code, w.Body)
	}

	jsonVal := []byte(`{ "email": "mayursiinh@gmail.com", "password": "abc" }`)
	wrongJsonValPassword := []byte(`{ "email": "mayursiinh@gmail.com", "password": "abcd" }`)
	var simulation = []test.SimulationData{
		test.SimulationData{
			Method:             "POST",
			API:                registerAPI,
			Body:               jsonVal,
			ExpectedStatusCode: http.StatusOK,
			ErrorMessage:       "Error while registering user.",
		},
		test.SimulationData{
			Method:             "POST",
			API:                registerAPI,
			Body:               jsonVal,
			ExpectedStatusCode: http.StatusConflict,
			ErrorMessage:       "Already registered user.",
		},
		test.SimulationData{
			Method:             "POST",
			API:                loginAPI,
			Body:               jsonVal,
			ExpectedStatusCode: http.StatusOK,
			ErrorMessage:       "Error logging in registered user.",
		},
		test.SimulationData{
			Method:             "POST",
			API:                loginAPI,
			Body:               wrongJsonValPassword,
			ExpectedStatusCode: http.StatusUnauthorized,
			ErrorMessage:       "Password was wrong. It's response must be unauthorized.",
		},
		test.SimulationData{
			Method:             "POST",
			API:                loginAPI,
			Body:               jsonVal,
			ExpectedStatusCode: http.StatusOK,
			ErrorMessage:       "With correct credentials, It must provide 200 status",
		},
	}

	for _, sim := range simulation {
		test.SimulateAPI(t, router, sim)
	}
}
