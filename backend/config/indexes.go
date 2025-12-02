package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes creates all necessary database indexes for optimal performance
func CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Users collection indexes
	usersCollection := GetCollection("users")
	usersIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "walletId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "cnic", Value: 1}},
		},
	}
	if _, err := usersCollection.Indexes().CreateMany(ctx, usersIndexes); err != nil {
		log.Printf("Warning: Failed to create users indexes: %v", err)
	}

	// Wallets collection indexes
	walletsCollection := GetCollection("wallets")
	walletsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "walletId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "isActive", Value: 1}},
		},
	}
	if _, err := walletsCollection.Indexes().CreateMany(ctx, walletsIndexes); err != nil {
		log.Printf("Warning: Failed to create wallets indexes: %v", err)
	}

	// UTXOs collection indexes
	utxosCollection := GetCollection("utxos")
	utxosIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "walletId", Value: 1}, {Key: "spent", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "transactionHash", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "spent", Value: 1}},
		},
	}
	if _, err := utxosCollection.Indexes().CreateMany(ctx, utxosIndexes); err != nil {
		log.Printf("Warning: Failed to create utxos indexes: %v", err)
	}

	// Transactions collection indexes
	transactionsCollection := GetCollection("transactions")
	transactionsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "hash", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "senderWalletId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "receiverWalletId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
	}
	if _, err := transactionsCollection.Indexes().CreateMany(ctx, transactionsIndexes); err != nil {
		log.Printf("Warning: Failed to create transactions indexes: %v", err)
	}

	// Blocks collection indexes
	blocksCollection := GetCollection("blocks")
	blocksIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "hash", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "index", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "minerWalletId", Value: 1}},
		},
	}
	if _, err := blocksCollection.Indexes().CreateMany(ctx, blocksIndexes); err != nil {
		log.Printf("Warning: Failed to create blocks indexes: %v", err)
	}

	// Zakat deductions collection indexes
	zakatCollection := GetCollection("zakatDeductions")
	zakatIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "walletId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "month", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
	}
	if _, err := zakatCollection.Indexes().CreateMany(ctx, zakatIndexes); err != nil {
		log.Printf("Warning: Failed to create zakat indexes: %v", err)
	}

	// System logs collection indexes
	systemLogsCollection := GetCollection("systemLogs")
	systemLogsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "eventType", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
	}
	if _, err := systemLogsCollection.Indexes().CreateMany(ctx, systemLogsIndexes); err != nil {
		log.Printf("Warning: Failed to create system logs indexes: %v", err)
	}

	// Transaction logs collection indexes
	transactionLogsCollection := GetCollection("transactionLogs")
	transactionLogsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "transactionHash", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "fromWallet", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "toWallet", Value: 1}},
		},
	}
	if _, err := transactionLogsCollection.Indexes().CreateMany(ctx, transactionLogsIndexes); err != nil {
		log.Printf("Warning: Failed to create transaction logs indexes: %v", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
