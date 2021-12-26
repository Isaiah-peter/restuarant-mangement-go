package controllers

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	helpers "golang-restaurant-management/helpers"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
		projectStage := bson.D{{
			"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

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
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*100)
		var user models.User

		var userId = c.Param("user_id")
		err := foodCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while getting this user"})
		}

		c.JSON(http.StatusOK, user)

	}
}

func SingUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		var user models.User

		//convert the JSON data coming from postman to something golang can understand
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "error while trying to convert data to golang understandable format"})
			return
		}
		//validate the data base on user struct
		validateErr := validate.Struct(&user)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while validating" + validateErr.Error()})
			return
		}

		// you'll check if the email hav been use by other user
		count, err := usercollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusConflict, gin.H{"error": "email already used" + " " + *user.Email + " " + *user.FirstName + " " + *user.LastName})
			return
		}
		// hash password
		password := HashPassword(*user.Password)
		user.Password = &password
		//you'll also check if the phone no have being already use by another user
		count, err = usercollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusConflict, gin.H{"error": "phone number already used"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "your email or password are alreay used"})
			return
		}
		//create some extra detail for the user object - created-at, updated-at
		user.CreatedAt = timetouse
		user.UpdatedAt = timetouse
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		//generate token and refreshtoken (generat all token function form helper)
		token, refreshToken, _ := helpers.GenerateAlltoken(*user.Email, *user.Email, *user.FirstName, *user.LastName, *user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken
		//if all ok, then insert the new user into the user collection
		reusltInsertNumber, insertErr := usercollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
		}

		//return status ok, and send the rusult back

		c.JSON(http.StatusOK, reusltInsertNumber)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		var user models.User
		var foundUser models.User
		//convert the login data from the postman which is in JSON to golang readable format
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while converting to golang readable format" + err.Error()})
			return
		}
		//find a user with the email and see if the user exist

		err := usercollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email does not exist"})
			return
		}
		//then you verify the password

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if passwordIsValid != true {
			c.JSON(http.StatusConflict, gin.H{"error": msg})
			return
		}
		//if all goes well the you generate tokens
		token, refreshToken, _ := helpers.GenerateAlltoken(*user.Email, *user.Email, *user.FirstName, *user.LastName, *user.UserId)

		//update tokens - token and refresh token
		helpers.UpdateAlltoken(token, refreshToken, foundUser.UserId)

		//return status ok and send the result back
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	result, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(result)
}

func VerifyPassword(userPasword string, providerPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providerPassword), []byte(userPasword))
	check := true
	msg := ""

	if err != nil {
		check = false
		msg = fmt.Sprintf("password is incorrect try again")
	}
	return check, msg
}
