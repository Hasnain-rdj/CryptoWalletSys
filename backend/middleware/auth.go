package middleware

import (
	"backend/services"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			services.LogSystemEvent("auth_failure", "Missing authorization header", "", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			services.LogSystemEvent("auth_failure", "Invalid authorization format", "", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token length (prevent empty or malformed tokens)
		if len(tokenString) < 20 {
			services.LogSystemEvent("auth_failure", "Token too short", "", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Parse and validate JWT token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			services.LogSystemEvent("auth_failure", "Invalid token: "+err.Error(), "", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			// Validate claims are not empty
			if claims.UserID == "" || claims.Email == "" {
				services.LogSystemEvent("auth_failure", "Empty claims in token", "", c.ClientIP())
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}

			// Set user info in context
			c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)
			c.Next()
		} else {
			services.LogSystemEvent("auth_failure", "Invalid token claims", "", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) string {
	userID, exists := c.Get("userID")
	if !exists {
		return ""
	}
	return userID.(string)
}

// GetEmail retrieves email from context
func GetEmail(c *gin.Context) string {
	email, exists := c.Get("email")
	if !exists {
		return ""
	}
	return email.(string)
}
