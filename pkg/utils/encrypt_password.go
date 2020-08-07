package utils

import "golang.org/x/crypto/bcrypt"

// https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72

// HashAndSalt to hash password and store it in database
func HashAndSalt(pwd string) (string, error) {
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePassword to compare hashed password while login
func ComparePassword(hashedPassword, plainPassword string) bool {
	byteHash := []byte(hashedPassword)
	bytePlainPass := []byte(plainPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPass)
	if err != nil {
		return false
	}

	return true
}
