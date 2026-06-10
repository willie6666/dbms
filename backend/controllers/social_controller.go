package controllers

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
)

type socialUserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type reviewAuthorDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func getSocialUser(id uint) socialUserDTO {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return socialUserDTO{ID: id, Username: ""}
	}
	return socialUserDTO{ID: user.UserID, Username: user.Username}
}

func getReviewAuthor(user models.User, fallbackID uint) reviewAuthorDTO {
	if user.UserID == 0 || user.Permission == "DELETED" || user.Role == "NULL" {
		return reviewAuthorDTO{ID: fallbackID, Username: "已刪除玩家"}
	}
	return reviewAuthorDTO{ID: user.UserID, Username: user.Username}
}

func getReviewAuthorByID(userID uint) reviewAuthorDTO {
	var user models.User
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return reviewAuthorDTO{ID: userID, Username: "已刪除玩家"}
	}
	return getReviewAuthor(user, userID)
}

func parseRoleFromContent(content string) (role string, cleanContent string) {
	if strings.HasPrefix(content, "[ROLE:ADMIN]") {
		return "ADMIN", strings.TrimPrefix(content, "[ROLE:ADMIN]")
	}
	if strings.HasPrefix(content, "[ROLE:CSR]") {
		return "CSR", strings.TrimPrefix(content, "[ROLE:CSR]")
	}
	if strings.HasPrefix(content, "[ROLE:AUTHOR]") {
		return "AUTHOR", strings.TrimPrefix(content, "[ROLE:AUTHOR]")
	}
	return "USERS", content
}

func buildRoleContent(content string, role string) string {
	if role == "ADMIN" || role == "CSR" || role == "AUTHOR" {
		return "[ROLE:" + role + "]" + content
	}
	return content
}

// GetReviews handles GET /api/games/:id/reviews
func GetReviews(c *gin.Context) {
	gameID := c.Param("id")
	var reviews []models.Review

	// Preload the User to get Username
	if err := database.DB.Where("game_id = ?", gameID).Order("created_at desc").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	type ReviewReplyDTO struct {
		ReviewReplyID uint            `json:"review_reply_id"`
		ReviewID      uint            `json:"review_id"`
		UserID        uint            `json:"user_id"`
		Author        reviewAuthorDTO `json:"user"`
		AuthorName    string          `json:"author_username"`
		Content       string          `json:"content"`
		PostedAsRole  string          `json:"posted_as_role"`
		CreatedAt     string          `json:"created_at"`
	}
	type ReviewDTO struct {
		ReviewID   uint             `json:"review_id"`
		GameID     uint             `json:"game_id"`
		UserID     uint             `json:"user_id"`
		Author     reviewAuthorDTO  `json:"user"`
		AuthorName string           `json:"author_username"`
		Content    string           `json:"content"`
		Attitude   string           `json:"attitude"`
		Status     string           `json:"status"`
		PostedAsRole string         `json:"posted_as_role"`
		CreatedAt  string           `json:"created_at"`
		Replies    []ReviewReplyDTO `json:"replies"`
	}

	var fullReviews []ReviewDTO
	for _, r := range reviews {
		var replies []models.ReviewReply
		database.DB.Where("review_id = ?", r.ReviewID).Order("created_at asc").Find(&replies)

		replyDTOs := make([]ReviewReplyDTO, 0, len(replies))
		for _, reply := range replies {
			author := getReviewAuthorByID(reply.UserID)
			replyRole, replyContent := parseRoleFromContent(reply.Content)
			replyDTOs = append(replyDTOs, ReviewReplyDTO{
				ReviewReplyID: reply.ReviewReplyID,
				ReviewID:      reply.ReviewID,
				UserID:        reply.UserID,
				Author:        author,
				AuthorName:    author.Username,
				Content:       replyContent,
				PostedAsRole:  replyRole,
				CreatedAt:     reply.CreatedAt.Format(time.RFC3339),
			})
		}

		author := getReviewAuthorByID(r.UserID)
		revRole, revContent := parseRoleFromContent(r.Content)
		fullReviews = append(fullReviews, ReviewDTO{
			ReviewID:   r.ReviewID,
			GameID:     r.GameID,
			UserID:     r.UserID,
			Author:     author,
			AuthorName: author.Username,
			Content:    revContent,
			Attitude:   r.Attitude,
			Status:     r.Status,
			PostedAsRole: revRole,
			CreatedAt:  r.CreatedAt.Format(time.RFC3339),
			Replies:    replyDTOs,
		})
	}

	c.JSON(http.StatusOK, fullReviews)
}

