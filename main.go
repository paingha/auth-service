package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
	"github.com/paingha/auth-service/middleware"

	user "github.com/paingha/auth-service/auth"
	//config "github.com/paingha/auth-service/config"
	db "github.com/paingha/auth-service/db"
)

//DB global database variable
var DB *gorm.DB

func homeHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user.Login(w, r, DB)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	user.Register(w, r, DB)
}

func verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	user.VerifyAuthToken(w, r)
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	user.VerifyEmail(w, r, DB)
}

func getVerifyHandler(w http.ResponseWriter, r *http.Request) {
	user.ResendVerifyEmail(w, r, DB)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	user.GetByID(w, r, DB)
}

func main() {
	dbconnection := db.DBConnect()
	DB = dbconnection
	//Connect to database
	//Set Env variables
	//config.SetupEnvVariables()
	//Routes Defined Here
	r := mux.NewRouter()
	router := r.PathPrefix("/api/v1/user").Subrouter()
	router.Use(middleware.AuthMiddleware)
	//User Api Home Route
	router.HandleFunc("/", homeHandler).Methods("GET")
	//User API Login Route
	router.HandleFunc("/login", loginHandler).Methods("POST")
	//User API Register Route
	router.HandleFunc("/register", registerHandler).Methods("POST")
	//User API Verify Email
	router.HandleFunc("/verify-email", verifyEmailHandler).Methods("GET")
	//User Api to get Verification Token
	router.HandleFunc("/get-verify-token", getVerifyHandler).Methods("POST")
	//User API Verify JWT Route
	router.HandleFunc("/token-verify", verifyTokenHandler).Methods("GET")
	//User API Get User Route
	router.HandleFunc("/{id:[0-9]+}", getUserHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
