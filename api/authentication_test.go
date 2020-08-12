package api_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

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

func setRequestHeaderJSON(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func responseError(t *testing.T, message string, expected, got int, body *bytes.Buffer) {
	t.Errorf("%v\n\nExpected status code %v, got %v\nBody: %v", message, expected, got, body)
}

func simulateAPI(t *testing.T, router *gin.Engine, s simulationData) {
	w := httptest.NewRecorder()
	jsonBody := bytes.NewBuffer(s.body)
	req, _ := http.NewRequest(s.method, s.api, jsonBody)
	setRequestHeaderJSON(req)
	router.ServeHTTP(w, req)

	if w.Code != s.expectedStatusCode {
		responseError(t, s.errorMessage, s.expectedStatusCode, w.Code, w.Body)
	}
}

type simulationData struct {
	method             string
	api                string
	body               []byte
	expectedStatusCode int
	errorMessage       string
}

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
		responseError(t, fmt.Sprintf("%v should not accept empty body.", registerAPI), http.StatusBadRequest, w.Code, w.Body)
	}

	jsonVal := []byte(`{ "email": "mayursiinh@gmail.com", "password": "abc" }`)
	wrongJsonValPassword := []byte(`{ "email": "mayursiinh@gmail.com", "password": "abcd" }`)
	var simulation = []simulationData{
		simulationData{
			method:             "POST",
			api:                registerAPI,
			body:               jsonVal,
			expectedStatusCode: http.StatusOK,
			errorMessage:       "Error while registering user.",
		},
		simulationData{
			method:             "POST",
			api:                registerAPI,
			body:               jsonVal,
			expectedStatusCode: http.StatusConflict,
			errorMessage:       "Already registered user.",
		},
		simulationData{
			method:             "POST",
			api:                loginAPI,
			body:               jsonVal,
			expectedStatusCode: http.StatusOK,
			errorMessage:       "Error logging in registered user.",
		},
		simulationData{
			method:             "POST",
			api:                loginAPI,
			body:               wrongJsonValPassword,
			expectedStatusCode: http.StatusUnauthorized,
			errorMessage:       "Password was wrong. It's response must be unauthorized.",
		},
		simulationData{
			method:             "POST",
			api:                loginAPI,
			body:               jsonVal,
			expectedStatusCode: http.StatusOK,
			errorMessage:       "With correct credentials, It must provide 200 status",
		},
	}

	for _, sim := range simulation {
		simulateAPI(t, router, sim)
	}
}
