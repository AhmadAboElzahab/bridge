package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/AhmadAboElzahab/bridge/initializers"
	"github.com/AhmadAboElzahab/bridge/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Claims struct to represent the JWT claims
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix from the token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Declare the claims struct to hold the JWT payload
		claims := &Claims{}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT_SECRET is not set"})
			c.Abort()
			return
		}

		// Parse the JWT token and extract claims
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Retrieve user from the database based on the user ID in the claims
		var user models.User
		if err := initializers.DB.First(&user, claims.UserID).Error; err != nil {
			// If user is not found, return an unauthorized error
			c.JSON(http.StatusUnauthorized, gin.H{"error": claims.UserID})
			c.Abort()
			return
		}

		// Store the user in the context
		c.Set("user", user)

		// Continue to the next handler
		c.Next()
	}
}
