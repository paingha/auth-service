package utils

import (
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

//Claims struct for Jwt
type Claims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

//VerifyJWT takes in token as a string and returns a boolean.
func VerifyJWT(jwtToken string) (bool, uint64) {
	var response = false
	var emptyString uint64 = 0
	// Initialize a new instance of `Claims`
	//fmt.Println(jwtToken)
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			//w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Token Error")
			response = false
			return response, emptyString
		}
		fmt.Println("Bad Request")
		response = false
		return response, emptyString
	}
	if !tkn.Valid {
		//w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Invalid Token")
		response = false
		return response, emptyString
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	response = true
	return response, uint64(claims.ID)
}

