package database

import (
	"errors"
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"

	"github.com/Marvin9/uptime-server-microservice/pkg/utils"
)

// RegisterUser returns status code & error (if)
func RegisterUser(email, password string) (int, error) {
	hashPassword, err := utils.HashAndSalt(password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	db, err := ConnectDB()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer db.Close()

	db.AutoMigrate(&models.Users{})

	uniqueID, err := utils.GenerateUniqueID()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Check if email already exist
	emailNotExist := db.Where("Email = ?", email).First(&models.User{}).RecordNotFound()
	if !emailNotExist {
		return http.StatusConflict, errors.New("Email already exists")
	}

	db.Create(&models.Users{
		UniqueID: uniqueID,
		Email:    email,
		Password: hashPassword,
	})

	return http.StatusOK, nil
}

// func FindInTable() error {
// 	db, err := ConnectDB()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()
// 	db.AutoMigrate(&TestTable{})
// 	var values []TestTable
// 	db.Find(&values)

// 	b, _ := json.Marshal(values)
// 	fmt.Println(string(b))
// 	return nil
// }
