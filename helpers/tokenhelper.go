package helpers

import (
	"golang-management-restaurant/database"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type SigninDetail struct {
	Email     string
	FirstName string
	LastName  string
	UserId    string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GenerateAlltoken(email string, firstname string, lastname string, uid string) {
	claims := &SigninDetail{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		UserId:    uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SigninDetail{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
}

func UpdateAlltoken() {

}

func ValidateAlltoken() {}
