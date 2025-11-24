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
	mux.Handle("/api/auth/register", http.HandlerFunc(handlers.RegisterHandler(UserRepo)))
	mux.Handle("/api/auth/login", http.HandlerFunc(handlers.LoginHandler(UserRepo, SessionRepo)))
	mux.Handle("/api/auth/logout", AuthMiddleware.RequireAuth(handlers.LogoutHandler(UserRepo, SessionRepo)))
	mux.Handle("/api/auth/me", AuthMiddleware.RequireAuth(handlers.GetCurrentUser()))

	// ===== EXISTING POST ROUTES =====
	// All routes protected for private forum
	mux.Handle("/api/posts", AuthMiddleware.RequireAuth(handlers.GetAllPostsHandler(PostRepo)))
	mux.Handle("/api/posts/view/{id}", AuthMiddleware.RequireAuth(handlers.GetSinglePostHandler(PostRepo)))
	mux.Handle("/api/posts/by-category/{id}", AuthMiddleware.RequireAuth(handlers.GetPostsByCategoryHandler(PostRepo)))
	mux.Handle("/api/posts/create", AuthMiddleware.RequireAuth(handlers.CreatePostHandler(PostRepo, CategoryRepo)))

	// ===== EXISTING CATEGORY ROUTES =====
	// Protected for private forum
	mux.Handle("/api/categories", AuthMiddleware.RequireAuth(handlers.GetAllCategoriesHandler(CategoryRepo, PostRepo)))

	// ===== EXISTING COMMENT ROUTES =====
	// All routes protected for private forum
	mux.Handle("/api/comments/for-post/{id}", AuthMiddleware.RequireAuth(handlers.GetCommentsByPostIDHandler(CommentRepo)))
	mux.Handle("/api/comments/create-on-post/{id}", AuthMiddleware.RequireAuth(handlers.CreateCommentHandler(CommentRepo)))

	// ===== WEBSOCKET ROUTES =====
	// Protected - requires authentication
	mux.Handle("/ws", AuthMiddleware.RequireAuth(handlers.WebSocketHandler(hub)))

	// ===== APPLY MIDDLEWARE =====
	handler := RateLimiter.Limit(mux)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.CORS(handler)
	return AuthMiddleware.Authenticate(handler)
}
