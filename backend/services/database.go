package services

import (
	"backend/config"
	"backend/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection names
const (
	UsersCollection               = "users"
	WalletsCollection             = "wallets"
	UTXOsCollection               = "utxos"
	TransactionsCollection        = "transactions"
	PendingTransactionsCollection = "pendingTransactions"
	BlocksCollection              = "blocks"
	ZakatDeductionsCollection     = "zakatDeductions"
	SystemLogsCollection          = "systemLogs"
	TransactionLogsCollection     = "transactionLogs"
)

// User operations

// SaveUser saves a user to MongoDB
func SaveUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(UsersCollection)

	// Create unique index on email
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, _ = collection.Indexes().CreateOne(ctx, indexModel)

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(UsersCollection)

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(UsersCollection)

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates a user
func UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	return SaveUser(user)
}

// GetAllUsers retrieves all users
func GetAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(UsersCollection)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// Wallet operations

// SaveWallet saves a wallet to MongoDB
func SaveWallet(wallet *models.Wallet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(WalletsCollection)

	// Create unique index on walletId
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "walletId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, _ = collection.Indexes().CreateOne(ctx, indexModel)

	filter := bson.M{"walletId": wallet.WalletID}
	update := bson.M{"$set": wallet}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetWalletByID retrieves a wallet by ID
func GetWalletByID(walletID string) (*models.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(WalletsCollection)

	var wallet models.Wallet
	err := collection.FindOne(ctx, bson.M{"walletId": walletID}).Decode(&wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// UpdateWallet updates a wallet
func UpdateWallet(wallet *models.Wallet) error {
	wallet.UpdatedAt = time.Now()
	return SaveWallet(wallet)
}

// UTXO operations

// SaveUTXO saves a UTXO to MongoDB
func SaveUTXO(utxo *models.UTXO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(UTXOsCollection)

	filter := bson.M{"_id": utxo.ID}
	update := bson.M{"$set": utxo}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetUTXOByID retrieves a UTXO by ID
func GetUTXOByID(utxoID string) (*models.UTXO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(UTXOsCollection)

	var utxo models.UTXO
	err := collection.FindOne(ctx, bson.M{"_id": utxoID}).Decode(&utxo)
	if err != nil {
		return nil, err
	}

	return &utxo, nil
}

// GetUnspentUTXOs retrieves all unspent UTXOs for a wallet
func GetUnspentUTXOs(walletID string) ([]models.UTXO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(UTXOsCollection)

	filter := bson.M{
		"walletId": walletID,
		"spent":    false,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var utxos []models.UTXO
	if err = cursor.All(ctx, &utxos); err != nil {
		return nil, err
	}

	return utxos, nil
}

// GetAllUTXOs retrieves all UTXOs
func GetAllUTXOs() ([]models.UTXO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(UTXOsCollection)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var utxos []models.UTXO
	if err = cursor.All(ctx, &utxos); err != nil {
		return nil, err
	}

	return utxos, nil
}

// UpdateUTXO updates a UTXO
func UpdateUTXO(utxo *models.UTXO) error {
	return SaveUTXO(utxo)
}

// Transaction operations

// SaveTransaction saves a transaction to MongoDB
func SaveTransaction(tx *models.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(TransactionsCollection)

	// Create unique index on hash
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "hash", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, _ = collection.Indexes().CreateOne(ctx, indexModel)

	filter := bson.M{"hash": tx.Hash}
	update := bson.M{"$set": tx}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetTransactionFromDB retrieves a transaction by hash
func GetTransactionFromDB(hash string) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(TransactionsCollection)

	var tx models.Transaction
	err := collection.FindOne(ctx, bson.M{"hash": hash}).Decode(&tx)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

// UpdateTransaction updates a transaction
func UpdateTransaction(tx models.Transaction) error {
	return SaveTransaction(&tx)
}

// GetWalletTransactions retrieves all transactions for a wallet
func GetWalletTransactions(walletID string) ([]models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(TransactionsCollection)

	filter := bson.M{
		"$or": []bson.M{
			{"senderWalletId": walletID},
			{"receiverWalletId": walletID},
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// Pending Transaction operations

// SavePendingTransaction saves a pending transaction
func SavePendingTransaction(tx models.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(PendingTransactionsCollection)

	pendingTx := models.PendingTransaction{
		Transaction: tx,
		ReceivedAt:  time.Now(),
	}

	filter := bson.M{"transaction.hash": tx.Hash}
	update := bson.M{"$set": pendingTx}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetPendingTransactions retrieves all pending transactions
func GetPendingTransactions() ([]models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(PendingTransactionsCollection)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pendingTxs []models.PendingTransaction
	if err = cursor.All(ctx, &pendingTxs); err != nil {
		return nil, err
	}

	transactions := make([]models.Transaction, len(pendingTxs))
	for i, pt := range pendingTxs {
		transactions[i] = pt.Transaction
	}

	return transactions, nil
}

// ClearPendingTransactions removes all pending transactions
func ClearPendingTransactions() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(PendingTransactionsCollection)

	_, err := collection.DeleteMany(ctx, bson.M{})
	return err
}

// RemovePendingTransaction removes a specific pending transaction
func RemovePendingTransaction(txHash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(PendingTransactionsCollection)

	_, err := collection.DeleteOne(ctx, bson.M{"transaction.hash": txHash})
	return err
}

// Block operations

// SaveBlock saves a block to MongoDB
func SaveBlock(block models.Block) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(BlocksCollection)

	// Create unique index on hash
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "hash", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, _ = collection.Indexes().CreateOne(ctx, indexModel)

	filter := bson.M{"hash": block.Hash}
	update := bson.M{"$set": block}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetAllBlocks retrieves all blocks
func GetAllBlocks() ([]models.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(BlocksCollection)

	opts := options.Find().SetSort(bson.D{{Key: "index", Value: 1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blocks []models.Block
	if err = cursor.All(ctx, &blocks); err != nil {
		return nil, err
	}

	return blocks, nil
}

// GetBlockByHash retrieves a block by hash
func GetBlockByHash(hash string) (*models.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(BlocksCollection)

	var block models.Block
	err := collection.FindOne(ctx, bson.M{"hash": hash}).Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetBlockByIndex retrieves a block by index
func GetBlockByIndex(index int64) (*models.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(BlocksCollection)

	var block models.Block
	err := collection.FindOne(ctx, bson.M{"index": index}).Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// Zakat operations

// SaveZakatDeduction saves a zakat deduction
func SaveZakatDeduction(deduction *models.ZakatDeduction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(ZakatDeductionsCollection)

	filter := bson.M{"_id": deduction.ID}
	update := bson.M{"$set": deduction}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetUserZakatDeductions retrieves all zakat deductions for a user
func GetUserZakatDeductions(userID string) ([]models.ZakatDeduction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(ZakatDeductionsCollection)

	opts := options.Find().SetSort(bson.D{{Key: "deductedAt", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deductions []models.ZakatDeduction
	if err = cursor.All(ctx, &deductions); err != nil {
		return nil, err
	}

	return deductions, nil
}

// GetAllZakatDeductions retrieves all zakat deductions
func GetAllZakatDeductions() ([]models.ZakatDeduction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(ZakatDeductionsCollection)

	opts := options.Find().SetSort(bson.D{{Key: "deductedAt", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deductions []models.ZakatDeduction
	if err = cursor.All(ctx, &deductions); err != nil {
		return nil, err
	}

	return deductions, nil
}

// Logging operations

// SaveSystemLog saves a system log
func SaveSystemLog(log *models.SystemLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(SystemLogsCollection)

	filter := bson.M{"_id": log.ID}
	update := bson.M{"$set": log}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetSystemLogs retrieves system logs with optional filters
func GetSystemLogs(logType string, limit int) ([]models.SystemLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(SystemLogsCollection)

	filter := bson.M{}
	if logType != "" {
		filter["type"] = logType
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.SystemLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// SaveTransactionLog saves a transaction log
func SaveTransactionLog(log *models.TransactionLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(TransactionLogsCollection)

	filter := bson.M{"_id": log.ID}
	update := bson.M{"$set": log}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetUserTransactionLogs retrieves transaction logs for a user
func GetUserTransactionLogs(userID string, limit int) ([]models.TransactionLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.GetCollection(TransactionLogsCollection)

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := collection.Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.TransactionLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}
