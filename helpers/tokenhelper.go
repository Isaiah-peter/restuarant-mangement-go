package helpers

import (
	"context"
	"fmt"
	"golang-management-restaurant/database"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GenerateAlltoken(email string, firstname string, lastname string, uid string) (signinToken string, refreshTokens string, err error) {
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

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAlltoken(signinToken, refreshToken, UUid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"token", signinToken})
	updateObj = append(updateObj, bson.E{"refresh_token", refreshToken})

	UpdateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"update_at", UpdateAt})
	upsert := true
	filter := bson.M{"user_id": UUid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	return
}

func ValidateAlltoken(signinToken string) (claims *SigninDetail, msg string) {
	token, err := jwt.ParseWithClaims(
		signinToken,
		&SigninDetail{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	// token is invalid
	claim, ok := token.Claims.(*SigninDetail)
	if !ok {
		msg = fmt.Sprintf("the toke is not vaild")
		msg = err.Error()
	}
	//token is expired
	if claim.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token exprired")
	}
	return claims, msg
}
