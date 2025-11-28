package routes

import (
	"database/sql"
	"net/http"
	"time"

	"platform.zone01.gr/git/gpapadopoulos/forum/config"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/handlers"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/middleware"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/repository"
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

	// ===== EXISTING MIDDLEWARE =====
	AuthMiddleware := middleware.NewMiddleware(UserRepo, SessionRepo)
	RateLimiter := middleware.NewRateLimiter(
		time.Duration(config.Config.RateLimitWindow)*time.Minute,
		config.Config.RateLimitRequests,
	)

	// ===== OAUTH HANDLER =====
	OAuthHandler := handlers.NewOAuthHandler(OAuthRepo, UserRepo, SessionRepo, &config.Config)

	// ===== EXISTING AUTH ROUTES =====
	mux.Handle("/api/auth/register", http.HandlerFunc(handlers.RegisterHandler(UserRepo)))
	mux.Handle("/api/auth/login", http.HandlerFunc(handlers.LoginHandler(UserRepo, SessionRepo)))
	mux.Handle("/api/auth/logout", AuthMiddleware.RequireAuth(handlers.LogoutHandler(UserRepo, SessionRepo)))
	mux.Handle("/api/auth/me", AuthMiddleware.RequireAuth(handlers.GetCurrentUser()))

	// ===== SIMPLIFIED OAUTH ROUTES (WEB ONLY) =====
	// GitHub OAuth initiation - SINGLE endpoint
	mux.Handle("/api/auth/github/login", http.HandlerFunc(OAuthHandler.ServeGitHubLogin))
	// GitHub OAuth callback - SINGLE endpoint
	mux.Handle("/api/auth/github/callback", http.HandlerFunc(OAuthHandler.ServeGitHubCallback))
	// Google OAuth initiation - SINGLE endpoint
	mux.Handle("/api/auth/google/login", http.HandlerFunc(OAuthHandler.ServeGoogleLogin))
	// Google OAuth callback - SINGLE endpoint
	mux.Handle("/api/auth/google/callback", http.HandlerFunc(OAuthHandler.ServeGoogleCallback))
	// =================
	// ===== EXISTING USER PROFILE ROUTES =====
	mux.Handle("/api/users/profile/{id}", AuthMiddleware.RequireAuth(handlers.GetUserProfileHandler(UserRepo)))
	mux.Handle("/api/users/posts/{id}", AuthMiddleware.RequireAuth(handlers.GetUserPostsProfileHandler(PostRepo)))
	mux.Handle("/api/users/liked-posts/{id}", AuthMiddleware.RequireAuth(handlers.GetUserLikedPostsProfileHandler(PostRepo)))
	mux.Handle("/api/users/commented-posts/{id}", AuthMiddleware.RequireAuth(handlers.GetUserCommentedPostsProfileHandler(PostRepo)))

	// ===== EXISTING POST ROUTES =====
	// Public GET routes
	mux.Handle("/api/posts", http.HandlerFunc(handlers.GetAllPostsHandler(PostRepo)))
	mux.Handle("/api/posts/view/{id}", http.HandlerFunc(handlers.GetSinglePostHandler(PostRepo)))
	mux.Handle("/api/posts/by-category/{id}", http.HandlerFunc(handlers.GetPostsByCategoryHandler(PostRepo)))

	// Protected POST routes (create only)
	mux.Handle("/api/posts/create", AuthMiddleware.RequireAuth(handlers.CreatePostHandler(PostRepo, CategoryRepo, PostImageRepo)))

	// Protected PUT/DELETE routes (clear naming)
	mux.Handle("/api/posts/edit/{id}", AuthMiddleware.RequireAuth(handlers.UpdatePostHandler(PostRepo, CategoryRepo, PostImageRepo)))
	mux.Handle("/api/posts/remove/{id}", AuthMiddleware.RequireAuth(handlers.DeletePostHandler(PostRepo, CategoryRepo, PostImageRepo)))

	// ===== EXISTING CATEGORY ROUTES =====
	mux.Handle("/api/categories", http.HandlerFunc(handlers.GetAllCategoriesHandler(CategoryRepo, PostRepo)))

	// ===== EXISTING COMMENT ROUTES =====
	// Public GET routes
	mux.Handle("/api/comments/for-post/{id}", http.HandlerFunc(handlers.GetCommentsByPostIDHandler(CommentRepo)))
	// ---  SERVE UPLOADED IMAGES  ---
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(config.Config.UploadDir))))
	// Protected routes
	mux.Handle("/api/comments/create-on-post/{id}", AuthMiddleware.RequireAuth(handlers.CreateCommentHandler(CommentRepo, NotificationRepo, PostRepo, UserRepo)))
	mux.Handle("/api/comments/edit/{id}", AuthMiddleware.RequireAuth(handlers.UpdateCommentHandler(CommentRepo)))
	mux.Handle("/api/comments/remove/{id}", AuthMiddleware.RequireAuth(handlers.DeleteCommentHandler(CommentRepo)))
	mux.Handle("/api/comments/view/{id}", http.HandlerFunc(handlers.GetSingleCommentHandler(CommentRepo)))

	// ===== EXISTING REACTION ROUTES =====
	// Post reactions
	mux.Handle("/api/reactions/posts/toggle", AuthMiddleware.RequireAuth(handlers.TogglePostReactionHandler(PostReactionRepo, NotificationRepo, PostRepo, UserRepo)))

	// Comment reactions
	mux.Handle("/api/reactions/comments/toggle", AuthMiddleware.RequireAuth(handlers.ToggleCommentReactionHandler(CommentReactionRepo, NotificationRepo, CommentRepo, UserRepo, PostRepo)))

	// ===== NEW NOTIFICATION ROUTES =====
	mux.Handle("/api/notifications", AuthMiddleware.RequireAuth(handlers.GetNotificationsHandler(NotificationRepo)))
	mux.Handle("/api/notifications/mark-read/{id}", AuthMiddleware.RequireAuth(handlers.MarkAsReadHandler(NotificationRepo)))

	// ===== APPLY MIDDLEWARE =====
	handler := RateLimiter.Limit(mux)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.CORS(handler)
	return AuthMiddleware.Authenticate(handler)
}
