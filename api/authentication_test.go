package api_test

import (
	"bytes"
	"encoding/json"
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
		t.Errorf("%v should not accept with empty body. Expected %v status code, got %v", registerAPI, http.StatusBadRequest, w.Code)
	}

	jsonVal, _ := json.Marshal(models.Users{
		Email:    "mayursiinh@gmail.com",
		Password: "abc",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", registerAPI, bytes.NewBuffer(jsonVal))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Error while registering user, Expected status code %v, got %v", http.StatusOK, w.Code)
	}
}
