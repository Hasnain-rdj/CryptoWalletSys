package handlers

import (
	"backend/middleware"
	"backend/services"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents registration request
type RegisterRequest struct {
	FullName string `json:"fullName" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	CNIC     string `json:"cnic" binding:"required,len=15"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// RegisterResponse represents registration response
type RegisterResponse struct {
	User       interface{} `json:"user"`
	Token      string      `json:"token"`
	PrivateKey string      `json:"privateKey"`
	Message    string      `json:"message"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate CNIC format (must be exactly 13 digits with hyphens: 12345-1234567-1)
	if !services.ValidateCNIC(req.CNIC) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CNIC format. Expected format: 12345-1234567-1"})
		return
	}

	// Validate password strength
	if err := services.ValidatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email is verified
	otpData, exists := otpStore[req.Email]
	if !exists || !otpData.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified. Please verify your email first."})
		return
	}

	// Check if user already exists
	existingUser, _ := services.GetUserByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		services.LogSystemEvent("registration_failure", "Failed to hash password", "", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user and wallet in our system
	user, privateKey, err := services.RegisterUser(req.FullName, req.Email, req.CNIC, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clear OTP after successful registration
	ClearOTPAfterRegistration(req.Email)

	// Send welcome email (async, don't block registration)
	go services.SendWelcomeEmail(req.Email, req.FullName)

	// Generate JWT token
	token, err := generateJWT(user.ID, user.Email)
	if err != nil {
		services.LogSystemEvent("token_generation_failure", "Failed to generate JWT", user.ID, c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	services.LogSystemEvent("registration_success", "User registered successfully", user.ID, c.ClientIP())

	c.JSON(http.StatusCreated, RegisterResponse{
		User:       user,
		Token:      token,
		PrivateKey: privateKey,
		Message:    "User registered successfully. Please save your private key securely!",
	})
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User    interface{} `json:"user"`
	Token   string      `json:"token"`
	Message string      `json:"message"`
}

// Login handles user login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from database
	user, err := services.GetUserByEmail(req.Email)
	if err != nil {
		services.LogSystemEvent("login_failure", "User not found: "+req.Email, "", c.ClientIP())
		c.JSON(http.StatusNotFound, gin.H{"error": "User doesn't exist. Please sign up first."})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		services.LogSystemEvent("login_failure", "Invalid password for: "+req.Email, user.ID, c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID, user.Email)
	if err != nil {
		services.LogSystemEvent("token_generation_failure", "Failed to generate JWT", user.ID, c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	services.LogSystemEvent("login_success", "User logged in", user.ID, c.ClientIP())

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		User:    user,
		Token:   token,
		Message: "Login successful",
	})
}

// GetProfile returns the authenticated user's profile
func GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get wallet info
	wallet, err := services.GetUserWallet(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallet"})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"wallet": wallet,
	})
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FullName string `json:"fullName" binding:"required"`
}

// UpdateProfile updates user profile
func UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update profile
	user.FullName = req.FullName
	user.UpdatedAt = time.Now()

	if err := services.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	services.LogSystemEvent("profile_update", "Profile updated", userID, c.ClientIP())

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// AddBeneficiaryRequest represents add beneficiary request
type AddBeneficiaryRequest struct {
	BeneficiaryWalletID string `json:"beneficiaryWalletId" binding:"required"`
}

// AddBeneficiary adds a beneficiary to user's list
func AddBeneficiary(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req AddBeneficiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate wallet ID exists
	_, err := services.GetWalletByID(req.BeneficiaryWalletID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if already added
	for _, b := range user.Beneficiaries {
		if b == req.BeneficiaryWalletID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Beneficiary already exists"})
			return
		}
	}

	// Add beneficiary
	user.Beneficiaries = append(user.Beneficiaries, req.BeneficiaryWalletID)
	user.UpdatedAt = time.Now()

	if err := services.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add beneficiary"})
		return
	}

	services.LogSystemEvent("beneficiary_added", "Beneficiary added: "+req.BeneficiaryWalletID, userID, c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"message":       "Beneficiary added successfully",
		"beneficiaries": user.Beneficiaries,
	})
}

// RemoveBeneficiary removes a beneficiary from user's list
func RemoveBeneficiary(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	beneficiaryID := c.Param("id")
	if beneficiaryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Beneficiary ID required"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Remove beneficiary
	newBeneficiaries := []string{}
	found := false
	for _, b := range user.Beneficiaries {
		if b != beneficiaryID {
			newBeneficiaries = append(newBeneficiaries, b)
		} else {
			found = true
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Beneficiary not found"})
		return
	}

	user.Beneficiaries = newBeneficiaries
	user.UpdatedAt = time.Now()

	if err := services.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove beneficiary"})
		return
	}

	services.LogSystemEvent("beneficiary_removed", "Beneficiary removed: "+beneficiaryID, userID, c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"message":       "Beneficiary removed successfully",
		"beneficiaries": user.Beneficiaries,
	})
}

// generateJWT generates a JWT token for a user
func generateJWT(userID, email string) (string, error) {
	claims := middleware.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "blockchain-wallet",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
