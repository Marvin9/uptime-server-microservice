package database_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/Marvin9/uptime-server-microservice/test"
)

func TestInstanceCRUD(t *testing.T) {
	test.FakeDB(test.CREATE)
	defer test.FakeDB(test.DROP)

	userUniqueID := "abc"
	url := "https://www.google.com"
	duration := time.Hour

	db, err := test.RetryConnection()
	if err != nil {
		t.Errorf("Error connecting database\n.%v", err)
	}
	db.AutoMigrate(&models.Instances{}, &models.Reports{})
	defer db.Close()

	newInstanceID, err := database.CreateInstance(userUniqueID, url, duration)
	if err != nil {
		t.Errorf("Error creating new instance\n.%v", err)
	}

	var newInstanceData models.Instances
	notFound := db.Where("unique_id = ?", newInstanceID).First(&newInstanceData).RecordNotFound()
	if notFound {
		t.Errorf("Instance was successfully created but was not stored in database")
	}

	if newInstanceData.Owner != userUniqueID {
		t.Errorf("Expected new instance owner as %v, got %v", newInstanceData.Owner, userUniqueID)
	}

	if newInstanceData.URL != url {
		t.Errorf("Expected new instance url as %v, got %v", newInstanceData.URL, url)
	}

	anotherURL := "https://www.facebook.com"
	anotherDuration := time.Hour * 2

	anotherNewInstanceID, err := database.CreateInstance(userUniqueID, anotherURL, anotherDuration)
	if err != nil {
		t.Errorf("Error creating new instance.\n%v", err)
	}

	statusCode, allInstances, err := database.GetInstances(userUniqueID)
	if statusCode != http.StatusOK {
		t.Errorf("Error getting instances.\n%v", err)
	}

	expectedAllInstances := []models.Instances{
		models.Instances{
			UniqueID: anotherNewInstanceID,
			Owner:    userUniqueID,
			URL:      anotherURL,
			Duration: anotherDuration,
		},
		models.Instances{
			UniqueID: newInstanceID,
			Owner:    userUniqueID,
			URL:      url,
			Duration: duration,
		},
	}

	if len(allInstances) != len(expectedAllInstances) {
		t.Errorf("Two instance was created, but got %v from database.", len(allInstances))
	}

	if allInstances[0].UniqueID != expectedAllInstances[0].UniqueID {
		t.Errorf("%v is latest instance, but got %v from database", expectedAllInstances[0], allInstances[0])
	}

	removeInstanceErr := database.RemoveInstance(newInstanceID)
	if removeInstanceErr != nil {
		t.Errorf("Error while removing instance.\n%v", removeInstanceErr)
	}

	_, allInstances, err = database.GetInstances(userUniqueID)
	if err != nil {
		t.Errorf("Error while getting instances.\n%v", err)
	}

	if len(allInstances) != 1 {
		t.Errorf("Removing one instance, expected %v instances, but got %v", 1, len(allInstances))
	}

	if allInstances[0].UniqueID == newInstanceID {
		t.Errorf("Instance %v was expected removed, but got the value from database", allInstances[0])
	}
}
