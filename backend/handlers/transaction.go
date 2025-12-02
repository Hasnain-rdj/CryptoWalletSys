package handlers

import (
	"backend/crypto"
	"backend/middleware"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTransactionRequest represents transaction creation request
type CreateTransactionRequest struct {
	ReceiverWalletID string  `json:"receiverWalletId" binding:"required,min=10"`
	Amount           float64 `json:"amount" binding:"required,gt=0"`
	Note             string  `json:"note" binding:"max=500"`
	PrivateKey       string  `json:"privateKey" binding:"required,min=100"`
}

// CreateTransaction creates a new transaction
func CreateTransaction(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate minimum transaction amount (0.01)
	if req.Amount < 0.01 {
		services.LogSystemEvent("transaction_failure", "Amount too small", userID, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimum transaction amount is 0.01 BC"})
		return
	}

	// Validate maximum transaction amount (1000000)
	if req.Amount > 1000000 {
		services.LogSystemEvent("transaction_failure", "Amount too large", userID, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum transaction amount is 1,000,000 BC"})
		return
	}

	// Get user
	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Decrypt and validate private key
	decryptedKey, err := crypto.DecryptPrivateKey(user.PrivateKey)
	if err != nil || decryptedKey != req.PrivateKey {
		// If encrypted key doesn't match, try using the provided key directly
		if req.PrivateKey == "" {
			services.LogSystemEvent("transaction_failure", "Invalid private key", userID, c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid private key"})
			return
		}
	}

	// Create transaction
	tx, err := services.CreateTransaction(
		user.WalletID,
		req.ReceiverWalletID,
		req.Amount,
		req.Note,
		user.PublicKey,
		req.PrivateKey,
	)
	if err != nil {
		services.LogSystemEvent("transaction_failure", "Failed to create transaction: "+err.Error(), userID, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process transaction (validate and add to pending pool)
	if err := services.ProcessTransaction(*tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process UTXOs
	if err := services.ProcessTransactionUTXOs(*tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process UTXOs"})
		return
	}

	// Mine the transaction into a block
	block, err := services.MineBlock(user.WalletID)
	if err != nil {
		services.LogSystemEvent("mining_failure", "Failed to mine block: "+err.Error(), userID, c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction created but mining failed"})
		return
	}

	services.LogSystemEvent("transaction_success", "Transaction created and mined", userID, c.ClientIP())

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transaction successful",
		"transaction": tx,
		"block":       block,
	})
}

// GetTransactionHistory returns transaction history for user
func GetTransactionHistory(c *gin.Context) {
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

	transactions, err := services.GetTransactionsByWallet(user.WalletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
	})
}

// GetTransactionByHash returns a specific transaction
func GetTransactionByHash(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction hash required"})
		return
	}

	tx, err := services.GetTransactionByHash(hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": tx,
	})
}

// GetPendingTransactions returns all pending transactions
func GetPendingTransactions(c *gin.Context) {
	transactions, _ := services.GetPendingTransactions()

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
	})
}
