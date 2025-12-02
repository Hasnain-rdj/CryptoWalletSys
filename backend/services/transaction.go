package services

import (
	"backend/crypto"
	"backend/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// CreateTransaction creates a new transaction
func CreateTransaction(senderWalletID, receiverWalletID string, amount float64, note, senderPublicKey, privateKeyStr string) (*models.Transaction, error) {
	// Prevent self-transfer
	if senderWalletID == receiverWalletID {
		return nil, fmt.Errorf("cannot send money to yourself")
	}

	// Validate minimum amount
	if amount < 0.01 {
		return nil, fmt.Errorf("minimum transaction amount is 0.01 BC")
	}

	// Validate sender wallet exists
	senderWallet, err := GetWalletByID(senderWalletID)
	if err != nil {
		return nil, fmt.Errorf("invalid sender wallet ID: %v", err)
	}

	// Validate receiver wallet exists
	_, err = GetWalletByID(receiverWalletID)
	if err != nil {
		return nil, fmt.Errorf("invalid receiver wallet ID: %v", err)
	}

	// Check balance
	balance, err := CalculateBalance(senderWalletID)
	if err != nil {
		return nil, err
	}

	if balance < amount {
		return nil, fmt.Errorf("insufficient balance: have %.2f, need %.2f", balance, amount)
	}

	// Select UTXOs to spend
	selectedUTXOs, total, err := SelectUTXOs(senderWalletID, amount)
	if err != nil {
		return nil, err
	}

	// Create transaction
	timestamp := time.Now()
	tx := &models.Transaction{
		Hash:             uuid.New().String(),
		SenderWalletID:   senderWalletID,
		ReceiverWalletID: receiverWalletID,
		Amount:           amount,
		Note:             note,
		Timestamp:        timestamp,
		SenderPublicKey:  senderPublicKey,
		Type:             "transfer",
		Status:           "pending",
	}

	// Create signature payload
	payload := crypto.CreateTransactionPayload(
		senderWalletID,
		receiverWalletID,
		amount,
		timestamp.String(),
		note,
	)

	// Sign transaction
	privateKey, err := crypto.StringToPrivateKey(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	signature, err := crypto.SignData(payload, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	tx.Signature = signature

	// Add input UTXOs
	for _, utxo := range selectedUTXOs {
		tx.InputUTXOs = append(tx.InputUTXOs, utxo.ID)
	}

	// Create output UTXOs
	// 1. Payment to receiver
	tx.OutputUTXOs = append(tx.OutputUTXOs, models.UTXOOutput{
		WalletID: receiverWalletID,
		Amount:   amount,
	})

	// 2. Change back to sender (if any)
	change := total - amount
	if change > 0 {
		tx.OutputUTXOs = append(tx.OutputUTXOs, models.UTXOOutput{
			WalletID: senderWalletID,
			Amount:   change,
		})
	}

	// Calculate transaction hash
	tx.Hash = crypto.HashSHA256(fmt.Sprintf("%s%s%.8f%d", senderWalletID, receiverWalletID, amount, timestamp.Unix()))

	// Update wallet balance (cached)
	senderWallet.Balance -= amount
	senderWallet.UpdatedAt = time.Now()
	if err := UpdateWallet(senderWallet); err != nil {
		log.Printf("Warning: failed to update sender wallet balance: %v", err)
	}

	return tx, nil
}

// ValidateTransaction validates a transaction
func ValidateTransaction(tx models.Transaction) error {
	// 1. Validate wallet IDs exist
	_, err := GetWalletByID(tx.SenderWalletID)
	if err != nil {
		return fmt.Errorf("invalid sender wallet ID")
	}

	_, err = GetWalletByID(tx.ReceiverWalletID)
	if err != nil {
		return fmt.Errorf("invalid receiver wallet ID")
	}

	// 2. Verify digital signature
	payload := crypto.CreateTransactionPayload(
		tx.SenderWalletID,
		tx.ReceiverWalletID,
		tx.Amount,
		tx.Timestamp.String(),
		tx.Note,
	)

	publicKey, err := crypto.StringToPublicKey(tx.SenderPublicKey)
	if err != nil {
		return fmt.Errorf("invalid public key: %v", err)
	}

	if err := crypto.VerifySignature(payload, tx.Signature, publicKey); err != nil {
		return fmt.Errorf("invalid signature: %v", err)
	}

	// 3. Validate UTXOs
	if err := ValidateUTXOs(tx.InputUTXOs); err != nil {
		return fmt.Errorf("invalid UTXOs: %v", err)
	}

	// 4. Check if inputs cover outputs
	inputTotal := 0.0
	for _, utxoID := range tx.InputUTXOs {
		utxo, err := GetUTXOByID(utxoID)
		if err != nil {
			return fmt.Errorf("UTXO not found: %v", err)
		}
		inputTotal += utxo.Amount
	}

	outputTotal := 0.0
	for _, output := range tx.OutputUTXOs {
		outputTotal += output.Amount
	}

	if inputTotal < outputTotal {
		return fmt.Errorf("insufficient inputs: have %.2f, need %.2f", inputTotal, outputTotal)
	}

	return nil
}

// ProcessTransaction processes and adds transaction to pending pool
func ProcessTransaction(tx models.Transaction) error {
	// Validate transaction
	if err := ValidateTransaction(tx); err != nil {
		// Log failed transaction
		LogSystemEvent("validation_failure", fmt.Sprintf("Transaction validation failed: %v", err), "", "")
		return err
	}

	// Save transaction to database
	if err := SaveTransaction(&tx); err != nil {
		return err
	}

	// Add to pending pool
	if err := AddPendingTransaction(tx); err != nil {
		return err
	}

	// Log transaction
	LogTransactionEvent(tx.Hash, "pending", "", tx.SenderWalletID, "", tx.Amount, "pending")

	log.Printf("Transaction %s processed and added to pending pool", tx.Hash)
	return nil
}

// GetTransactionsByWallet retrieves all transactions for a wallet
func GetTransactionsByWallet(walletID string) ([]models.Transaction, error) {
	return GetWalletTransactions(walletID)
}

// GetTransactionByHash retrieves a transaction by its hash
func GetTransactionByHash(hash string) (*models.Transaction, error) {
	return GetTransactionFromDB(hash)
}

// CreateZakatTransaction creates a zakat deduction transaction
func CreateZakatTransaction(walletID string, amount float64, month string) (*models.Transaction, error) {
	zakatPoolWallet := GetZakatPoolWallet()

	// Select UTXOs
	selectedUTXOs, total, err := SelectUTXOs(walletID, amount)
	if err != nil {
		return nil, err
	}

	tx := &models.Transaction{
		Hash:             uuid.New().String(),
		SenderWalletID:   walletID,
		ReceiverWalletID: zakatPoolWallet,
		Amount:           amount,
		Note:             fmt.Sprintf("Zakat deduction for %s", month),
		Timestamp:        time.Now(),
		Type:             "zakat_deduction",
		Status:           "pending",
	}

	// Add input UTXOs
	for _, utxo := range selectedUTXOs {
		tx.InputUTXOs = append(tx.InputUTXOs, utxo.ID)
	}

	// Create output UTXOs
	tx.OutputUTXOs = append(tx.OutputUTXOs, models.UTXOOutput{
		WalletID: zakatPoolWallet,
		Amount:   amount,
	})

	// Change back to sender
	change := total - amount
	if change > 0 {
		tx.OutputUTXOs = append(tx.OutputUTXOs, models.UTXOOutput{
			WalletID: walletID,
			Amount:   change,
		})
	}

	return tx, nil
}
