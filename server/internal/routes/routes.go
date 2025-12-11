package routes

import (
	"database/sql"
	"net/http"
	"time"

	"real-time-forum/config"
	"real-time-forum/internal/handlers"
	"real-time-forum/internal/middleware"
	"real-time-forum/internal/repository"
	ws "real-time-forum/internal/websocket"
)

func SetupRoutes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// ===== EXISTING REPOSITORIES =====
	UserRepo := repository.NewUserRepository(db)
	SessionRepo := repository.NewSessionRepository(db)
	PostRepo := repository.NewPostsRepository(db)
	CategoryRepo := repository.NewCategoryRepository(db)
	CommentRepo := repository.NewCommentRepository(db)
	PostReactionRepo := repository.NewPostReactionRepository(db)
	CommentReactionRepo := repository.NewCommentReactionRepository(db)
	OAuthRepo := repository.NewOAuthRepository(db)
	PostImageRepo := repository.NewPostImagesRepository(db)
	NotificationRepo := repository.NewNotificationRepository(db)
	MessageImageRepo := repository.NewMessageImageRepository(db)
	MessageRepo := repository.NewMessageRepository(db, MessageImageRepo)

	// ===== WEBSOCKET HUB =====
	hub := ws.NewHub()
	go hub.Run() // Start the hub in a goroutine

	// ===== EXISTING MIDDLEWARE =====
	AuthMiddleware := middleware.NewMiddleware(UserRepo, SessionRepo)
	RateLimiter := middleware.NewRateLimiter(
		time.Duration(config.Config.RateLimitWindow)*time.Minute,
		config.Config.RateLimitRequests,
	)

	// ===== OAUTH HANDLER =====
	OAuthHandler := handlers.NewOAuthHandler(OAuthRepo, UserRepo, SessionRepo, &config.Config)

	// ===== EXISTING AUTH ROUTES =====
	mux.Handle("POST /api/auth/register", http.HandlerFunc(handlers.RegisterHandler(UserRepo)))
	mux.Handle("POST /api/auth/login", http.HandlerFunc(handlers.LoginHandler(UserRepo, SessionRepo)))
	mux.Handle("POST /api/auth/logout", AuthMiddleware.RequireAuth(handlers.LogoutHandler(UserRepo, SessionRepo)))
	mux.Handle("POST /api/auth/me", AuthMiddleware.RequireAuth(handlers.GetCurrentUser()))

	// ===== SIMPLIFIED OAUTH ROUTES (WEB ONLY) =====
	// GitHub OAuth initiation - SINGLE endpoint
	mux.Handle("GET /api/auth/github/login", http.HandlerFunc(OAuthHandler.ServeGitHubLogin))
	// GitHub OAuth callback - SINGLE endpoint
	mux.Handle("GET /api/auth/github/callback", http.HandlerFunc(OAuthHandler.ServeGitHubCallback))
	// Google OAuth initiation - SINGLE endpoint
	mux.Handle("GET /api/auth/google/login", http.HandlerFunc(OAuthHandler.ServeGoogleLogin))
	// Google OAuth callback - SINGLE endpoint
	mux.Handle("GET /api/auth/google/callback", http.HandlerFunc(OAuthHandler.ServeGoogleCallback))
	// =================
	// ===== EXISTING USER PROFILE ROUTES =====
	mux.Handle("GET /api/users/profile/{id}", AuthMiddleware.RequireAuth(handlers.GetUserProfileHandler(UserRepo)))
	mux.Handle("GET /api/users/posts/{id}", AuthMiddleware.RequireAuth(handlers.GetUserPostsProfileHandler(PostRepo)))
	mux.Handle("GET /api/users/liked-posts/{id}", AuthMiddleware.RequireAuth(handlers.GetUserLikedPostsProfileHandler(PostRepo)))
	mux.Handle("GET /api/users/commented-posts/{id}", AuthMiddleware.RequireAuth(handlers.GetUserCommentedPostsProfileHandler(PostRepo)))

	// ===== EXISTING POST ROUTES =====
	// Private GET routes
	mux.Handle("GET /api/posts", AuthMiddleware.RequireAuth(http.HandlerFunc(handlers.GetAllPostsHandler(PostRepo))))
	mux.Handle("GET /api/posts/view/{id}", AuthMiddleware.RequireAuth(http.HandlerFunc(handlers.GetSinglePostHandler(PostRepo))))
	mux.Handle("GET /api/posts/by-category/{id}", AuthMiddleware.RequireAuth(http.HandlerFunc(handlers.GetPostsByCategoryHandler(PostRepo))))

	// Protected POST routes (create only)
	mux.Handle("POST /api/posts/create", AuthMiddleware.RequireAuth(handlers.CreatePostHandler(PostRepo, CategoryRepo, PostImageRepo)))

	// Protected PUT/DELETE routes (clear naming)
	mux.Handle("PUT /api/posts/edit/{id}", AuthMiddleware.RequireAuth(handlers.UpdatePostHandler(PostRepo, CategoryRepo, PostImageRepo)))
	mux.Handle("DELETE /api/posts/remove/{id}", AuthMiddleware.RequireAuth(handlers.DeletePostHandler(PostRepo, CategoryRepo, PostImageRepo)))

	// ===== EXISTING CATEGORY ROUTES =====
	mux.Handle("GET /api/categories", http.HandlerFunc(handlers.GetAllCategoriesHandler(CategoryRepo, PostRepo)))

	// ===== EXISTING COMMENT ROUTES =====
	// Private GET routes
	mux.Handle("GET /api/comments/for-post/{id}", AuthMiddleware.RequireAuth(http.HandlerFunc(handlers.GetCommentsByPostIDHandler(CommentRepo))))

	// ---  SERVE STATIC FILES  ---
	// Serve uploaded images (user content)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(config.Config.UploadDir))))
	// Serve client static assets (logos, icons, etc.)
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("client/images"))))
	// Protected routes
	mux.Handle("POST /api/comments/create-on-post/{id}", AuthMiddleware.RequireAuth(handlers.CreateCommentHandler(CommentRepo, NotificationRepo, PostRepo, UserRepo)))
	mux.Handle("PUT /api/comments/edit/{id}", AuthMiddleware.RequireAuth(handlers.UpdateCommentHandler(CommentRepo)))
	mux.Handle("DELETE /api/comments/remove/{id}", AuthMiddleware.RequireAuth(handlers.DeleteCommentHandler(CommentRepo)))
	mux.Handle("GET /api/comments/view/{id}", AuthMiddleware.RequireAuth(http.HandlerFunc(handlers.GetSingleCommentHandler(CommentRepo))))

	// ===== EXISTING REACTION ROUTES =====
	// Post reactions
	mux.Handle("POST /api/reactions/posts/toggle", AuthMiddleware.RequireAuth(handlers.TogglePostReactionHandler(PostReactionRepo, NotificationRepo, PostRepo, UserRepo)))

	// Comment reactions
	mux.Handle("POST /api/reactions/comments/toggle", AuthMiddleware.RequireAuth(handlers.ToggleCommentReactionHandler(CommentReactionRepo, NotificationRepo, CommentRepo, UserRepo, PostRepo)))

	// ===== NEW NOTIFICATION ROUTES =====
	mux.Handle("GET /api/notifications", AuthMiddleware.RequireAuth(handlers.GetNotificationsHandler(NotificationRepo)))
	mux.Handle("POST /api/notifications/mark-read/{id}", AuthMiddleware.RequireAuth(handlers.MarkAsReadHandler(NotificationRepo)))

	// ===== MESSAGE ROUTES =====
	// All routes protected - requires authentication
	mux.Handle("POST /api/messages/send", AuthMiddleware.RequireAuth(handlers.SendMessageHandler(MessageRepo, hub)))
	mux.Handle("GET /api/messages/{id}", AuthMiddleware.RequireAuth(handlers.GetMessagesHandler(MessageRepo)))
	mux.Handle("GET /api/messages/unread-count", AuthMiddleware.RequireAuth(handlers.GetUnreadCountHandler(MessageRepo)))
	mux.Handle("GET /api/conversations", AuthMiddleware.RequireAuth(handlers.GetConversationsHandler(MessageRepo, hub)))

	// ===== USER ROUTES =====
	// All routes protected - requires authentication
	// Note: User list with online/offline status is available via /api/conversations

	// ===== WEBSOCKET ROUTES =====
	// Protected - requires authentication
	mux.Handle("/ws", AuthMiddleware.RequireAuth(handlers.WebSocketHandler(hub)))

	// ===== APPLY MIDDLEWARE =====
	handler := RateLimiter.Limit(mux)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.CORS(handler)
	return AuthMiddleware.Authenticate(handler)
}
