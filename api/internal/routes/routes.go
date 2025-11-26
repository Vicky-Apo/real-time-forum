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
	MessageRepo := repository.NewMessageRepository(db)

	// ===== WEBSOCKET HUB =====
	hub := ws.NewHub()
	go hub.Run() // Start the hub in a goroutine

	// ===== EXISTING MIDDLEWARE =====
	AuthMiddleware := middleware.NewMiddleware(UserRepo, SessionRepo)
	RateLimiter := middleware.NewRateLimiter(
		time.Duration(config.Config.RateLimitWindow)*time.Minute,
		config.Config.RateLimitRequests,
	)

	// ===== EXISTING AUTH ROUTES =====
	mux.Handle("POST /api/auth/register", http.HandlerFunc(handlers.RegisterHandler(UserRepo)))
	mux.Handle("POST /api/auth/login", http.HandlerFunc(handlers.LoginHandler(UserRepo, SessionRepo)))
	mux.Handle("POST /api/auth/logout", AuthMiddleware.RequireAuth(handlers.LogoutHandler(UserRepo, SessionRepo)))
	mux.Handle("POST /api/auth/me", AuthMiddleware.RequireAuth(handlers.GetCurrentUser()))

	// ===== EXISTING POST ROUTES =====
	// All routes protected for private forum
	mux.Handle("GET /api/posts", AuthMiddleware.RequireAuth(handlers.GetAllPostsHandler(PostRepo)))
	mux.Handle("GET /api/posts/view/{id}", AuthMiddleware.RequireAuth(handlers.GetSinglePostHandler(PostRepo)))
	mux.Handle("GET /api/posts/by-category/{id}", AuthMiddleware.RequireAuth(handlers.GetPostsByCategoryHandler(PostRepo)))
	mux.Handle("POST /api/posts/create", AuthMiddleware.RequireAuth(handlers.CreatePostHandler(PostRepo, CategoryRepo)))

	// ===== EXISTING CATEGORY ROUTES =====
	// Protected for private forum
	mux.Handle("GET /api/categories", AuthMiddleware.RequireAuth(handlers.GetAllCategoriesHandler(CategoryRepo, PostRepo)))

	// ===== EXISTING COMMENT ROUTES =====
	// All routes protected for private forum
	mux.Handle("GET /api/comments/for-post/{id}", AuthMiddleware.RequireAuth(handlers.GetCommentsByPostIDHandler(CommentRepo)))
	mux.Handle("POST /api/comments/create-on-post/{id}", AuthMiddleware.RequireAuth(handlers.CreateCommentHandler(CommentRepo)))

	// ===== MESSAGE ROUTES =====
	// All routes protected - requires authentication
	mux.Handle("POST /api/messages/send", AuthMiddleware.RequireAuth(handlers.SendMessageHandler(MessageRepo, hub)))
	mux.Handle("GET /api/messages/{id}", AuthMiddleware.RequireAuth(handlers.GetMessagesHandler(MessageRepo)))
	mux.Handle("GET /api/messages/unread-count", AuthMiddleware.RequireAuth(handlers.GetUnreadCountHandler(MessageRepo)))
	mux.Handle("GET /api/conversations", AuthMiddleware.RequireAuth(handlers.GetConversationsHandler(MessageRepo, hub)))

	// ===== USER ROUTES =====
	// All routes protected - requires authentication
	mux.Handle("GET /api/users/online", AuthMiddleware.RequireAuth(handlers.GetOnlineUsersHandler(hub)))

	// ===== WEBSOCKET ROUTES =====
	// Protected - requires authentication
	mux.Handle("/ws", AuthMiddleware.RequireAuth(handlers.WebSocketHandler(hub)))

	// ===== APPLY MIDDLEWARE =====
	handler := RateLimiter.Limit(mux)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.CORS(handler)
	return AuthMiddleware.Authenticate(handler)
}
