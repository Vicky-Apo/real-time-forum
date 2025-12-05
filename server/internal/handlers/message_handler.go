package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"real-time-forum/config"
	"real-time-forum/internal/middleware"
	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/utils"
	ws "real-time-forum/internal/websocket"
)

// SendMessageHandler handles sending a message via HTTP POST
// After saving to DB, it broadcasts to WebSocket if recipient is online
func SendMessageHandler(mr *repository.MessageRepository, hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Check Content-Type to determine if this is a multipart form or JSON
		contentType := r.Header.Get("Content-Type")
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

		var recipientID, content string
		var images []models.MessageImage

		if isMultipart {
			// Parse multipart form (25MB limit)
			err := r.ParseMultipartForm(25 << 20)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse form")
				return
			}

			// Extract text fields
			recipientID = r.FormValue("recipient_id")
			content = r.FormValue("content")

			// Process images (if any)
			files := r.MultipartForm.File["images"]
			if len(files) > 0 {
				images, err = utils.ProcessMessageImageUploads(files, config.Config.MaxImagesPerMessage, config.Config.MaxMessageImageSize)
				if err != nil {
					utils.RespondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
			}
		} else {
			// Parse JSON request (backwards compatibility)
			var req models.SendMessageRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
				return
			}
			recipientID = req.RecipientID
			content = req.Content
		}

		// Validate recipient
		if recipientID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Recipient ID is required")
			return
		}

		// Validate recipient is not sender
		if recipientID == user.ID {
			utils.RespondWithError(w, http.StatusBadRequest, "Cannot send message to yourself")
			return
		}

		// Validate content length
		if len(content) > 512 {
			utils.RespondWithError(w, http.StatusBadRequest, "Message content too long (max 512 characters)")
			return
		}

		// Require either content or images
		if len(content) == 0 && len(images) == 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "Message must have content or images")
			return
		}

		// Save message to database
		var response *models.SendMessageResponse
		var err error

		if len(images) > 0 {
			response, err = mr.SaveMessageWithImages(user.ID, recipientID, content, images)
		} else {
			response, err = mr.SaveMessage(user.ID, recipientID, content)
		}

		if err != nil {
			log.Printf("Failed to save message: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to send message")
			return
		}

		// Fetch saved images with complete data (message_id, uploaded_at) for WebSocket broadcast
		var savedImages []models.MessageImage
		if len(images) > 0 {
			savedImages, err = mr.GetImagesForMessage(response.MessageID)
			if err != nil {
				log.Printf("Failed to fetch saved images: %v", err)
				// Don't fail the request, just send empty images array
				savedImages = []models.MessageImage{}
			}
		}

		// Broadcast to WebSocket if recipient is online
		hub.SendMessageToUser(recipientID, models.EventTypeReceiveMessage, models.ReceiveMessagePayload{
			SenderID:   user.ID,
			SenderName: user.Username,
			Content:    content,
			SentAt:     response.CreatedAt,
			Images:     savedImages,
		})

		// Return success response
		utils.RespondWithSuccess(w, http.StatusCreated, response)
	}
}

// GetMessagesHandler retrieves message history between two users
func GetMessagesHandler(mr *repository.MessageRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Get user ID from URL path
		otherUserID := r.PathValue("id")
		if otherUserID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "User ID is required")
			return
		}

		// Parse query parameters
		limitStr := r.URL.Query().Get("limit")
		limit := 10 // Default
		if limitStr != "" {
			parsedLimit, err := strconv.Atoi(limitStr)
			if err != nil || parsedLimit < 1 || parsedLimit > 50 {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid limit parameter (must be 1-50)")
				return
			}
			limit = parsedLimit
		}

		// Parse optional "before" timestamp for pagination
		var beforeTimestamp *time.Time
		beforeStr := r.URL.Query().Get("before")
		if beforeStr != "" {
			parsed, err := time.Parse(time.RFC3339, beforeStr)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid before parameter (use RFC3339 format)")
				return
			}
			beforeTimestamp = &parsed
		}

		// Get messages from database
		response, err := mr.GetMessages(user.ID, otherUserID, limit, beforeTimestamp)
		if err != nil {
			log.Printf("Failed to get messages: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve messages")
			return
		}

		// Mark messages from the other user as read
		err = mr.MarkMessagesAsRead(user.ID, otherUserID)
		if err != nil {
			log.Printf("Failed to mark messages as read: %v", err)
			// Don't fail the request, just log the error
		}

		// Return messages
		utils.RespondWithSuccess(w, http.StatusOK, response)
	}
}

// GetUnreadCountHandler returns the count of unread messages for the current user
func GetUnreadCountHandler(mr *repository.MessageRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Get unread count
		count, err := mr.GetUnreadCount(user.ID)
		if err != nil {
			log.Printf("Failed to get unread count: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get unread count")
			return
		}

		// Return count
		utils.RespondWithSuccess(w, http.StatusOK, map[string]int{"unread_count": count})
	}
}

// GetConversationsHandler returns all conversations for the current user
// Sorted by last message timestamp (Discord-style), with users without messages alphabetically at the end
func GetConversationsHandler(mr *repository.MessageRepository, hub interface{ GetOnlineUsers() []models.UserStatusPayload }) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Get conversations from database
		conversations, err := mr.GetConversations(user.ID)
		if err != nil {
			log.Printf("Failed to get conversations: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get conversations")
			return
		}

		// Get online users from hub
		onlineUsers := hub.GetOnlineUsers()
		onlineMap := make(map[string]bool)
		for _, u := range onlineUsers {
			onlineMap[u.UserID] = true
		}

		// Add online status to each conversation
		for i := range conversations {
			conversations[i].IsOnline = onlineMap[conversations[i].UserID]
		}

		// Return conversations
		utils.RespondWithSuccess(w, http.StatusOK, map[string]interface{}{
			"conversations": conversations,
		})
	}
}
