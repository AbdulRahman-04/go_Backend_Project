package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"Go_Backend/config"
)

// Claims struct for our JWT payload with correct JSON field name
type Claims struct {
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates the JWT token and attaches the user information to the request context.
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

		tokenString := strings.TrimSpace(parts[1])
		cfg := config.LoadConfig()
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtKey), nil
		})
		if err != nil || !token.Valid {
			log.Printf("JWT parse/validation error: %v", err)
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

		// Attach the user info for later use
		c.Set("user", claims.UserID)
		c.Next()
	}
}