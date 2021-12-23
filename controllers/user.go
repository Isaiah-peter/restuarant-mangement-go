package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var usercollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10

		}
		pages, err1 := strconv.Atoi(c.Query("page"))

		if err1 != nil || pages < 1 {
			pages = 1
		}

		startIndex := (pages - 1) * recordPerPage

		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				},
			}}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage,
		})

		defer cancel()

		//either pass an error

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while getting all data "})
		}

		var allUser []bson.M

		if err = result.All(ctx, &allUser); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allUser)

		//ideally want to return all users base on the query
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SingUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//convert the JSON data coming from postman to something golang can understand

		//validate the data base on user struct

		// you'll check if the email hav been use by other user

		// hash password

		//you'll also check if the phone no have being already use by another user

		//create some extra detail for the user object - created-at, updated-at

		//generate token and refreshtoken (generat all token function form helper)

		//if all ok, then insert the new user into the user collection

		//return status ok, and send the rusult back
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		//convert the login data from the postman which is in JSON to golang readable format

		//find a user with the email and see if the user exist

		//then you verify the password

		//if all goes well the you generate tokens

		//update tokens - token and refresh token

		//return status ok and send the result back
	}
}

func HashPassword(password string) string {
	return ""
}

func VerifyPassword(userPasword string, providerPsword string) (bool, string) {
	return false, ""
}
