package middleware

import (
	"golang-sosmed-gin/errorhandler"
	"golang-sosmed-gin/helper"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			errorhandler.HandleError(c, &errorhandler.UnauthorizedError{Message: "Unauthorize!"})
			c.Abort()
			return
		}

		userID, err := helper.ValidateToken(tokenString)
		if err != nil {
			errorhandler.HandleError(c, &errorhandler.UnauthorizedError{Message: "Unauthorize!!"})
			c.Abort()
			return
		}

		c.Set("UserID", *userID)
		c.Next()
	}

}
