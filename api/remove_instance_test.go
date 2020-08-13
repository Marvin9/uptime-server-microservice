package api_test

import (
	"net/http"
	"testing"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/api/setup"
	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/Marvin9/uptime-server-microservice/pkg/utils"
	"github.com/dgrijalva/jwt-go"

	"github.com/Marvin9/uptime-server-microservice/test"
)

func TestRemoveInstanceAPI(t *testing.T) {
	test.FakeDB(test.CREATE)
	defer test.FakeDB(test.DROP)

	db, err := database.ConnectDB()
	if err != nil {
		t.Errorf("Error connecting database.\n%v", err)
	}
	defer db.Close()
	db.AutoMigrate(&models.Users{}, &models.Instances{}, &models.Reports{})

	router := setup.Router()
	jwtToken, err := generateLogInCookie("m@gmail.com", "abc")
	jwtClaims := &models.Claims{}
	jwt.ParseWithClaims(jwtToken, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTKey(), nil
	})
	if err != nil {
		t.Errorf("Error generating login cookie.")
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

	var instances []models.Instances
	db.Where("owner = ?", jwtClaims.UniqueID).Find(&instances)

	if len(instances) != 1 {
		t.Errorf("Instance was not added in database.\n")
	}

	removeInstance := addInstance
	removeInstance.Method = "DELETE"
	removeInstance.Body = []byte(`{ "instance_id" : "` + instances[0].UniqueID + `" }`)
	removeInstance.ErrorMessage = "Error removing instance."
	test.SimulateAPI(t, router, removeInstance)

	prevInstance := instances[0]
	instances = []models.Instances{}
	db.Where("owner = ?", jwtClaims.UniqueID).Find(&instances)

	if len(instances) != 0 {
		t.Errorf("Instance was not removed from database.\n")
	}

	if database.IsInstanceRunning(prevInstance.Owner, prevInstance.URL) {
		t.Errorf("Instance was removed from database but still running in memory.\n")
	}
}
