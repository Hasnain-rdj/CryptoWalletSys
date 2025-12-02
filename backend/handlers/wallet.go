package handlers

import (
	"backend/middleware"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetWallet returns wallet information
func GetWallet(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	wallet, err := services.GetUserWallet(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get UTXOs
	utxos, err := services.GetUTXOsByWallet(wallet.WalletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get UTXOs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet": wallet,
		"utxos":  utxos,
	})
}

// GetBalance returns wallet balance
func GetBalance(c *gin.Context) {
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

	balance, err := services.CalculateBalance(user.WalletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance":  balance,
		"walletId": user.WalletID,
	})
}

// ValidateWalletID validates if a wallet ID exists
func ValidateWalletID(c *gin.Context) {
	walletID := c.Param("walletId")
	if walletID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet ID required"})
		return
	}

	wallet, err := services.GetWalletByID(walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"valid":   false,
			"message": "Invalid Wallet ID",
		})
		return
	}

	// Get user info
	user, _ := services.GetUserByID(wallet.UserID)
	displayName := "Unknown"
	if user != nil {
		displayName = user.FullName
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":       true,
		"walletId":    wallet.WalletID,
		"displayName": displayName,
	})
}

// GetWalletUTXOs returns all UTXOs for a wallet
func GetWalletUTXOs(c *gin.Context) {
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

	utxos, err := services.GetUTXOsByWallet(user.WalletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"utxos": utxos,
		"count": len(utxos),
	})
}
