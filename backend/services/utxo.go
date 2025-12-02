package services

import (
	"backend/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// CreateUTXO creates a new UTXO
func CreateUTXO(walletID string, amount float64, txHash string, outputIndex int) (*models.UTXO, error) {
	utxo := &models.UTXO{
		ID:              uuid.New().String(),
		TransactionHash: txHash,
		OutputIndex:     outputIndex,
		WalletID:        walletID,
		Amount:          amount,
		Spent:           false,
		CreatedAt:       time.Now(),
	}

	if err := SaveUTXO(utxo); err != nil {
		return nil, err
	}

	return utxo, nil
}

// GetUTXOsByWallet retrieves all unspent UTXOs for a wallet
func GetUTXOsByWallet(walletID string) ([]models.UTXO, error) {
	return GetUnspentUTXOs(walletID)
}

// CalculateBalance calculates the balance from UTXOs
func CalculateBalance(walletID string) (float64, error) {
	utxos, err := GetUTXOsByWallet(walletID)
	if err != nil {
		return 0, err
	}

	balance := 0.0
	for _, utxo := range utxos {
		if !utxo.Spent {
			balance += utxo.Amount
		}
	}

	return balance, nil
}

// SelectUTXOs selects UTXOs to spend for a given amount
func SelectUTXOs(walletID string, amount float64) ([]models.UTXO, float64, error) {
	utxos, err := GetUTXOsByWallet(walletID)
	if err != nil {
		return nil, 0, err
	}

	var selectedUTXOs []models.UTXO
	total := 0.0

	for _, utxo := range utxos {
		if !utxo.Spent {
			selectedUTXOs = append(selectedUTXOs, utxo)
			total += utxo.Amount

			if total >= amount {
				break
			}
		}
	}

	if total < amount {
		return nil, 0, fmt.Errorf("insufficient balance: need %.2f, have %.2f", amount, total)
	}

	return selectedUTXOs, total, nil
}

// SpendUTXO marks a UTXO as spent
func SpendUTXO(utxoID string, txHash string) error {
	utxo, err := GetUTXOByID(utxoID)
	if err != nil {
		return err
	}

	if utxo.Spent {
		return fmt.Errorf("UTXO %s already spent", utxoID)
	}

	utxo.Spent = true
	utxo.SpentInTxHash = txHash
	utxo.SpentAt = time.Now()

	return UpdateUTXO(utxo)
}

// ValidateUTXOs checks if UTXOs are valid and unspent
func ValidateUTXOs(utxoIDs []string) error {
	for _, utxoID := range utxoIDs {
		utxo, err := GetUTXOByID(utxoID)
		if err != nil {
			return fmt.Errorf("UTXO %s not found", utxoID)
		}

		if utxo.Spent {
			return fmt.Errorf("UTXO %s already spent in transaction %s", utxoID, utxo.SpentInTxHash)
		}
	}

	return nil
}

// ProcessTransactionUTXOs handles UTXO creation and spending for a transaction
func ProcessTransactionUTXOs(tx models.Transaction) error {
	// Mark input UTXOs as spent
	for _, utxoID := range tx.InputUTXOs {
		if err := SpendUTXO(utxoID, tx.Hash); err != nil {
			return fmt.Errorf("failed to spend UTXO %s: %v", utxoID, err)
		}
	}

	// Create output UTXOs
	for idx, output := range tx.OutputUTXOs {
		_, err := CreateUTXO(output.WalletID, output.Amount, tx.Hash, idx)
		if err != nil {
			return fmt.Errorf("failed to create output UTXO: %v", err)
		}
	}

	log.Printf("Processed UTXOs for transaction %s", tx.Hash)
	return nil
}

// GetTotalSupply calculates the total supply of cryptocurrency
func GetTotalSupply() (float64, error) {
	allUTXOs, err := GetAllUTXOs()
	if err != nil {
		return 0, err
	}

	total := 0.0
	for _, utxo := range allUTXOs {
		if !utxo.Spent {
			total += utxo.Amount
		}
	}

	return total, nil
}

// RecalculateWalletBalance recalculates and updates wallet balance
func RecalculateWalletBalance(walletID string) error {
	balance, err := CalculateBalance(walletID)
	if err != nil {
		return err
	}

	wallet, err := GetWalletByID(walletID)
	if err != nil {
		return err
	}

	wallet.Balance = balance
	wallet.UpdatedAt = time.Now()

	return UpdateWallet(wallet)
}
