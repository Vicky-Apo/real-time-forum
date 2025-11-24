package database

// TableCreationStatements contains all SQL statements for optimized table creation
var TableCreationStatements = []string{
	// Users table - updated with new fields for real-time-forum
	`CREATE TABLE IF NOT EXISTS users (
		user_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		nickname VARCHAR(15) NOT NULL UNIQUE,
		age INTEGER NOT NULL CHECK (age >= 13 AND age <= 120),
		gender TEXT NOT NULL,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`,

	// Sessions table - unchanged, already optimized
	`CREATE TABLE IF NOT EXISTS sessions (
		user_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		session_id TEXT NOT NULL UNIQUE,
		ip_address TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	);`,

	// Categories table - unchanged, already optimized
	`CREATE TABLE IF NOT EXISTS categories (
		category_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		category_name TEXT NOT NULL UNIQUE
	);`,

	// Posts table - CLEAN, no denormalized counts
	`CREATE TABLE IF NOT EXISTS posts (
		post_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL,


		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	);`,

	// Post categories junction table - unchanged
	`CREATE TABLE IF NOT EXISTS post_categories (
		post_id TEXT NOT NULL,
		category_id TEXT NOT NULL,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
	);`,

	// Comments table - CLEAN, no denormalized counts
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

	// Messages table - for direct user-to-user messaging
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

	// Message indexes (chat functionality)
	`CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(sender_id, recipient_id, created_at DESC);`, // Messages between two users
	`CREATE INDEX IF NOT EXISTS idx_messages_recipient_time ON messages(recipient_id, created_at DESC);`,          // Unread messages for a user
}

// // WALModeStatements contains SQL statements for enabling WAL mode and performance optimization
// var WALModeStatements = []string{
// 	// Enable WAL mode for better concurrency
// 	`PRAGMA journal_mode=WAL;`,

// 	// Optimize WAL performance
// 	`PRAGMA synchronous=NORMAL;`,  // Faster than FULL, still safe
// 	`PRAGMA cache_size=10000;`,    // Larger cache for better performance
// 	`PRAGMA temp_store=memory;`,   // Keep temporary data in memory
// 	`PRAGMA mmap_size=268435456;`, // 256MB memory mapping for larger databases

// 	// WAL checkpoint optimization
// 	`PRAGMA wal_autocheckpoint=1000;`, // Checkpoint every 1000 pages
// }
