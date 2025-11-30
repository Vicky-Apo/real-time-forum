package database

// TableCreationStatements contains all SQL statements for optimized table creation
var TableCreationStatements = []string{

	`CREATE TABLE IF NOT EXISTS users (
		user_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		username VARCHAR(15) NOT NULL UNIQUE,
		age INTEGER NOT NULL CHECK (age >= 13 AND age <= 120),
		gender TEXT NOT NULL,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL, -- Store hashed password
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`,

	`CREATE TABLE IF NOT EXISTS oauth_user_accounts (
    user_id TEXT NOT NULL,
    provider TEXT NOT NULL,                    -- 'github', 'google'
    provider_user_id TEXT NOT NULL,            -- GitHub user ID (e.g., "12345678")
    provider_email TEXT,                       -- Email from GitHub
    provider_username TEXT,                    -- GitHub username
    access_token TEXT NOT NULL,                -- OAuth access token
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, provider),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    UNIQUE(provider, provider_user_id)         -- One GitHub account = one forum user
	);`,

	`CREATE TABLE IF NOT EXISTS oauth_flow_states (
    state_id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,                    -- 'github', 'google'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,             -- Expires in 15 minutes
    
    CHECK (provider IN ('github', 'google'))
    );`,

	`CREATE TABLE IF NOT EXISTS sessions (
		user_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		session_id TEXT NOT NULL UNIQUE,
		ip_address TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	);`,

	`CREATE TABLE IF NOT EXISTS categories (
		category_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		category_name TEXT NOT NULL UNIQUE
	);`,

	`CREATE TABLE IF NOT EXISTS posts (
    post_id TEXT PRIMARY KEY NOT NULL UNIQUE,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,

    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);`,

	`CREATE TABLE IF NOT EXISTS post_categories (
		post_id TEXT NOT NULL,
		category_id TEXT NOT NULL,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
	);`,

	`CREATE TABLE IF NOT EXISTS comments (
		comment_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,
		
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	);`,

	`CREATE TABLE IF NOT EXISTS post_reactions (
		user_id TEXT NOT NULL,
		post_id TEXT NOT NULL,
		reaction_type INTEGER NOT NULL, -- 1 for like, 2 for dislike
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		-- Natural primary key - no UUID needed, prevents duplicate reactions
		PRIMARY KEY (user_id, post_id),
		
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
		
		-- Ensure valid reaction types
		CHECK (reaction_type IN (1, 2))
	);`,

	`CREATE TABLE IF NOT EXISTS comment_reactions (
		user_id TEXT NOT NULL,
		comment_id TEXT NOT NULL,
		reaction_type INTEGER NOT NULL, -- 1 for like, 2 for dislike
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		-- Natural primary key - no UUID needed, prevents duplicate reactions
		PRIMARY KEY (user_id, comment_id),
		
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE,
		
		-- Ensure valid reaction types
		CHECK (reaction_type IN (1, 2))
	);`,

	`CREATE TABLE IF NOT EXISTS post_images (
		image_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		post_id TEXT NOT NULL,
		image_url TEXT NOT NULL,
		original_filename TEXT,
		uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
	);`,

	`CREATE TABLE IF NOT EXISTS notifications (
		notification_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		user_id TEXT NOT NULL,              -- who gets the notification
		trigger_username TEXT NOT NULL,     -- who caused it (e.g., "John")
		post_content_preview TEXT NOT NULL, -- first 50 chars of post content
		post_id TEXT NOT NULL,              -- link to the post
		action TEXT NOT NULL,               -- flexible action description
		is_read BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		-- REMOVED: Restrictive CHECK constraint to allow flexible action text
	);`,

	`CREATE TABLE IF NOT EXISTS messages (
		message_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		sender_id TEXT NOT NULL,
		recipient_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		is_read BOOLEAN NOT NULL DEFAULT 0,

		FOREIGN KEY (sender_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (recipient_id) REFERENCES users(user_id) ON DELETE CASCADE
	);`,
}

// IndexCreationStatements contains ESSENTIAL indexes only - what you'll actually use
var IndexCreationStatements = []string{
	// Authentication indexes (used every request)
	`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,                 // Login by email
	`CREATE INDEX IF NOT EXISTS idx_sessions_session_id ON sessions(session_id);`, // Session validation

	// Core post browsing indexes (main forum functionality)
	`CREATE INDEX IF NOT EXISTS idx_posts_created_desc ON posts(created_at DESC);`,                           // Homepage post list
	`CREATE INDEX IF NOT EXISTS idx_post_categories_category_post ON post_categories(category_id, post_id);`, // Posts by category

	// Comment indexes (viewing posts with comments)
	`CREATE INDEX IF NOT EXISTS idx_comments_post_created ON comments(post_id, created_at ASC);`, // Comments for a post

	// Reaction indexes (like/dislike counts)
	`CREATE INDEX IF NOT EXISTS idx_post_reactions_post_type ON post_reactions(post_id, reaction_type);`,             // Post reaction counts
	`CREATE INDEX IF NOT EXISTS idx_comment_reactions_comment_type ON comment_reactions(comment_id, reaction_type);`, // Comment reaction counts

	`CREATE INDEX IF NOT EXISTS idx_oauth_user_accounts_user ON oauth_user_accounts(user_id);`, // Fast lookup by user ID

	// Fast lookup during OAuth callback
	`CREATE INDEX IF NOT EXISTS idx_oauth_user_accounts_provider ON oauth_user_accounts(provider, provider_user_id);`,

	// Cleanup expired states
	`CREATE INDEX IF NOT EXISTS idx_oauth_flow_states_expires ON oauth_flow_states(expires_at);`,

	// Quickly fetch images for a post
	`CREATE INDEX IF NOT EXISTS idx_post_images_post_id ON post_images(post_id);`,

	// User's notifications ordered by time
	`CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);`,
}
