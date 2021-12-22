package controllers

import (
	"golang-restaurant-management/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var usercollection *mongo.Collection = database.OpenCollection(database.Client, "user")
func GetUsers() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func SingUp() gin.HandlerFunc {
	return func(context *gin.Context) {
		//convert the JSON data coming from postman to something golang can understand

		//validate the data base on user struct

		// you'll check if the email hav been use by other user

		// hash password

		//you'll also check if the phone no have being already use by another user

		//create some extra detail for the user object - created-at, updated-at

		//generatetoken and re
	}
}

func Login() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func HashPassword(password string) string {
	return ""
}

func VerifyPassword(userPasword string, providerPsword string) (bool, string) {
	return false, ""
}
