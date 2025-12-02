package services

import (
	"backend/crypto"
	"backend/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	blockchain          []models.Block
	pendingTransactions []models.Transaction
	blockchainMutex     sync.RWMutex
	difficulty          int
)

// InitBlockchain initializes the blockchain with genesis block
func InitBlockchain() {
	diffStr := os.Getenv("MINING_DIFFICULTY")
	if diffStr == "" {
		difficulty = 4
	} else {
		difficulty, _ = strconv.Atoi(diffStr)
	}

	// Check if blockchain exists in database
	blocks, err := GetAllBlocks()
	if err != nil || len(blocks) == 0 {
		// Create genesis block
		genesisBlock := createGenesisBlock()
		blockchain = append(blockchain, genesisBlock)

		// Save to database
		if err := SaveBlock(genesisBlock); err != nil {
			log.Printf("Error saving genesis block: %v", err)
		}

		log.Println("Genesis block created")
	} else {
		blockchain = blocks
		log.Printf("Loaded %d blocks from database", len(blocks))
	}
}

// createGenesisBlock creates the first block in the blockchain
func createGenesisBlock() models.Block {
	genesisTransaction := models.Transaction{
		Hash:             "genesis",
		SenderWalletID:   "GENESIS",
		ReceiverWalletID: os.Getenv("ZAKAT_POOL_WALLET_ID"),
		Amount:           0,
		Note:             "Genesis Block",
		Timestamp:        time.Now(),
		Type:             "genesis",
		Status:           "confirmed",
	}

	block := models.Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: []models.Transaction{genesisTransaction},
		PreviousHash: "0",
		Nonce:        0,
		Difficulty:   difficulty,
	}

	block.MerkleRoot = calculateMerkleRoot(block.Transactions)
	block.Hash = calculateBlockHash(block)

	return block
}

// GetLatestBlock returns the last block in the chain
func GetLatestBlock() models.Block {
	blockchainMutex.RLock()
	defer blockchainMutex.RUnlock()

	if len(blockchain) == 0 {
		return models.Block{}
	}
	return blockchain[len(blockchain)-1]
}

// AddPendingTransaction adds a transaction to the pending pool
func AddPendingTransaction(tx models.Transaction) error {
	blockchainMutex.Lock()
	defer blockchainMutex.Unlock()

	pendingTransactions = append(pendingTransactions, tx)

	// Save to database
	if err := SavePendingTransaction(tx); err != nil {
		return err
	}

	log.Printf("Transaction %s added to pending pool", tx.Hash)
	return nil
}

// MineBlock mines a new block with pending transactions
func MineBlock(minerWalletID string) (models.Block, error) {
	blockchainMutex.Lock()
	defer blockchainMutex.Unlock()

	if len(pendingTransactions) == 0 {
		return models.Block{}, fmt.Errorf("no pending transactions to mine")
	}

	latestBlock := blockchain[len(blockchain)-1]

	newBlock := models.Block{
		Index:        latestBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: pendingTransactions,
		PreviousHash: latestBlock.Hash,
		Difficulty:   difficulty,
		MinedBy:      minerWalletID,
	}

	newBlock.MerkleRoot = calculateMerkleRoot(newBlock.Transactions)

	// Proof of Work
	log.Printf("Mining block %d with %d transactions...", newBlock.Index, len(newBlock.Transactions))
	newBlock = proofOfWork(newBlock)

	// Add block to chain
	blockchain = append(blockchain, newBlock)

	// Clear pending transactions
	pendingTransactions = []models.Transaction{}

	// Update transaction statuses
	for _, tx := range newBlock.Transactions {
		tx.Status = "confirmed"
		tx.BlockHash = newBlock.Hash
		if err := UpdateTransaction(tx); err != nil {
			log.Printf("Error updating transaction: %v", err)
		}
	}

	// Save block to database
	if err := SaveBlock(newBlock); err != nil {
		return models.Block{}, err
	}

	// Clear pending transactions from database
	if err := ClearPendingTransactions(); err != nil {
		log.Printf("Error clearing pending transactions: %v", err)
	}

	log.Printf("Block %d mined successfully with hash: %s", newBlock.Index, newBlock.Hash)

	// Log mining event
	LogSystemEvent("mining", fmt.Sprintf("Block %d mined by %s", newBlock.Index, minerWalletID), minerWalletID, "")

	return newBlock, nil
}

