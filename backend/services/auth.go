package services

import (
	"backend/crypto"
	"backend/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RegisterUser creates a new user account with wallet
func RegisterUser(fullName, email, cnic, hashedPassword string) (*models.User, string, error) {
	// Check if user already exists
	existingUser, _ := GetUserByEmail(email)
	if existingUser != nil {
		return nil, "", fmt.Errorf("user with email %s already exists", email)
	}

	// Generate keypair
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate keypair: %v", err)
	}

	// Convert keys to strings
	privateKeyStr := crypto.PrivateKeyToString(privateKey)
	publicKeyStr := crypto.PublicKeyToString(publicKey)

	// Encrypt private key
	encryptedPrivateKey, err := crypto.EncryptPrivateKey(privateKeyStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt private key: %v", err)
	}

	// Generate wallet ID from public key
	walletID := crypto.GenerateWalletID(publicKey)

	// Create user
	user := &models.User{
		ID:            uuid.New().String(),
		FullName:      fullName,
		Email:         email,
		Password:      hashedPassword,
		CNIC:          cnic,
		WalletID:      walletID,
		PublicKey:     publicKeyStr,
		PrivateKey:    encryptedPrivateKey,
		Beneficiaries: []string{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ZakatTracking: models.ZakatInfo{
			LastDeduction:   time.Time{},
			TotalDeducted:   0,
			MonthlyDeducted: 0,
		},
	}

	// Save user to database
	if err := SaveUser(user); err != nil {
		return nil, "", fmt.Errorf("failed to save user: %v", err)
	}

	// Create wallet
	wallet := &models.Wallet{
		WalletID:  walletID,
		UserID:    user.ID,
		PublicKey: publicKeyStr,
		Balance:   1000, // Initial balance
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	if err := SaveWallet(wallet); err != nil {
		return nil, "", fmt.Errorf("failed to create wallet: %v", err)
	}

	// Create initial UTXO with 1000 BC for new users
	initialTxHash := fmt.Sprintf("genesis-%s", user.ID)
	_, err = CreateUTXO(walletID, 1000.0, initialTxHash, 0)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create initial UTXO: %v", err)
	}

	// Log registration
	LogSystemEvent("registration", fmt.Sprintf("New user registered: %s with initial balance: 1000 BC", email), user.ID, "")

	return user, privateKeyStr, nil
}

// UpdateUserProfile updates user profile information
func UpdateUserProfile(userID, fullName, email string) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	// Check if email is being changed
	emailChanged := user.Email != email

	user.FullName = fullName
	user.Email = email
	user.UpdatedAt = time.Now()

	if err := UpdateUser(user); err != nil {
		return err
	}

	// Log profile update
	LogSystemEvent("profile_update", fmt.Sprintf("User %s updated profile", userID), userID, "")

	// If email changed, require re-verification
	if emailChanged {
		LogSystemEvent("email_change", fmt.Sprintf("User %s changed email to %s", userID, email), userID, "")
	}

	return nil
}

// AddBeneficiary adds a beneficiary wallet ID to user's list
func AddBeneficiary(userID, beneficiaryWalletID string) error {
	// Validate beneficiary wallet exists
	_, err := GetWalletByID(beneficiaryWalletID)
	if err != nil {
		return fmt.Errorf("invalid beneficiary wallet ID")
	}

	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	// Check if already added
	for _, id := range user.Beneficiaries {
		if id == beneficiaryWalletID {
			return fmt.Errorf("beneficiary already exists")
		}
	}

	user.Beneficiaries = append(user.Beneficiaries, beneficiaryWalletID)
	user.UpdatedAt = time.Now()

	return UpdateUser(user)
}

// RemoveBeneficiary removes a beneficiary from user's list
func RemoveBeneficiary(userID, beneficiaryWalletID string) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	newBeneficiaries := []string{}
	found := false

	for _, id := range user.Beneficiaries {
		if id != beneficiaryWalletID {
			newBeneficiaries = append(newBeneficiaries, id)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("beneficiary not found")
	}

	user.Beneficiaries = newBeneficiaries
	user.UpdatedAt = time.Now()

	return UpdateUser(user)
}

// GetUserWallet retrieves user's wallet information
func GetUserWallet(userID string) (*models.Wallet, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	wallet, err := GetWalletByID(user.WalletID)
	if err != nil {
		return nil, err
	}

	// Recalculate balance from UTXOs
	balance, err := CalculateBalance(wallet.WalletID)
	if err != nil {
		return nil, err
	}

	wallet.Balance = balance
	wallet.UpdatedAt = time.Now()

	// Update cached balance
	if err := UpdateWallet(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetUserByWalletID retrieves user by wallet ID
func GetUserByWalletID(walletID string) (*models.User, error) {
	wallet, err := GetWalletByID(walletID)
	if err != nil {
		return nil, err
	}

	return GetUserByID(wallet.UserID)
}

// ValidateCNIC validates Pakistani CNIC format (13 digits with hyphens)
func ValidateCNIC(cnic string) bool {
	if len(cnic) != 15 {
		return false
	}
	// Format: 12345-1234567-1
	if cnic[5] != '-' || cnic[13] != '-' {
		return false
	}
	// Check if all other characters are digits
	for i, ch := range cnic {
		if i != 5 && i != 13 {
			if ch < '0' || ch > '9' {
				return false
			}
		}
	}
	return true
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	if len(password) > 100 {
		return fmt.Errorf("password must not exceed 100 characters")
	}

	// Check for at least one letter and one number
	hasLetter := false
	hasNumber := false
	for _, ch := range password {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			hasLetter = true
		}
		if ch >= '0' && ch <= '9' {
			hasNumber = true
		}
	}
	if !hasLetter || !hasNumber {
		return fmt.Errorf("password must contain at least one letter and one number")
	}
	return nil
}
