package models

import (
	"github.com/dgrijalva/jwt-go"
)

// Claims is jwt claim
type Claims struct {
	UniqueID string `json:"uniqueId"`
	jwt.StandardClaims
}
