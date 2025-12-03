// routes.go
package routes

import (
	"net/http"

	"frontend-service/config"
	"frontend-service/internal/handlers"
	"frontend-service/internal/services"
)

// App holds all application dependencies in one place
type App struct {
	Services *Services
	Config   *config.Config
}

// Services contains all service dependencies
type Services struct {
	Auth            *services.AuthService
	Post            *services.PostService
	Category        *services.CategoryService
	User            *services.UserService
	Comment         *services.CommentService
	PostReaction    *services.PostReactionService
	CommentReaction *services.CommentReactionService
	Notification    *services.NotificationService
	OAuth           *services.OAuthService
	Template        *services.TemplateService
}

// NewApp creates the application with all dependencies
func NewApp(cfg *config.Config) (*App, error) {
	baseClient := services.NewBaseClient(cfg.APIBaseURL)

	templateService, err := services.NewTemplateService(cfg.TemplatesDir)
	if err != nil {
		return nil, err
	}

	return &App{
		Services: &Services{
			Auth:            services.NewAuthService(baseClient, cfg),
			Post:            services.NewPostService(baseClient),
			Category:        services.NewCategoryService(baseClient),
			User:            services.NewUserService(baseClient),
			Comment:         services.NewCommentService(baseClient),
			PostReaction:    services.NewPostReactionService(baseClient),
			CommentReaction: services.NewCommentReactionService(baseClient),
			Notification:    services.NewNotificationService(baseClient),
			OAuth:           services.NewOAuthService(baseClient, cfg),
			Template:        templateService,
		},
		Config: cfg,
	}, nil
}

// SetupRoutes configures all routes - ONE parameter!
func SetupRoutes(app *App) *http.ServeMux {
	mux := http.NewServeMux()
	s := app.Services // shorthand

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static/"))))

	// Create handlers with services
	authH := handlers.NewAuthHandler(s.Auth, s.Template, app.Config)
	oauthH := handlers.NewOAuthHandler(s.Auth, s.OAuth, s.Template, app.Config)
	homeH := handlers.NewHomeHandler(s.Auth, s.Post, s.Category, s.Template)
	categoryH := handlers.NewCategoryHandler(s.Auth, s.Post, s.Category, s.Template)
	postH := handlers.NewPostHandler(s.Auth, s.Post, s.Template)
	createPostH := handlers.NewCreatePostHandler(s.Auth, s.Post, s.Category, s.Template)
	editPostH := handlers.NewEditPostHandler(s.Auth, s.Post, s.Category, s.Template)
	deletePostH := handlers.NewDeletePostHandler(s.Auth, s.Post, s.Template)
	profileH := handlers.NewProfileHandler(s.Auth, s.User, s.Template)
	commentH := handlers.NewCommentHandler(s.Auth, s.Comment, s.Post, s.Template)
	reactionH := handlers.NewPostReactionHandler(s.Auth, s.PostReaction, s.CommentReaction, s.Template)
	notifH := handlers.NewNotificationHandler(s.Auth, s.Notification, s.Template)

	// Route registration - compact and clear
	routes := map[string]http.HandlerFunc{
		// Auth routes
		"/register":          authH.ServeRegister,
		"/login":             authH.ServeLogin,
		"/logout":            authH.ServeLogout,
		"/auth/github/login": oauthH.ServeGitHubLogin,
		"/auth/google/login": oauthH.ServeGoogleLogin,

		// Content routes
		"/":                 homeH.ServeHome,
		"/category/{id}":    categoryH.ServeCategoryPosts,
		"/post/{id}":        postH.ServePostView,
		"/create-post":      createPostH.ServeCreatePost,
		"/edit-post/{id}":   editPostH.ServeEditPost,
		"/delete-post/{id}": deletePostH.ServeDeletePost,

		// Profile routes
		"/profile":                 profileH.ServeProfile,
		"/profile/my-posts":        profileH.ServeUserPosts,
		"/profile/liked-posts":     profileH.ServeUserLikedPosts,
		"/profile/commented-posts": profileH.ServeUserCommentedPosts,

		// Comment routes
		"/api/comments/create/{post_id}":    commentH.ServeCreateComment,
		"/api/comments/edit/{comment_id}":   commentH.ServeEditComment,
		"/api/comments/delete/{comment_id}": commentH.ServeDeleteComment,
		"/edit-comment/{id}":                commentH.ServeEditCommentForm,
		"/edit-comment/{id}/submit":         commentH.ServeEditCommentSubmit,

		// Reaction & notification routes
		"/reactions/posts/toggle":       reactionH.ServeTogglePostReaction,
		"/reactions/comments/toggle":    reactionH.ServeToggleCommentReaction,
		"/notifications":                notifH.ServeNotifications,
		"/api/notifications":            notifH.GetNotificationsAPI,
		"/notifications/mark-read/{id}": notifH.HandleMarkAsRead,
	}

	// Register all routes in a loop
	for pattern, handler := range routes {
		mux.HandleFunc(pattern, handler)
	}

	return mux
}
