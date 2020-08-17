package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/api/setup"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/test"
)

const getReportsAPI = "/api/report"

type reportsResponseData struct {
	Error bool             `json:"error"`
	Data  []models.Reports `json:"data"`
}

func TestGetReport(t *testing.T) {
	test.FakeDB(test.CREATE)
	defer test.FakeDB(test.DROP)

	db, err := database.ConnectDB()
	if err != nil {
		t.Errorf("Error connecting database.\n%v", err)
	}
	defer db.Close()
	db.AutoMigrate(&models.Users{}, &models.Instances{}, &models.Reports{})

	router := setup.Router()
	jwtToken, err := generateLogInCookie("mayur@gmail.com", "abc")
	if err != nil {
		t.Errorf("Error generating login cookie.\n%v", err)
	}
	cookie := http.Cookie{
		Name:  middlewares.JWTCookieName,
		Value: jwtToken,
	}

	addInstance := test.SimulationData{
		Method:             "POST",
		API:                addInstanceAPI,
		Body:               []byte(`{ "url": "https://www.google.com", "duration": 3600000000000 }`),
		ExpectedStatusCode: http.StatusOK,
		ErrorMessage:       "Error adding instance.",
		Cookie:             cookie,
	}
	test.SimulateAPI(t, router, addInstance)

	// wait until report is added, because that process uses goroutine
	time.Sleep(time.Millisecond * 1000)

	getReports := test.SimulationData{
		Method:             "GET",
		API:                getReportsAPI,
		ExpectedStatusCode: http.StatusOK,
		ErrorMessage:       "Error getting reports.",
		Cookie:             cookie,
	}

	response := test.SimulateAPI(t, router, getReports)

	var storeResponseData reportsResponseData
	buf, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal([]byte(buf), &storeResponseData)

	if len(storeResponseData.Data) != 1 {
		t.Errorf("Instance was added and must give one report instantly. Expected report %v, got %v", 1, len(storeResponseData.Data))
	}

	addInstance.Body = []byte(`{ "url": "https://www.yahoo.com", "duration": 3600000000000 }`)
	test.SimulateAPI(t, router, addInstance)

	time.Sleep(time.Second * 5)

	response = test.SimulateAPI(t, router, getReports)
	buf, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal([]byte(buf), &storeResponseData)

	if len(storeResponseData.Data) != 2 {
		t.Errorf("Report for 2 instances must be there because both were added. Expected reports %v, got %v", 2, len(storeResponseData.Data))
	}
}
