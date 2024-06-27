package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
	"github.com/mananKoyawala/whatsapp-clone/internal/user"
)

var unauthorized string = "permission denied"

func AuthMiddleware(userHandler *user.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("X-AUTH-TOKEN")

		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorized})
			c.Abort()
			return
		}

		// validate token claims, expiry
		claims, msg := helper.ValidateToken(clientToken)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorized})
			c.Abort()
			return
		}

		// check token is token or not
		if claims.TokenType != "token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorized})
			c.Abort()
			return
		}

		// check the user exits with the id
		user, err := userHandler.Service.GetUserById(c.Request.Context(), claims.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorized})
			c.Abort()
			return
		}

		if claims.ID != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorized})
			c.Abort()
			return
		}

		c.Set("id", user.ID)
		c.Set("name", user.Name)
		c.Set("mobile", user.Mobile)
		c.Set("about", user.About)
		c.Set("image", user.Image)
		c.Next()

	}
}
