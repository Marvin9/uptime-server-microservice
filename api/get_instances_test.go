package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"

	"github.com/Marvin9/uptime-server-microservice/api/setup"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/Marvin9/uptime-server-microservice/test"
)

const getInstancesAPI = "/api/instances"

type instancesResponseData struct {
	Error bool               `json:"error"`
	Data  []models.Instances `json:"data"`
}

func TestGetInstance(t *testing.T) {
	test.FakeDB(test.CREATE)
	defer test.FakeDB(test.DROP)

	db, err := test.RetryConnection()
	if err != nil {
		t.Errorf("Error connecting database.\n%v", err)
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

	instanceBodies := [][]byte{
		[]byte(`{ "url": "https://www.flipkart.com", "duration": 3600000000000 }`),
		[]byte(`{ "url": "https://www.google.com", "duration":  3600000000000 }`),
	}

	for _, instanceBody := range instanceBodies {
		addInstance := test.SimulationData{
			Method:             "POST",
			API:                addInstanceAPI,
			Body:               instanceBody,
			ExpectedStatusCode: http.StatusOK,
			ErrorMessage:       "Error adding new instance.",
			Cookie:             cookie,
		}
		test.SimulateAPI(t, router, addInstance)
	}

	getInstances := test.SimulationData{
		Method:             "GET",
		API:                getInstancesAPI,
		ExpectedStatusCode: http.StatusOK,
		ErrorMessage:       "Error getting instances.",
		Cookie:             cookie,
	}
	response := test.SimulateAPI(t, router, getInstances)

	var storeResponseData instancesResponseData
	buf, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal([]byte(buf), &storeResponseData)

	if len(storeResponseData.Data) != 2 {
		t.Errorf("%v instances were added but was stored %v in database.", 2, len(storeResponseData.Data))
	}
}
