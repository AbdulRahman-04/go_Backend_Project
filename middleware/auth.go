package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"Go_Backend/config"
)

// Claims struct for our JWT payload
type Claims struct {
	UserID string `json:"user"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates the JWT token and attaches the user information to the request context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "No token provided ❌"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Token format is invalid ❌"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// Load config to get the JWT key
		cfg := config.LoadConfig()
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Use the JwtKey from our configuration file
			return []byte(cfg.JwtKey), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token ❌"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Failed to parse token claims ❌"})
			c.Abort()
			return
		}

		// Attach user info to the context (for later use in controllers)
		c.Set("user", claims.UserID)
		c.Next()
	}
}