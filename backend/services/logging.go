package services

import (
	"backend/models"
	"log"
	"time"

	"github.com/google/uuid"
)

// LogSystemEvent logs a system event
func LogSystemEvent(logType, message, userID, ipAddress string) {
	systemLog := &models.SystemLog{
		ID:        uuid.New().String(),
		Type:      logType,
		Message:   message,
		UserID:    userID,
		IPAddress: ipAddress,
		Timestamp: time.Now(),
	}

	if err := SaveSystemLog(systemLog); err != nil {
		log.Printf("Error saving system log: %v", err)
	}
}

// LogSystemEventWithMetadata logs a system event with additional metadata
func LogSystemEventWithMetadata(logType, message, userID, ipAddress string, metadata map[string]interface{}) {
	systemLog := &models.SystemLog{
		ID:        uuid.New().String(),
		Type:      logType,
		Message:   message,
		UserID:    userID,
		IPAddress: ipAddress,
		Timestamp: time.Now(),
		Metadata:  metadata,
	}

	if err := SaveSystemLog(systemLog); err != nil {
		log.Printf("Error saving system log with metadata: %v", err)
	}
}

// LogTransactionEvent logs a transaction-specific event
func LogTransactionEvent(txHash, action, userID, walletID, ipAddress string, amount float64, status string) {
	txLog := &models.TransactionLog{
		ID:              uuid.New().String(),
		TransactionHash: txHash,
		Action:          action,
		UserID:          userID,
		WalletID:        walletID,
		Amount:          amount,
		Status:          status,
		IPAddress:       ipAddress,
		Timestamp:       time.Now(),
	}

	if err := SaveTransactionLog(txLog); err != nil {
		log.Printf("Error saving transaction log: %v", err)
	}
}
