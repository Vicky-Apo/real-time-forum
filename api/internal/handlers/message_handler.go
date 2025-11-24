package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get authenticated user
		user := middleware.GetCurrentUser(r)
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Parse request
		var req models.SendMessageRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Validate content
		if len(req.Content) == 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "Message content cannot be empty")
			return
		}

		if len(req.Content) > 512 {
			utils.RespondWithError(w, http.StatusBadRequest, "Message content too long (max 512 characters)")
			return
		}

		// Validate recipient is not sender
		if req.RecipientID == user.ID {
			utils.RespondWithError(w, http.StatusBadRequest, "Cannot send message to yourself")
			return
		}

		// Save message to database
		response, err := mr.SaveMessage(user.ID, req.RecipientID, req.Content)
		if err != nil {
			log.Printf("Failed to save message: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to send message")
			return
		}

		// Broadcast to WebSocket if recipient is online
		hub.SendMessageToUser(req.RecipientID, models.EventTypeReceiveMessage, models.ReceiveMessagePayload{
			SenderID:   user.ID,
			SenderName: user.Nickname,
			Content:    req.Content,
			SentAt:     response.CreatedAt,
		})

		// Return success response
		utils.RespondWithSuccess(w, http.StatusCreated, response)
	}
}

// GetMessagesHandler retrieves message history between two users
func GetMessagesHandler(mr *repository.MessageRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get authenticated user
		user := middleware.GetCurrentUser(r)
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

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
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get authenticated user
		user := middleware.GetCurrentUser(r)
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

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
