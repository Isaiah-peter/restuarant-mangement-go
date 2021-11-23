package controllers

import (
	"github.com/gin-gonic/gin"
)

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

	}
}

func Login() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func HashPassword(password string) string  {
	return ""
}

func VerifyPassword(userPasword string, providerPsword string) (bool, string)  {
	return false, ""
}
