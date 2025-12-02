package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	FullName      string    `bson:"fullName" json:"fullName"`
	Email         string    `bson:"email" json:"email"`
	Password      string    `bson:"password" json:"-"` // Hashed password, not sent in JSON
	CNIC          string    `bson:"cnic" json:"cnic"`
	WalletID      string    `bson:"walletId" json:"walletId"`
	PublicKey     string    `bson:"publicKey" json:"publicKey"`
	PrivateKey    string    `bson:"privateKey" json:"privateKey"` // Encrypted
	Beneficiaries []string  `bson:"beneficiaries" json:"beneficiaries"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`
	ZakatTracking ZakatInfo `bson:"zakatTracking" json:"zakatTracking"`
}

// ZakatInfo tracks zakat deductions for a user
type ZakatInfo struct {
	LastDeduction   time.Time `bson:"lastDeduction" json:"lastDeduction"`
	TotalDeducted   float64   `bson:"totalDeducted" json:"totalDeducted"`
	MonthlyDeducted float64   `bson:"monthlyDeducted" json:"monthlyDeducted"`
}

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	WalletID  string    `bson:"walletId" json:"walletId"`
	UserID    string    `bson:"userId" json:"userId"`
	PublicKey string    `bson:"publicKey" json:"publicKey"`
	Balance   float64   `bson:"balance" json:"balance"` // Cached balance
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	IsActive  bool      `bson:"isActive" json:"isActive"`
}

// UTXO represents an Unspent Transaction Output
type UTXO struct {
	ID              string    `bson:"_id,omitempty" json:"id"`
	TransactionHash string    `bson:"transactionHash" json:"transactionHash"`
	OutputIndex     int       `bson:"outputIndex" json:"outputIndex"`
	WalletID        string    `bson:"walletId" json:"walletId"`
	Amount          float64   `bson:"amount" json:"amount"`
	Spent           bool      `bson:"spent" json:"spent"`
	SpentInTxHash   string    `bson:"spentInTxHash,omitempty" json:"spentInTxHash,omitempty"`
	CreatedAt       time.Time `bson:"createdAt" json:"createdAt"`
	SpentAt         time.Time `bson:"spentAt,omitempty" json:"spentAt,omitempty"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash             string       `bson:"hash" json:"hash"`
	SenderWalletID   string       `bson:"senderWalletId" json:"senderWalletId"`
	ReceiverWalletID string       `bson:"receiverWalletId" json:"receiverWalletId"`
	Amount           float64      `bson:"amount" json:"amount"`
	Note             string       `bson:"note,omitempty" json:"note,omitempty"`
	Timestamp        time.Time    `bson:"timestamp" json:"timestamp"`
	SenderPublicKey  string       `bson:"senderPublicKey" json:"senderPublicKey"`
	Signature        string       `bson:"signature" json:"signature"`
	InputUTXOs       []string     `bson:"inputUtxos" json:"inputUtxos"`   // UTXO IDs being spent
	OutputUTXOs      []UTXOOutput `bson:"outputUtxos" json:"outputUtxos"` // New UTXOs created
	Type             string       `bson:"type" json:"type"`               // "transfer", "zakat_deduction", "mining_reward"
	Status           string       `bson:"status" json:"status"`           // "pending", "confirmed", "failed"
	BlockHash        string       `bson:"blockHash,omitempty" json:"blockHash,omitempty"`
}

// UTXOOutput represents a new UTXO created in a transaction
type UTXOOutput struct {
	WalletID string  `bson:"walletId" json:"walletId"`
	Amount   float64 `bson:"amount" json:"amount"`
}

// Block represents a block in the blockchain
type Block struct {
	Index        int64         `bson:"index" json:"index"`
	Timestamp    time.Time     `bson:"timestamp" json:"timestamp"`
	Transactions []Transaction `bson:"transactions" json:"transactions"`
	PreviousHash string        `bson:"previousHash" json:"previousHash"`
	Nonce        int64         `bson:"nonce" json:"nonce"`
	Hash         string        `bson:"hash" json:"hash"`
	MerkleRoot   string        `bson:"merkleRoot" json:"merkleRoot"`
	Difficulty   int           `bson:"difficulty" json:"difficulty"`
	MinedBy      string        `bson:"minedBy,omitempty" json:"minedBy,omitempty"`
}

// SystemLog represents system-wide logs
type SystemLog struct {
	ID        string                 `bson:"_id,omitempty" json:"id"`
	Type      string                 `bson:"type" json:"type"` // "login", "failed_login", "mining", "validation_failure", etc.
	Message   string                 `bson:"message" json:"message"`
	UserID    string                 `bson:"userId,omitempty" json:"userId,omitempty"`
	IPAddress string                 `bson:"ipAddress,omitempty" json:"ipAddress,omitempty"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// TransactionLog represents transaction-specific logs
type TransactionLog struct {
	ID              string    `bson:"_id,omitempty" json:"id"`
	TransactionHash string    `bson:"transactionHash" json:"transactionHash"`
	Action          string    `bson:"action" json:"action"` // "sent", "received", "mined", "zakat_deducted"
	UserID          string    `bson:"userId" json:"userId"`
	WalletID        string    `bson:"walletId" json:"walletId"`
	Amount          float64   `bson:"amount" json:"amount"`
	BlockHash       string    `bson:"blockHash,omitempty" json:"blockHash,omitempty"`
	Status          string    `bson:"status" json:"status"`
	IPAddress       string    `bson:"ipAddress,omitempty" json:"ipAddress,omitempty"`
	Note            string    `bson:"note,omitempty" json:"note,omitempty"`
	Timestamp       time.Time `bson:"timestamp" json:"timestamp"`
}

// ZakatDeduction represents a monthly zakat deduction
type ZakatDeduction struct {
	ID              string    `bson:"_id,omitempty" json:"id"`
	UserID          string    `bson:"userId" json:"userId"`
	WalletID        string    `bson:"walletId" json:"walletId"`
	Amount          float64   `bson:"amount" json:"amount"`
	BalanceBefore   float64   `bson:"balanceBefore" json:"balanceBefore"`
	BalanceAfter    float64   `bson:"balanceAfter" json:"balanceAfter"`
	TransactionHash string    `bson:"transactionHash" json:"transactionHash"`
	BlockHash       string    `bson:"blockHash,omitempty" json:"blockHash,omitempty"`
	Month           string    `bson:"month" json:"month"` // "2024-12"
	DeductedAt      time.Time `bson:"deductedAt" json:"deductedAt"`
	Status          string    `bson:"status" json:"status"` // "pending", "completed", "failed"
}

// PendingTransaction represents transactions waiting to be mined
type PendingTransaction struct {
	Transaction Transaction `bson:"transaction" json:"transaction"`
	ReceivedAt  time.Time   `bson:"receivedAt" json:"receivedAt"`
}
