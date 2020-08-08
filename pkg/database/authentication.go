package database

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"

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

	uniqueID, err := utils.GenerateUniqueID()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Check if email already exist
	emailNotExist := db.Where("Email = ?", email).First(&models.Users{}).RecordNotFound()
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

// LoginUser returns http status code, jwt token, error
func LoginUser(email, password string) (int, string, error) {
	db, err := ConnectDB()
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	defer db.Close()

	// Check if email exist
	var detailsOfUser models.Users
	user := db.Where("Email = ?", email).First(&detailsOfUser)
	if user.RecordNotFound() {
		return http.StatusNotFound, "", errors.New("Email do not exist, please register first")
	}

	// Match password
	originalPasswordOfThatEmail := detailsOfUser.Password
	isPasswordMatched := utils.ComparePassword(originalPasswordOfThatEmail, password)
	if !isPasswordMatched {
		return http.StatusUnauthorized, "", errors.New("Wrong password")
	}

	// Generate token
	jwtClaims := &models.Claims{
		UniqueID: detailsOfUser.UniqueID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: utils.JWTExpirationTime,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtTokenString, err := jwtToken.SignedString(utils.GetJWTKey())
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, jwtTokenString, nil
}
