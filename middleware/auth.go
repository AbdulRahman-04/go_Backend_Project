package middleware

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"Go_Backend/config"
)

// Claims struct for our JWT payload.
type Claims struct {
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

var (
	// jwtKey holds the secret key loaded once.
	jwtKey []byte
	// tokenCache caches token string to its claims.
	tokenCache sync.Map // map[string]*Claims
)

func init() {
	// Load configuration once and cache the JWT key.
	cfg := config.LoadConfig()
	jwtKey = []byte(cfg.JwtKey)
}

// getCachedClaims returns the cached claims for a token if still valid.
func getCachedClaims(tokenString string) (*Claims, bool) {
	if cached, ok := tokenCache.Load(tokenString); ok {
		if claims, ok := cached.(*Claims); ok {
			// Optional: verify that the token hasn't expired.
			if claims.ExpiresAt != nil && claims.ExpiresAt.Time.After(time.Now()) {
				return claims, true
			}
			// If expired, remove it from the cache.
			tokenCache.Delete(tokenString)
		}
	}
	return nil, false
}

// cacheClaims stores the claims for a token in the cache.
func cacheClaims(tokenString string, claims *Claims) {
	tokenCache.Store(tokenString, claims)
}

// AuthMiddleware validates the JWT token and attaches the user info to the request context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header.
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

		// Check if claims for this token are cached.
		if claims, found := getCachedClaims(tokenString); found {
			c.Set("user", claims.UserID)
			c.Next()
			return
		}

		// Parse the JWT token.
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
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

		// Cache the claims for future requests.
		cacheClaims(tokenString, claims)

		// Attach the user information for later use.
		c.Set("user", claims.UserID)
		c.Next()
	}
}