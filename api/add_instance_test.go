package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/api/setup"
	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/Marvin9/uptime-server-microservice/test"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const addInstanceAPI = "/api/instance"

type responseData struct {
	Data string `json:"data"`
}

func TestAddInstanceAPI(t *testing.T) {
	test.FakeDB(test.CREATE)
	defer test.FakeDB(test.DROP)

	db, err := test.RetryConnection()
	if err != nil {
		t.Errorf("Error connecting database.\n%v\n", err)
	}
	defer db.Close()
	db.AutoMigrate(&models.Users{}, &models.Instances{}, &models.Reports{})

	router := setup.Router()

	credentials := test.GenerateFakeCredentials()

	jwtToken, err := generateLogInCookie(credentials.Email, credentials.Password)
	if err != nil {
		t.Errorf("Error generating login cookie.\n%v", err)
	}

	cookie := http.Cookie{
		Name:  middlewares.JWTCookieName,
		Value: jwtToken,
	}

	// 1 ^ 9 = 1 second
	instanceBody := []byte(`{ "url": "https://www.google.com", "duration": 1000000000 }`)
	addInstance := test.SimulationData{
		Method:             "POST",
		API:                addInstanceAPI,
		Body:               instanceBody,
		ExpectedStatusCode: http.StatusBadRequest,
		ErrorMessage:       "Duration of instance should be minimum 1 hour, but it accepted less than that",
		Cookie:             cookie,
	}
	test.SimulateAPI(t, router, addInstance)

	// 1 hour
	instanceBody = []byte(`{ "url": "https://www.google.com", "duration": 3600000000000 }`)
	addInstance.Body = instanceBody
	addInstance.ExpectedStatusCode = http.StatusOK
	addInstance.ErrorMessage = "Instance was not added."

	resp := test.SimulateAPI(t, router, addInstance)
	var storeResponseData responseData
	buf, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(buf), &storeResponseData)
	newInstanceID := storeResponseData.Data

	var instanceData models.Instances
	addedInstanceNotFound := db.Where("unique_id = ?", newInstanceID).First(&instanceData).RecordNotFound()
	if addedInstanceNotFound {
		t.Errorf("Instance %v was not stored in database.\n", newInstanceID)
	}

	if !database.IsInstanceRunning(instanceData.Owner, instanceData.URL) {
		t.Errorf("Instance was created but not running in memory.\n")
	}
}
