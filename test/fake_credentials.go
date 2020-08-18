package test

import "fmt"

// Credential type
type Credential struct {
	Email    string
	Password string
}

var start = 0

const password = "abcd"

/*GenerateFakeCredentials will generate unique fake credentials to use in tests.
We have dropping database after each test case. But in travis build it is either not dropping or running parallely database conflicts.
*/
func GenerateFakeCredentials() Credential {
	start++
	email := fmt.Sprintf("%v@gmail.com", start)
	return Credential{
		email, password,
	}
}
