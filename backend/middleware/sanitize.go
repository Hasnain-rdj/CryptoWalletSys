package middleware

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// Sanitize input strings to prevent injection attacks
func SanitizeInput(input string) string {
	// Remove potential SQL injection patterns
	sqlPatterns := []string{
		`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute|script|javascript|<script)`,
		`--|;|\/\*|\*\/|xp_|sp_`,
	}

	sanitized := input
	for _, pattern := range sqlPatterns {
		re := regexp.MustCompile(pattern)
		sanitized = re.ReplaceAllString(sanitized, "")
	}

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// SanitizeMiddleware sanitizes request inputs
func SanitizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				values[i] = SanitizeInput(value)
			}
			c.Request.URL.Query()[key] = values
		}

		c.Next()
	}
}

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateWalletID checks if wallet ID format is valid
func ValidateWalletID(walletID string) bool {
	// Wallet ID should be alphanumeric and at least 10 characters
	if len(walletID) < 10 {
		return false
	}
	walletRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return walletRegex.MatchString(walletID)
}
