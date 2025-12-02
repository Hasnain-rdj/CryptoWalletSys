package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// Create rate limiters
	authLimiter := middleware.NewRateLimiter(20)        // 20 requests per minute for auth
	apiLimiter := middleware.NewRateLimiter(100)        // 100 requests per minute for general API
	transactionLimiter := middleware.NewRateLimiter(30) // 30 requests per minute for transactions

	api := r.Group("/api")
	{
		// Public routes (no authentication required)
		public := api.Group("/")
		public.Use(apiLimiter.RateLimit())
		{
			// Authentication with stricter rate limiting
			auth := public.Group("/")
			auth.Use(authLimiter.RateLimit())
			{
				auth.POST("/register", handlers.Register)
				auth.POST("/login", handlers.Login)
				auth.POST("/otp/generate", handlers.GenerateOTP)
				auth.POST("/otp/verify", handlers.VerifyOTP)
			}

			// Email verification
			public.GET("/otp/check", handlers.CheckEmailVerification)

			// Blockchain explorer (public)
			public.GET("/blockchain", handlers.GetBlockchain)
			public.GET("/blockchain/stats", handlers.GetBlockchainStats)
			public.GET("/blockchain/validate", handlers.ValidateBlockchain)
			public.GET("/block/hash/:hash", handlers.GetBlockByHash)
			public.GET("/block/index/:index", handlers.GetBlockByIndex)
			public.GET("/block/latest", handlers.GetLatestBlock)

			// Wallet validation (public)
			public.GET("/wallet/validate/:walletId", handlers.ValidateWalletID)

			// Transaction (public read)
			public.GET("/transaction/:hash", handlers.GetTransactionByHash)
			public.GET("/transactions/pending", handlers.GetPendingTransactions)
		}

		// Protected routes (authentication required)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		protected.Use(apiLimiter.RateLimit())
		{
			// User profile
			protected.GET("/profile", handlers.GetProfile)
			protected.PUT("/profile", handlers.UpdateProfile)

			// Beneficiaries
			protected.POST("/beneficiary", handlers.AddBeneficiary)
			protected.DELETE("/beneficiary/:walletId", handlers.RemoveBeneficiary)

			// Wallet
			protected.GET("/wallet", handlers.GetWallet)
			protected.GET("/balance", handlers.GetBalance)
			protected.GET("/wallet/utxos", handlers.GetWalletUTXOs)

			// Transactions with separate rate limiter
			transactions := protected.Group("/")
			transactions.Use(transactionLimiter.RateLimit())
			{
				transactions.POST("/transaction", handlers.CreateTransaction)
				transactions.POST("/mine", handlers.MineBlockManual)
			}

			// Transaction history (read-only, less restrictive)
			protected.GET("/transactions", handlers.GetTransactionHistory)

			// Zakat
			protected.GET("/zakat/history", handlers.GetZakatHistory)
			protected.GET("/zakat/summary", handlers.GetZakatSummary)
			protected.POST("/zakat/deduct", handlers.TriggerZakatDeduction)

			// Logs
			protected.GET("/logs/system", handlers.GetSystemLogs)
			protected.GET("/logs/transactions", handlers.GetTransactionLogs)
			protected.GET("/logs/transactions/all", handlers.GetAllTransactionLogs)

			// Reports
			protected.GET("/reports", handlers.GetReports)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "blockchain-wallet-api",
		})
	})
}
