package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

//HashSaltPassword hash user password
func HashSaltPassword(password []byte) string {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	passHash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(passHash)
}

//VerifyHash verify hashed password and login password
func VerifyHash(hashedPassword []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
