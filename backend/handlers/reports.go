package handlers

import (
	"backend/middleware"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetZakatHistory returns zakat deduction history for user
func GetZakatHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	history, err := services.GetZakatHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"count":   len(history),
	})
}

// GetZakatSummary returns zakat summary for user
func GetZakatSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	summary, err := services.GetZakatSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// TriggerZakatDeduction manually triggers zakat deduction (admin/testing)
func TriggerZakatDeduction(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := services.ManualZakatDeduction(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Zakat deduction processed successfully",
	})
}

// GetSystemLogs returns system logs
func GetSystemLogs(c *gin.Context) {
	logType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "100")

	limit, _ := strconv.Atoi(limitStr)

	logs, err := services.GetSystemLogs(logType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"count": len(logs),
	})
}

// GetTransactionLogs returns transaction logs for user
func GetTransactionLogs(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := services.GetUserTransactionLogs(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"count": len(logs),
	})
}

// GetAllTransactionLogs returns all transaction logs (admin)
func GetAllTransactionLogs(c *gin.Context) {
	// For now, just return empty as we need user-specific logs
	c.JSON(http.StatusOK, gin.H{
		"logs":    []interface{}{},
		"count":   0,
		"message": "Use user-specific transaction logs endpoint",
	})
}

// GetReports generates various reports
func GetReports(c *gin.Context) {
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

	// Get transactions
	transactions, _ := services.GetTransactionsByWallet(user.WalletID)

	// Calculate totals
	totalSent := 0.0
	totalReceived := 0.0

	for _, tx := range transactions {
		if tx.SenderWalletID == user.WalletID {
			totalSent += tx.Amount
		}
		if tx.ReceiverWalletID == user.WalletID {
			totalReceived += tx.Amount
		}
	}

	// Get zakat summary
	zakatSummary, _ := services.GetZakatSummary(userID)

	// Get current balance
	balance, _ := services.CalculateBalance(user.WalletID)

	report := map[string]interface{}{
		"walletId":         user.WalletID,
		"currentBalance":   balance,
		"totalSent":        totalSent,
		"totalReceived":    totalReceived,
		"transactionCount": len(transactions),
		"zakatSummary":     zakatSummary,
	}

	c.JSON(http.StatusOK, report)
}