// PostReview handles POST /api/social/games/:id/reviews
func PostReview(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	gameID := c.Param("id")

	var input struct {
		Attitude   string `json:"attitude" binding:"required"` // POSITIVE or NEGATIVE
		Content    string `json:"content" binding:"required"`
		PostAsRole string `json:"post_as_role"`
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
	if game.Status == "TAKEN_DOWN" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot review a game that has been taken down"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	canBypass := false
	if input.PostAsRole == "ADMIN" && user.Role == "ADMIN" {
		canBypass = true
	} else if input.PostAsRole == "CSR" && user.Role == "CSR" {
		canBypass = true
	} else if input.PostAsRole == "AUTHOR" && user.Role == "DEVELOPER" && game.DeveloperID == userID {
		canBypass = true
	} else {
		input.PostAsRole = "USERS"
	}

	if !canBypass {
		// VERIFY OWNERSHIP: Only players who own the game can leave a review
		var license models.GameLicense
		if err := database.DB.Where("user_id = ? AND game_id = ? AND status = ?", userID, game.GameID, "ACTIVE").First(&license).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You must own the game to leave a review"})
			return
		}
	}

	finalContent := buildRoleContent(input.Content, input.PostAsRole)

	review := models.Review{
		GameID:   game.GameID,
		UserID:   userID,
		Attitude: input.Attitude,
		Content:  finalContent,
	}

	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post review"})
		return
	}
	refreshStoredGameRating(game.GameID)

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

	// 1.5 Verify License is Active
	if license.Status != "ACTIVE" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This game license is not active or has already been refunded"})
		return
	}

	// 2. Prevent Duplicate Refunds
	var existing []models.RefundRequest
	if err := database.DB.Where("transaction_item_id = ?", input.TransactionItemID).Find(&existing).Error; err == nil {
		for _, req := range existing {
			if req.Status == "PENDING" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "A refund request is already pending for this item"})
				return
			}
			if req.Status == "APPROVED" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "This item has already been refunded"})
				return
			}
		}
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
	database.DB.Where("(sender_id = ? OR receiver_id = ?) AND status = ?", userID, userID, "ACCEPTED").Order("created_at desc").Find(&friends)

	type FriendDTO struct {
		FriendshipID  uint          `json:"friendship_id"`
		ID            uint          `json:"id"`
		Username      string        `json:"username"`
		User          socialUserDTO `json:"user"`
		CreatedAt     time.Time     `json:"created_at"`
		LastMessage   string        `json:"last_message"`
		LastMessageAt time.Time     `json:"last_message_at"`
		HasUnread     bool          `json:"has_unread"`
	}

	result := make([]FriendDTO, 0, len(friends))
	for _, friend := range friends {
		friendID := friend.SenderID
		if friendID == userID {
			friendID = friend.ReceiverID
		}
		friendUser := getSocialUser(friendID)
		
		// 檢查黑名單，過濾已被自己封鎖的好友
		var blocks []models.Blacklist
		database.DB.Where("blocker_id = ? AND blocked_id = ?", userID, friendID).Limit(1).Find(&blocks)

		if len(blocks) > 0 {
			continue // 已經封鎖的好友不顯示在好友名單中
		}

		query := database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, friendID, friendID, userID)

		var lastMsg models.Message
		hasMsg := query.Order("sent_at desc, message_id desc").Limit(1).Find(&lastMsg).RowsAffected > 0

		lastMsgContent := ""
		lastMsgAt := friend.CreatedAt // Fallback
		if hasMsg {
			lastMsgContent = lastMsg.Content
			lastMsgAt = lastMsg.SentAt
		}

		// 檢查是否有未讀訊息
		var unreadCount int64
		database.DB.Model(&models.Message{}).Where("sender_id = ? AND receiver_id = ? AND is_read = ?", friendID, userID, false).Count(&unreadCount)

		result = append(result, FriendDTO{
			FriendshipID:  friend.FriendshipID,
			ID:            friendID,
			Username:      friendUser.Username,
			User:          friendUser,
			CreatedAt:     friend.CreatedAt,
			LastMessage:   lastMsgContent,
			LastMessageAt: lastMsgAt,
			HasUnread:     unreadCount > 0,
		})
	}

	// 依照 LastMessageAt (最新訊息時間) 降冪排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastMessageAt.After(result[j].LastMessageAt)
	})

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// SendFriendRequest handles POST /api/social/friends/request
func SendFriendRequest(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var input struct {
		ReceiverID uint   `json:"receiver_id"`
		Username   string `json:"username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.ReceiverID == 0 && input.Username != "" {
		var receiver models.User
		if err := database.DB.Where("username = ?", input.Username).First(&receiver).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if receiver.Permission != "ACTIVE" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is no longer active"})
			return
		}
		input.ReceiverID = receiver.UserID
	}
	if input.ReceiverID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receiver_id or username is required"})
		return
	}
	if input.ReceiverID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send a friend request to yourself"})
		return
	}

	// Double check if ID is given directly instead of username
	if input.Username == "" {
		var receiver models.User
		if err := database.DB.First(&receiver, input.ReceiverID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if receiver.Permission != "ACTIVE" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is no longer active"})
			return
		}
	}

	var blocks []models.Blacklist
	if database.DB.Where("(blocker_id = ? AND blocked_id = ?) OR (blocker_id = ? AND blocked_id = ?)", userID, input.ReceiverID, input.ReceiverID, userID).Limit(1).Find(&blocks).RowsAffected > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "無法發送好友邀請"})
		return
	}

	var existings []models.Friendship
	if database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, input.ReceiverID, input.ReceiverID, userID).Limit(1).Find(&existings).RowsAffected > 0 {
		existing := existings[0]
		if existing.Status == "PENDING" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Friend request already exists"})
			return
		} else if existing.Status == "ACCEPTED" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You are already friends"})
			return
		} else if existing.Status == "DECLINED" {
			existing.SenderID = userID
			existing.ReceiverID = input.ReceiverID
			existing.Status = "PENDING"
			if err := database.DB.Save(&existing).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send friend request"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"message": "Friend request sent"})
			return
		}
	}

	friend := models.Friendship{
		SenderID:   userID,
		ReceiverID: input.ReceiverID,
		Status:     "PENDING",
	}

	if err := database.DB.Create(&friend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send friend request"})
		return
	}
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

	var friend models.Friendship
	if err := database.DB.Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND status = ?", userID, input.ReceiverID, input.ReceiverID, userID, "ACCEPTED").First(&friend).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "無法發送訊息：你們不是好友"})
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

	// Check if I blocked the other user
	var blocks []models.Blacklist
	database.DB.Where("blocker_id = ? AND blocked_id = ?", myID, otherID).Limit(1).Find(&blocks)
	hasBlocked := len(blocks) > 0

	// Mark all unread messages sent by the other user to me as read
	readQuery := database.DB.Model(&models.Message{}).Where("sender_id = ? AND receiver_id = ? AND is_read = ?", otherID, myID, false)
	if hasBlocked {
		readQuery = readQuery.Where("sent_at <= ?", blocks[0].CreatedAt)
	}
	readQuery.Update("is_read", true)

	var messages []models.Message
	if hasBlocked {
		database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ? AND sent_at <= ?)", myID, otherID, otherID, myID, blocks[0].CreatedAt).Order("sent_at asc, message_id asc").Find(&messages)
	} else {
		database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", myID, otherID, otherID, myID).Order("sent_at asc, message_id asc").Find(&messages)
	}

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
		PostAsRole    string `json:"post_as_role"`
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

	var game models.Game
	if err := database.DB.First(&game, review.GameID).Error; err == nil {
		var user models.User
		if err := database.DB.First(&user, userID).Error; err == nil {
			canBypass := false
			if input.PostAsRole == "ADMIN" && user.Role == "ADMIN" {
				canBypass = true
			} else if input.PostAsRole == "CSR" && user.Role == "CSR" {
				canBypass = true
			} else if input.PostAsRole == "AUTHOR" && user.Role == "DEVELOPER" && game.DeveloperID == userID {
				canBypass = true
			} else {
				input.PostAsRole = "USERS"
			}

			if !canBypass {
				// We don't block users from replying even if they don't own the game in original code,
				// but we shouldn't allow them to pretend they are ADMIN.
			}
		}
	} else {
		input.PostAsRole = "USERS"
	}

	finalContent := buildRoleContent(input.Content, input.PostAsRole)

	reply := models.ReviewReply{
		ReviewID:      review.ReviewID,
		UserID:        userID,
		ParentReplyID: input.ParentReplyID,
		Content:       finalContent,
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
		BlockedID uint `json:"blocked_id"`
		UserID    uint `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.BlockedID == 0 {
		input.BlockedID = input.UserID
	}
	if input.BlockedID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "blocked_id is required"})
		return
	}
	if input.BlockedID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot blacklist yourself"})
		return
	}

	blacklist := models.Blacklist{
		BlockerID: userID,
		BlockedID: input.BlockedID,
	}

	if err := database.DB.Create(&blacklist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to blacklist"})
		return
	}

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
	database.DB.Where("(receiver_id = ? OR sender_id = ?) AND status = ?", userID, userID, "PENDING").Order("created_at desc").Find(&requests)

	type RequestDTO struct {
		ID           uint          `json:"id"`
		FriendshipID uint          `json:"friendship_id"`
		SenderID     uint          `json:"sender_id"`
		ReceiverID   uint          `json:"receiver_id"`
		Sender       socialUserDTO `json:"sender"`
		Receiver     socialUserDTO `json:"receiver"`
		Status       string        `json:"status"`
	}

	result := make([]RequestDTO, 0, len(requests))
	for _, req := range requests {
		result = append(result, RequestDTO{
			ID:           req.FriendshipID,
			FriendshipID: req.FriendshipID,
			SenderID:     req.SenderID,
			ReceiverID:   req.ReceiverID,
			Sender:       getSocialUser(req.SenderID),
			Receiver:     getSocialUser(req.ReceiverID),
			Status:       req.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetBlacklist handles GET /api/social/blacklist
func GetBlacklist(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var blacklisted []models.Blacklist
	database.DB.Where("blocker_id = ?", userID).Order("created_at desc").Find(&blacklisted)

	type BlacklistDTO struct {
		ID          uint          `json:"id"`
		BlacklistID uint          `json:"blacklist_id"`
		BlockedID   uint          `json:"blocked_id"`
		Username    string        `json:"username"`
		User        socialUserDTO `json:"user"`
	}

	result := make([]BlacklistDTO, 0, len(blacklisted))
	for _, item := range blacklisted {
		blockedUser := getSocialUser(item.BlockedID)
		result = append(result, BlacklistDTO{ID: item.BlockedID, BlacklistID: item.BlacklistID, BlockedID: item.BlockedID, Username: blockedUser.Username, User: blockedUser})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetMyRefunds handles GET /api/protected/refunds
// It returns the refund history for the current user.
func GetMyRefunds(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var refunds []models.RefundRequest
	// We might need to manually populate Game details later, but for now we fetch the requests.
	if err := database.DB.Where("buyer_id = ?", userID).Order("created_at desc").Find(&refunds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch refunds"})
		return
	}

	// Fetch transaction items and their game details so frontend can display titles
	type RefundWithGameDTO struct {
		models.RefundRequest
		GameTitle string `json:"game_title"`
		GameCover string `json:"game_cover"`
	}

	var results []RefundWithGameDTO
	for _, req := range refunds {
		var item models.TransactionItem
		title := "未知遊戲"
		cover := ""
		if err := database.DB.Preload("Game").Where("item_id = ?", req.TransactionItemID).First(&item).Error; err == nil {
			title = item.Game.Title
			// We won't fetch the full media array here to keep it fast, or we could if needed.
		}

		results = append(results, RefundWithGameDTO{
			RefundRequest: req,
			GameTitle:     title,
			GameCover:     cover,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
