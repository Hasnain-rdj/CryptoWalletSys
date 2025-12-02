package services

import (
	"backend/models"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var zakatPoolWalletID string

// InitZakatPool initializes the zakat pool wallet
func InitZakatPool() {
	zakatPoolWalletID = os.Getenv("ZAKAT_POOL_WALLET_ID")
	if zakatPoolWalletID == "" {
		zakatPoolWalletID = "ZAKAT_POOL_WALLET"
	}
}

// GetZakatPoolWallet returns the zakat pool wallet ID
func GetZakatPoolWallet() string {
	if zakatPoolWalletID == "" {
		InitZakatPool()
	}
	return zakatPoolWalletID
}

// StartZakatScheduler starts the monthly zakat deduction scheduler
func StartZakatScheduler() {
	InitZakatPool()

	log.Println("Zakat scheduler started")

	// Run zakat deduction on the 1st of every month
	ticker := time.NewTicker(24 * time.Hour) // Check daily
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		// Check if it's the 1st of the month and hasn't run yet today
		if now.Day() == 1 {
			log.Println("Running monthly zakat deduction...")
			if err := ProcessMonthlyZakat(); err != nil {
				log.Printf("Error processing monthly zakat: %v", err)
				LogSystemEvent("zakat_error", fmt.Sprintf("Failed to process monthly zakat: %v", err), "", "")
			}
		}
	}
}

// ProcessMonthlyZakat processes zakat deduction for all users
func ProcessMonthlyZakat() error {
	users, err := GetAllUsers()
	if err != nil {
		return err
	}

	zakatPercentage := getZakatPercentage()
	currentMonth := time.Now().Format("2006-01")

	successCount := 0
	failCount := 0

	for _, user := range users {
		if err := deductZakat(user, zakatPercentage, currentMonth); err != nil {
			log.Printf("Failed to deduct zakat for user %s: %v", user.ID, err)
			failCount++
		} else {
			successCount++
		}
	}

	log.Printf("Zakat deduction completed: %d successful, %d failed", successCount, failCount)

	// Log system event
	LogSystemEvent("zakat_deduction",
		fmt.Sprintf("Monthly zakat deduction completed: %d successful, %d failed", successCount, failCount),
		"", "")

	return nil
}

// deductZakat deducts zakat from a single user
func deductZakat(user models.User, percentage float64, month string) error {
	// Check if already deducted this month
	if user.ZakatTracking.LastDeduction.Format("2006-01") == month {
		log.Printf("Zakat already deducted for user %s in month %s", user.ID, month)
		return nil
	}

	// Calculate balance
	balance, err := CalculateBalance(user.WalletID)
	if err != nil {
		return err
	}

	// Skip if balance is zero or negative
	if balance <= 0 {
		log.Printf("User %s has zero balance, skipping zakat", user.ID)
		return nil
	}

	// Calculate zakat amount
	zakatAmount := balance * (percentage / 100)

	if zakatAmount < 0.01 {
		log.Printf("Zakat amount too small for user %s, skipping", user.ID)
		return nil
	}

	// Create zakat transaction
	tx, err := CreateZakatTransaction(user.WalletID, zakatAmount, month)
	if err != nil {
		return err
	}

	// Process transaction
	if err := ProcessTransaction(*tx); err != nil {
		return err
	}

	// Process UTXOs
	if err := ProcessTransactionUTXOs(*tx); err != nil {
		return err
	}

	// Record zakat deduction
	zakatDeduction := models.ZakatDeduction{
		ID:              uuid.New().String(),
		UserID:          user.ID,
		WalletID:        user.WalletID,
		Amount:          zakatAmount,
		BalanceBefore:   balance,
		BalanceAfter:    balance - zakatAmount,
		TransactionHash: tx.Hash,
		Month:           month,
		DeductedAt:      time.Now(),
		Status:          "completed",
	}

	if err := SaveZakatDeduction(&zakatDeduction); err != nil {
		return err
	}

	// Update user's zakat tracking
	user.ZakatTracking.LastDeduction = time.Now()
	user.ZakatTracking.TotalDeducted += zakatAmount
	user.ZakatTracking.MonthlyDeducted = zakatAmount
	user.UpdatedAt = time.Now()

	if err := UpdateUser(&user); err != nil {
		return err
	}

	// Recalculate wallet balance
	if err := RecalculateWalletBalance(user.WalletID); err != nil {
		log.Printf("Warning: failed to recalculate balance: %v", err)
	}

	log.Printf("Zakat deducted for user %s: %.2f (%.2f%%)", user.ID, zakatAmount, percentage)

	// Log transaction
	LogTransactionEvent(tx.Hash, "zakat_deducted", user.ID, user.WalletID, "", zakatAmount, "completed")

	return nil
}

// getZakatPercentage returns the zakat percentage from environment
func getZakatPercentage() float64 {
	percentStr := os.Getenv("ZAKAT_PERCENTAGE")
	if percentStr == "" {
		return 2.5
	}

	percent, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		return 2.5
	}

	return percent
}

// GetZakatHistory retrieves zakat history for a user
func GetZakatHistory(userID string) ([]models.ZakatDeduction, error) {
	return GetUserZakatDeductions(userID)
}

// GetZakatSummary calculates zakat summary for a user
func GetZakatSummary(userID string) (map[string]interface{}, error) {
	history, err := GetZakatHistory(userID)
	if err != nil {
		return nil, err
	}

	totalDeducted := 0.0
	monthlyDeductions := make(map[string]float64)

	for _, deduction := range history {
		totalDeducted += deduction.Amount
		monthlyDeductions[deduction.Month] = deduction.Amount
	}

	summary := map[string]interface{}{
		"totalDeducted":     totalDeducted,
		"monthlyDeductions": monthlyDeductions,
		"deductionCount":    len(history),
		"lastDeduction":     nil,
	}

	if len(history) > 0 {
		summary["lastDeduction"] = history[len(history)-1]
	}

	return summary, nil
}

// ManualZakatDeduction allows manual zakat deduction (for testing or admin)
func ManualZakatDeduction(userID string) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	zakatPercentage := getZakatPercentage()
	currentMonth := time.Now().Format("2006-01")

	return deductZakat(*user, zakatPercentage, currentMonth)
}
