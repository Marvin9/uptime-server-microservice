package utils

import (
	"os"
	"time"
)

// GetJWTKey returns secret key to generate token
func GetJWTKey() []byte {
	// JWTKey is secret json webtoken key stored in environment variable
	var JWTKey = []byte(os.Getenv("JWT_KEY"))
	return JWTKey
}

// JWTExpireAfter is delay, after which token will expire
var JWTExpireAfter = 30 * time.Minute

// JWTCookieExpireAfter is for cookie expiration same as jwt expiration
var JWTCookieExpireAfter = 60 * 30

// JWTExpirationTime is used while generating token as well as in cookie
var JWTExpirationTime = time.Now().Add(JWTExpireAfter).Unix()
