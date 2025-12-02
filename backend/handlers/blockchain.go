package handlers

import (
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBlockchain returns the entire blockchain
func GetBlockchain(c *gin.Context) {
	blockchain := services.GetBlockchain()

	c.JSON(http.StatusOK, gin.H{
		"blockchain": blockchain,
		"length":     len(blockchain),
	})
}

// GetBlockByHash returns a specific block by hash
func GetBlockByHash(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Block hash required"})
		return
	}

	block, err := services.GetBlockByHash(hash)
	if err != nil || block == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"block": block,
	})
}

// GetBlockByIndex returns a specific block by index
func GetBlockByIndex(c *gin.Context) {
	indexStr := c.Param("index")
	if indexStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Block index required"})
		return
	}

	index, err := strconv.ParseInt(indexStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block index"})
		return
	}

	block, err := services.GetBlockByIndex(index)
	if err != nil || block == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"block": block,
	})
}

// GetLatestBlock returns the latest block
func GetLatestBlock(c *gin.Context) {
	block := services.GetLatestBlock()

	c.JSON(http.StatusOK, gin.H{
		"block": block,
	})
}

// ValidateBlockchain validates the entire blockchain
func ValidateBlockchain(c *gin.Context) {
	isValid := services.ValidateChain()

	c.JSON(http.StatusOK, gin.H{
		"valid": isValid,
	})
}

// MineBlockManual manually triggers block mining
func MineBlockManual(c *gin.Context) {
	// This would typically be automated, but provided for testing
	walletID := c.Query("walletId")
	if walletID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet ID required"})
		return
	}

	block, err := services.MineBlock(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Block mined successfully",
		"block":   block,
	})
}

// GetBlockchainStats returns blockchain statistics
func GetBlockchainStats(c *gin.Context) {
	blockchain := services.GetBlockchain()
	pendingTxs, _ := services.GetPendingTransactions()

	totalTransactions := 0
	for _, block := range blockchain {
		totalTransactions += len(block.Transactions)
	}

	totalSupply, _ := services.GetTotalSupply()

	stats := map[string]interface{}{
		"totalBlocks":         len(blockchain),
		"totalTransactions":   totalTransactions,
		"pendingTransactions": len(pendingTxs),
		"totalSupply":         totalSupply,
		"latestBlock":         services.GetLatestBlock(),
	}

	c.JSON(http.StatusOK, stats)
}