// proofOfWork performs the mining process
func proofOfWork(block models.Block) models.Block {
	target := strings.Repeat("0", difficulty)

	for {
		block.Hash = calculateBlockHash(block)

		if strings.HasPrefix(block.Hash, target) {
			break
		}

		block.Nonce++
	}

	return block
}

// calculateBlockHash calculates the hash of a block
func calculateBlockHash(block models.Block) string {
	transactionsJSON, _ := json.Marshal(block.Transactions)

	data := fmt.Sprintf("%d%s%s%s%d%s",
		block.Index,
		block.Timestamp.String(),
		string(transactionsJSON),
		block.PreviousHash,
		block.Nonce,
		block.MerkleRoot,
	)

	return crypto.HashSHA256(data)
}

// calculateMerkleRoot calculates the merkle root of transactions
func calculateMerkleRoot(transactions []models.Transaction) string {
	if len(transactions) == 0 {
		return ""
	}

	var hashes []string
	for _, tx := range transactions {
		hashes = append(hashes, tx.Hash)
	}

	for len(hashes) > 1 {
		if len(hashes)%2 != 0 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}

		var newHashes []string
		for i := 0; i < len(hashes); i += 2 {
			combined := hashes[i] + hashes[i+1]
			newHashes = append(newHashes, crypto.HashSHA256(combined))
		}
		hashes = newHashes
	}

	return hashes[0]
}

// ValidateChain validates the entire blockchain
func ValidateChain() bool {
	blockchainMutex.RLock()
	defer blockchainMutex.RUnlock()

	for i := 1; i < len(blockchain); i++ {
		currentBlock := blockchain[i]
		previousBlock := blockchain[i-1]

		// Check if previous hash matches
		if currentBlock.PreviousHash != previousBlock.Hash {
			log.Printf("Invalid previous hash at block %d", i)
			return false
		}

		// Recalculate hash
		calculatedHash := calculateBlockHash(currentBlock)
		if currentBlock.Hash != calculatedHash {
			log.Printf("Invalid hash at block %d", i)
			return false
		}

		// Check PoW
		target := strings.Repeat("0", currentBlock.Difficulty)
		if !strings.HasPrefix(currentBlock.Hash, target) {
			log.Printf("Invalid PoW at block %d", i)
			return false
		}
	}

	return true
}

// GetBlockchain returns the entire blockchain
func GetBlockchain() []models.Block {
	blockchainMutex.RLock()
	defer blockchainMutex.RUnlock()
	return blockchain
}

// GetPendingTransactionsFromMemory returns all pending transactions from memory
func GetPendingTransactionsFromMemory() []models.Transaction {
	blockchainMutex.RLock()
	defer blockchainMutex.RUnlock()
	return pendingTransactions
}

// GetBlockByHashFromMemory retrieves a block by its hash from memory
func GetBlockByHashFromMemory(hash string) *models.Block {
	blockchainMutex.RLock()
	defer blockchainMutex.RUnlock()

	for _, block := range blockchain {
		if block.Hash == hash {
			return &block
		}
	}
	return nil
}

// GetBlockByIndexFromMemory retrieves a block by its index from memory
func GetBlockByIndexFromMemory(index int64) *models.Block {
	blockchainMutex.RLock()
	defer blockchainMutex.RUnlock()

	if index < 0 || int(index) >= len(blockchain) {
		return nil
	}
	return &blockchain[index]
}
