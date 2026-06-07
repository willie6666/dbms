package controllers

import (
	"net/http"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
)

// GetReviews handles GET /api/games/:id/reviews
func GetReviews(c *gin.Context) {
	gameID := c.Param("id")
	var reviews []models.Review

	// Preload the User to get Username
	if err := database.DB.Preload("User").Where("game_id = ?", gameID).Order("created_at desc").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	type ReviewWithReplies struct {
		models.Review
		Replies []models.ReviewReply `json:"replies"`
	}

	var fullReviews []ReviewWithReplies
	for _, r := range reviews {
		var replies []models.ReviewReply
		database.DB.Preload("User").Where("review_id = ?", r.ReviewID).Order("created_at asc").Find(&replies)
		fullReviews = append(fullReviews, ReviewWithReplies{Review: r, Replies: replies})
	}

	c.JSON(http.StatusOK, fullReviews)
}

// PostReview handles POST /api/social/games/:id/reviews
func PostReview(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	gameID := c.Param("id")

	var input struct {
		Attitude string `json:"attitude" binding:"required"` // POSITIVE or NEGATIVE
		Content  string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// VERIFY OWNERSHIP: Only players who own the game can leave a review
	var license models.GameLicense
	if err := database.DB.Where("user_id = ? AND game_id = ? AND status = ?", userID, game.GameID, "ACTIVE").First(&license).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You must own the game to leave a review"})
		return
	}

	review := models.Review{
		GameID:   game.GameID,
		UserID:   userID,
		Attitude: input.Attitude,
		Content:  input.Content,
	}

	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Review posted successfully"})
}

// ApplyRefund handles POST /api/social/refunds
func ApplyRefund(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var input struct {
		TransactionItemID uint   `json:"transaction_item_id" binding:"required"`
		Reason            string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Verify Ownership: The user must own the transaction item
	var license models.GameLicense
	if err := database.DB.Where("user_id = ? AND transaction_item_id = ?", userID, input.TransactionItemID).First(&license).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Transaction item not found in your library"})
		return
	}

	// 2. Prevent Duplicate Refunds
	var existing models.RefundRequest
	if err := database.DB.Where("transaction_item_id = ?", input.TransactionItemID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A refund request already exists for this item"})
		return
	}

	req := models.RefundRequest{
		BuyerID:           userID,
		TransactionItemID: input.TransactionItemID,
		Reason:            input.Reason,
	}

	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit refund request"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Refund request submitted. A CSR will review it shortly."})
}

// GetFriends handles GET /api/social/friends
func GetFriends(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var friends []models.Friendship
	database.DB.Where("(sender_id = ? OR receiver_id = ?) AND status = ?", userID, userID, "ACCEPTED").Find(&friends)

	c.JSON(http.StatusOK, gin.H{"data": friends})
}

// SendFriendRequest handles POST /api/social/friends/request
func SendFriendRequest(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var input struct {
		ReceiverID uint `json:"receiver_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	friend := models.Friendship{
		SenderID:   userID,
		ReceiverID: input.ReceiverID,
		Status:     "PENDING",
	}

	database.DB.Create(&friend)
	c.JSON(http.StatusCreated, gin.H{"message": "Friend request sent"})
}

// SendMessage handles POST /api/social/messages
func SendMessage(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var input struct {
		ReceiverID uint   `json:"receiver_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := models.Message{
		SenderID:   userID,
		ReceiverID: input.ReceiverID,
		Content:    input.Content,
	}
	database.DB.Create(&msg)
	c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
}

// GetMessages handles GET /api/social/messages/:user_id
func GetMessages(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	myID := uint(userIDFloat.(float64))
	otherID := c.Param("user_id")

	var messages []models.Message
	database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", myID, otherID, otherID, myID).Order("sent_at asc").Find(&messages)

	c.JSON(http.StatusOK, gin.H{"data": messages})
}

// ReplyToReview handles POST /api/social/reviews/:review_id/replies
func ReplyToReview(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	reviewID := c.Param("review_id")

	var input struct {
		ParentReplyID *uint  `json:"parent_reply_id"` // Optional
		Content       string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var review models.Review
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	reply := models.ReviewReply{
		ReviewID:      review.ReviewID,
		UserID:        userID,
		ParentReplyID: input.ParentReplyID,
		Content:       input.Content,
	}

	if err := database.DB.Create(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post reply"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Reply posted successfully", "data": reply})
}

// DeleteReviewReply handles DELETE /api/social/reviews/replies/:reply_id
func DeleteReviewReply(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	replyID := c.Param("reply_id")

	var reply models.ReviewReply
	if err := database.DB.First(&reply, replyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reply not found"})
		return
	}

	if reply.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only delete your own replies"})
		return
	}

	if err := database.DB.Delete(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reply"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reply deleted successfully"})
}

// AcceptFriendRequest handles PUT /api/social/friends/request/:id/accept
func AcceptFriendRequest(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	reqID := c.Param("id")

	var friend models.Friendship
	if err := database.DB.First(&friend, reqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
		return
	}

	if friend.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You are not the receiver"})
		return
	}

	friend.Status = "ACCEPTED"
	database.DB.Save(&friend)
	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted"})
}

// DeclineFriendRequest handles PUT /api/social/friends/request/:id/decline
func DeclineFriendRequest(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	reqID := c.Param("id")

	var friend models.Friendship
	if err := database.DB.First(&friend, reqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
		return
	}

	if friend.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You are not the receiver"})
		return
	}

	friend.Status = "DECLINED"
	database.DB.Save(&friend)
	c.JSON(http.StatusOK, gin.H{"message": "Friend request declined"})
}

// RevokeFriendRequest handles DELETE /api/social/friends/request/:id
func RevokeFriendRequest(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	reqID := c.Param("id")

	var friend models.Friendship
	if err := database.DB.First(&friend, reqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
		return
	}

	if friend.SenderID != userID && friend.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	database.DB.Delete(&friend)
	c.JSON(http.StatusOK, gin.H{"message": "Friend request revoked / removed"})
}

// AddBlacklist handles POST /api/social/blacklist
func AddBlacklist(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var input struct {
		BlockedID uint `json:"blocked_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blacklist := models.Blacklist{
		BlockerID: userID,
		BlockedID: input.BlockedID,
	}

	database.DB.Create(&blacklist)
	c.JSON(http.StatusCreated, gin.H{"message": "User added to blacklist"})
}

// RemoveBlacklist handles DELETE /api/social/blacklist/:user_id
func RemoveBlacklist(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	blockedID := c.Param("user_id")

	if err := database.DB.Where("blocker_id = ? AND blocked_id = ?", userID, blockedID).Delete(&models.Blacklist{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from blacklist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from blacklist"})
}

// GetFriendRequests handles GET /api/social/friends/requests
func GetFriendRequests(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var requests []models.Friendship
	database.DB.Where("receiver_id = ? AND status = ?", userID, "PENDING").Find(&requests)

	c.JSON(http.StatusOK, gin.H{"data": requests})
}

// GetBlacklist handles GET /api/social/blacklist
func GetBlacklist(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var blacklisted []models.Blacklist
	database.DB.Where("blocker_id = ?", userID).Find(&blacklisted)

	c.JSON(http.StatusOK, gin.H{"data": blacklisted})
}
