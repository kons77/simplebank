package util

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates bcrypt hash of the password
func HashPassword(pswd string) (string, error) {

	cost := 12 // default for production
	if os.Getenv("TEST_ENV") == "true" {
		cost = bcrypt.MinCost // much faster for tests
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pswd), cost)
	if err != nil {
		return "", err
	}

	/* base64 is better because it safely converts binary data into a compact,
	printable string suitable for storage or transmission.
	hashedPassword1 := base64.StdEncoding.EncodeToString(hashedPassword)*/

	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	/* decodedHash, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return err
	}	*/
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
