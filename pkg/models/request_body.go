package models

// User is for /register post body
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
