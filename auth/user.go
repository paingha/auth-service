package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"

	//Need the response and request handler types in the controllers
	_ "github.com/gorilla/mux"
	response "github.com/paingha/auth-service/json"
	mailer "github.com/paingha/auth-service/mailer"
	randomstring "github.com/paingha/auth-service/randomstring"
	"github.com/paingha/auth-service/utils"
)

//User struct
type User struct {
	ID            uint       `gorm:"primary_key" json:"id"`
	FirstName     string     `gorm:"not null" json:"firstName"`
	LastName      string     `gorm:"not null" json:"lastName"`
	Email         string     `gorm:"unique;not null" json:"email"`
	Password      string     `gorm:"not null" json:"-"`
	IsAdmin       bool       `gorm:"default:false" json:"isAdmin"`
	EmailVerified bool       `gorm:"default:false" json:"emailVerified"`
	VerifyToken   string     `json:"verifyToken"`
	ContactPhone  string     `json:"phone"`
	Gender        string     `json:"gender"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

//Claims struct for Jwt
type Claims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

//Register User register controller function
func Register(w http.ResponseWriter, r *http.Request, DB *gorm.DB) {
	var newUser User
	var currentUser User
	createdResponse := response.JsonResponse("Account created successfully, Check your email and verify your email address", 200)
	existsResponse := response.JsonResponse("Account already exists. Please register", 501)
	errorResponse := response.JsonResponse("An error occured", 500)
	reqBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newUser)
	DB.Where("email = ?", newUser.Email).First(&currentUser)
	if currentUser.Email == "" {
		if err != nil {
			json.NewEncoder(w).Encode(errorResponse)
		}
		password := []byte(newUser.Password)
		newUser.Password = utils.HashSaltPassword(password)
		//Sanitize these values so users cannot get higher permissions by setting json values
		newUser.IsAdmin = false
		newUser.EmailVerified = false
		//Force FirstName and Last Name to have first character to be uppercase
		newUser.FirstName = utils.UppercaseName(newUser.FirstName)
		newUser.LastName = utils.UppercaseName(newUser.LastName)
		//Set VerifyToken for new User
		newUser.VerifyToken = randomstring.GenerateRandomString(30)
		//send token as query params in email in a link
		parentURL := "localhost:8080/api/v1/user/verify-email?token="
		verifyURL := parentURL + newUser.VerifyToken
		template1 := "<html><body><h1>Welcome to paingha.me</h1><br />" + "<a href='" + verifyURL + "'>Verify Email</a></body></html>"
		fmt.Println(template1)
		DB.Create(&newUser)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		to := []string{newUser.Email}
		mailerResponse := mailer.SendMail(to, "Welcome! Please Verify your Email", template1)
		fmt.Println(mailerResponse)
		json.NewEncoder(w).Encode(createdResponse)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(existsResponse)
	}

}

//Login User login controller function
func Login(w http.ResponseWriter, r *http.Request, DB *gorm.DB) {
	jwtSecretByte := []byte(os.Getenv("JWT_SECRET"))
	expiresAt := time.Now().Add(30 * time.Minute)
	var currentUser, dbUser User
	errorResponse := response.JsonResponse("An error occured", 500)
	invalidResponse := response.JsonResponse("Invalid credentials", 501)
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
	}
	json.Unmarshal(requestBody, &currentUser)
	passwordByte := []byte(currentUser.Password)
	DB.Where("email = ?", currentUser.Email).First(&dbUser)
	if dbUser.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(invalidResponse)
	} else {
		dbPasswordHash := []byte(dbUser.Password)
		verifyPassword := utils.VerifyHash(dbPasswordHash, passwordByte)
		if !verifyPassword {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(invalidResponse)
		} else {
			if dbUser.EmailVerified != true {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				verifyEmailResponse := response.JsonResponse("Verify your Email to Login", 500)
				json.NewEncoder(w).Encode(verifyEmailResponse)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				//////
				claims := &Claims{
					ID: dbUser.ID, //change to have the user id in here
					StandardClaims: jwt.StandardClaims{
						// In JWT, the expiry time is expressed as unix milliseconds
						ExpiresAt: expiresAt.Unix(),
					},
				}
				// Declare the token with the algorithm used for signing, and the claims
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				// Create the JWT string
				tokenString, err := token.SignedString(jwtSecretByte)
				if err != nil {
					panic(err)
				}
				loginResponse := map[string]interface{}{
					"token": tokenString,
					"user":  dbUser,
				}
				json.NewEncoder(w).Encode(loginResponse)
			}
		}
	}
}

//VerifyAuthToken User JWT verification controller function
func VerifyAuthToken(w http.ResponseWriter, r *http.Request) {
	//Endpoint to verify user JWT
	//Useful for other services
	token := r.Header.Get("Authorization")
	result, _ := utils.VerifyJWT(token)
	verifyAuthTokenResp := response.JsonResponse("Token is valid", 200)
	invalidAuthTokenResp := response.JsonResponse("Token is invalid or expired", 500)
	if result {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(verifyAuthTokenResp)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(invalidAuthTokenResp)
	}
}

//ResendVerifyEmail resends the verify email to user. Post Params: Email
func ResendVerifyEmail(w http.ResponseWriter, r *http.Request, DB *gorm.DB) {
	var currentUser User
	var dbUser User
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse := response.JsonResponse("An error occured", 500)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
	}
	json.Unmarshal(requestBody, &currentUser)
	DB.Where("Email = ?", currentUser.Email).First(&dbUser)
	fmt.Println(dbUser.VerifyToken)
	if dbUser.VerifyToken != "" {
		verifyTokenResponse := response.JsonResponse("Check your email for verification token", 200)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		parentURL := "localhost:8080/api/v1/user/verify-email?token="
		verifyURL := parentURL + dbUser.VerifyToken
		template1 := "<html><body><h1>Welcome to paingha.me</h1><br />" + "<a href='" + verifyURL + "'>Verify Email</a></body></html>"
		to := []string{dbUser.Email}
		mailerResponse := mailer.SendMail(to, "Please Verify your Email", template1)
		fmt.Println(mailerResponse)
		json.NewEncoder(w).Encode(verifyTokenResponse)
	} else {
		verifyResponse := response.JsonResponse("Email Already Verified. Login", 200)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(verifyResponse)
	}
}

//VerifyEmail User email verification controller
func VerifyEmail(w http.ResponseWriter, r *http.Request, DB *gorm.DB) {
	//Endpoint to verify user email
	//create new route that pulls the query params and checks if it is the same as what has been saved in the user table
	//if it is the same then update the user row and change verified to true
	//if not the same ask them to request new verification token to email
	token := r.URL.Query().Get("token")
	var dbUser User
	//find the user
	DB.Where("verify_token = ?", token).First(&dbUser)
	//Check if the user has already verified their account to save a db call if they have
	if dbUser != (User{}) {
		if dbUser.EmailVerified != true {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			DB.Model(&dbUser).Where("Verify_Token = ?", token).Updates(map[string]interface{}{"EmailVerified": true, "VerifyToken": ""})
			updateResponse := response.JsonResponse("Email Verified Successfully", 200)
			json.NewEncoder(w).Encode(updateResponse)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			alreadyResponse := response.JsonResponse("Email Already Verified", 409)
			json.NewEncoder(w).Encode(alreadyResponse)
		}
	} else {
		errorResponse := response.JsonResponse("An error occured while verifying your email", 500)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
	}

}
func isAdmin(token string, DB *gorm.DB) bool {
	var currentUser User
	_, ID := utils.VerifyJWT(token)
	DB.Where("ID = ?", ID).First(&currentUser)
	return currentUser.IsAdmin
}

//GetByID get user by the user_id
func GetByID(w http.ResponseWriter, r *http.Request, DB *gorm.DB) {
	var dbUser User
	params := mux.Vars(r)
	userID := params["id"]
	//Need to make sure that the user that is requesting user info is either the user or an admin user
	token := r.Header.Get("Authorization")
	result, ID := utils.VerifyJWT(token)
	myID := strconv.FormatUint(uint64(ID), 10)
	//results := utils.IsAdmin(token, DB)
	//fmt.Printf("%v", results)
	if (result && userID == myID) || isAdmin(token, DB) {
		DB.Where("ID = ?", userID).First(&dbUser)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dbUser)
	} else {
		notauthorizedResponse := response.JsonResponse("Not authorized", 409)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(notauthorizedResponse)
	}

}
