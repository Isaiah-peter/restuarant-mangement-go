package middlewares

import (
	"fmt"
	"golang-management-restaurant/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claim, err := helpers.ValidateAlltoken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claim.Email)
		c.Set("first_name", claim.FirstName)
		c.Set("last_name", claim.LastName)
		c.Set("user_id", claim.UserId)

		c.Next()
	}
}
