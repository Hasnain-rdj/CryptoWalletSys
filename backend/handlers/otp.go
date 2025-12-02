package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type OTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

// In-memory OTP storage (in production, use Redis or database)
var otpStore = make(map[string]OTPData)

type OTPData struct {
	Code      string
	ExpiresAt time.Time
	Verified  bool
}

// GenerateOTP generates and sends OTP to email
func GenerateOTP(c *gin.Context) {
	var req OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	existingUser, err := services.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Generate 6-digit OTP
	otp, err := generateRandomOTP(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	// Store OTP with 5 minute expiration
	otpStore[req.Email] = OTPData{
		Code:      otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Verified:  false,
	}

	// Send OTP via email
	err = services.SendOTPEmail(req.Email, otp)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to send OTP email: %v\n", err)
	}

	// Check if we're in development mode (no SMTP configured)
	isDevelopment := os.Getenv("SMTP_HOST") == ""

	response := gin.H{"message": "OTP sent successfully"}
	if isDevelopment {
		// Only return OTP in development mode
		response["otp"] = otp
		response["dev_mode"] = true
	}

	c.JSON(http.StatusOK, response)
}

// VerifyOTP verifies the OTP code
func VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if OTP exists
	otpData, exists := otpStore[req.Email]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "OTP not found. Please request a new one."})
		return
	}

	// Check if OTP expired
	if time.Now().After(otpData.ExpiresAt) {
		delete(otpStore, req.Email)
		c.JSON(http.StatusGone, gin.H{"error": "OTP has expired. Please request a new one."})
		return
	}

	// Verify OTP
	if otpData.Code != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	// Mark as verified
	otpData.Verified = true
	otpStore[req.Email] = otpData

	c.JSON(http.StatusOK, gin.H{
		"message":  "Email verified successfully",
		"verified": true,
	})
}

// CheckEmailVerification checks if email is verified
func CheckEmailVerification(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	otpData, exists := otpStore[email]
	if !exists {
		c.JSON(http.StatusOK, gin.H{"verified": false})
		return
	}

	// Check if verified and not expired
	verified := otpData.Verified && time.Now().Before(otpData.ExpiresAt)

	c.JSON(http.StatusOK, gin.H{"verified": verified})
}

// Helper function to generate random OTP
func generateRandomOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}
	return string(otp), nil
}

// ClearOTPAfterRegistration clears OTP after successful registration
func ClearOTPAfterRegistration(email string) {
	delete(otpStore, email)
}
