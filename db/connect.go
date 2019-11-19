package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	user "github.com/paingha/auth-service/auth"

	//Not needed so escaped
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	dbname   = "auth-service"
	username = "postgres"
	password = "123456"
)

//DBconnection  is the database connection string used to make requests to the data
var DBconnection *gorm.DB
var err error

func connect() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	DBconnection, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer DBconnection.Close()
	fmt.Println("Successfully connected!")
	DBconnection.AutoMigrate(&user.User{})
	return DBconnection

}

//DBConnect is set to the connect function and exported
var DBConnect = connect
var DB *gorm.DB = DBconnection
