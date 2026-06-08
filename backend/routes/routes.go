package routes

import (
	"vapor_auror_backend/controllers"
	"vapor_auror_backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 32 << 20 // 32 MB
	r.Static("/media/images", "./assets/images")

	api := r.Group("/api")
	{
		// ==========================================
		// Phase 2: Auth Routes (Public)
		// ==========================================
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/logout", middleware.RequireAuth(), controllers.Logout)
		}

		// ==========================================
		// Phase 1: Game Routes (Public)
		// ==========================================
		api.GET("/games", controllers.GetGames)
		api.GET("/games/:id", controllers.GetGameByID)
		api.GET("/games/:id/reviews", controllers.GetReviews)

		// ==========================================
		// Phase 3: Protected Routes (Require JWT Token)
		// ==========================================
		protected := api.Group("/protected").Use(middleware.RequireAuth())
		{
			// 1. Shopping Cart
			protected.GET("/cart", controllers.GetCart)
			protected.POST("/cart", controllers.AddToCart)
			protected.DELETE("/cart/:game_id", controllers.RemoveFromCart)

			// 2. Checkout & Transactions
			protected.POST("/checkout", controllers.Checkout)
			protected.GET("/transactions", controllers.GetTransactions)
			protected.GET("/refunds", controllers.GetMyRefunds)

			// 3. Library & Wishlist
			protected.GET("/library", controllers.GetLibrary)
			protected.GET("/wishlist", controllers.GetWishlist)
			protected.POST("/wishlist", controllers.AddToWishlist)
			protected.DELETE("/wishlist/:game_id", controllers.RemoveFromWishlist)
			protected.GET("/library/:game_id/play", controllers.PlayGame)
			protected.GET("/library/:game_id/download", controllers.DownloadGame)
		}

		// ==========================================
		// Phase 4: Developer Routes (Require JWT + DEVELOPER Role)
		// ==========================================
		developer := api.Group("/developer").Use(middleware.RequireAuth(), middleware.RequireRole("DEVELOPER"))
		{
			developer.GET("/games", controllers.GetDeveloperGames)
			developer.POST("/games", controllers.UploadGame)
			developer.PUT("/games/:id", controllers.UpdateGame)
			developer.PUT("/games/:id/publish", controllers.PublishGame)
			developer.DELETE("/games/:id", controllers.DeleteGame)
			developer.POST("/games/:id/media", controllers.UploadMedia)
			developer.DELETE("/games/:id/media/:media_id", controllers.DeleteMedia)
			developer.GET("/games/:id/stats", controllers.GetGameStats)
			developer.POST("/tags", controllers.CreateTag)
			developer.POST("/games/:id/tags", controllers.AddTagToGame)
			developer.DELETE("/games/:id/tags/:tag_id", controllers.RemoveTagFromGame)
		}

		api.GET("/tags", controllers.GetTags)

		// ==========================================
		// Phase 5: Complete Remaining API Groups
		// ==========================================

		// Protected User Routes
		users := api.Group("/users").Use(middleware.RequireAuth())
		{
			users.PUT("/profile", controllers.UpdateProfile)
		}

		// ADMIN Zone
		admin := api.Group("/admin").Use(middleware.RequireAuth(), middleware.RequireRole("ADMIN"))
		{
			admin.GET("/users", controllers.GetUsers)
			admin.PUT("/users/:id/suspend", controllers.SuspendUser)
			admin.DELETE("/users/:id", controllers.DeleteUser)
			admin.PUT("/users/:id/role", controllers.ChangeUserRole)
			admin.DELETE("/games/:id", controllers.AdminDeleteGame)
		}

		// CSR Zone
		csr := api.Group("/csr").Use(middleware.RequireAuth(), middleware.RequireRole("CSR"))
		{
			csr.GET("/refunds", controllers.GetRefundRequests)
			csr.PUT("/refunds/:id", controllers.ProcessRefund)
		}

		// SOCIAL Zone
		social := api.Group("/social").Use(middleware.RequireAuth())
		{
			social.POST("/games/:id/reviews", controllers.PostReview)
			social.POST("/reviews/:review_id/replies", controllers.ReplyToReview)
			social.DELETE("/reviews/replies/:reply_id", controllers.DeleteReviewReply)
			social.POST("/refunds", controllers.ApplyRefund)

			// Friends
			social.GET("/friends", controllers.GetFriends)
			social.GET("/friends/requests", controllers.GetFriendRequests)
			social.POST("/friends/request", controllers.SendFriendRequest)
			social.PUT("/friends/request/:id/accept", controllers.AcceptFriendRequest)
			social.PUT("/friends/request/:id/decline", controllers.DeclineFriendRequest)
			social.DELETE("/friends/request/:id", controllers.RevokeFriendRequest)

			// Blacklist
			social.GET("/blacklist", controllers.GetBlacklist)
			social.POST("/blacklist", controllers.AddBlacklist)
			social.DELETE("/blacklist/:user_id", controllers.RemoveBlacklist)

			// Messages
			social.POST("/messages", controllers.SendMessage)
			social.GET("/messages/:user_id", controllers.GetMessages)
		}
	}

	return r
}
