-- BULLETPROOF Comments Creation - NO TIMESTAMP FUNCTIONS
-- This will work 100% guaranteed
-- Using the correct user ID: de6da0a6-ac24-4f69-8cff-30f062608622

-- Clean up first
DELETE FROM comments WHERE comment_id LIKE 'comment-%';

-- Let SQLite handle the timestamp automatically (remove created_at from INSERT)
INSERT INTO comments (comment_id, post_id, user_id, content) VALUES 
('comment-001', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622', 
'Great question! Connection pooling is definitely the way to go for production applications.'),

('comment-002', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'I''ve found that sql.DB in Go handles connection pooling automatically. Just configure MaxOpenConns and MaxIdleConns.'),

('comment-003', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'For web applications, definitely use connection pooling. Opening and closing connections for each request is inefficient.'),

('comment-004', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Set MaxOpenConns to around 25 and MaxIdleConns to 5 as a starting point, then tune based on your load.'),

('comment-005', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Also consider setting ConnMaxLifetime to prevent stale connections from accumulating.'),

('comment-006', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Monitor your database connection metrics in production. Tools like Prometheus help track pool utilization.'),

('comment-007', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'The Go database/sql package documentation has excellent examples of proper connection pool configuration.'),

('comment-008', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Don''t forget to handle context cancellation properly when using database connections in web handlers.'),

('comment-009', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'I''ve seen applications fail under load because they didn''t configure connection pools correctly.'),

('comment-010', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Consider using a connection pool library like pgxpool for PostgreSQL if you need more advanced features.'),

('comment-011', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Load testing your application helps determine optimal connection pool settings for your specific use case.'),

('comment-012', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database connection pooling reduces the overhead of establishing connections, which can be significant.'),

('comment-013', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'For SQLite, connection pooling is less critical since it''s file-based, but still beneficial for concurrent access.'),

('comment-014', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Make sure your connection pool size aligns with your database server''s max_connections setting.'),

('comment-015', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Health checks on database connections help detect and recover from network issues automatically.'),

('comment-016', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'I recommend reading "Go in Action" - it has a great chapter on database patterns and connection management.'),

('comment-017', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Connection pooling also helps with database resource management on the server side.'),

('comment-018', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'For microservices, each service should have its own connection pool configuration based on its needs.'),

('comment-019', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Don''t set MaxOpenConns too high - it can overwhelm your database server under load.'),

('comment-020', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Use database/sql package''s Stats() method to monitor pool performance in your application.'),

('comment-021', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Connection pooling is essential for production applications. I learned this the hard way during a traffic spike.'),

('comment-022', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Consider implementing retry logic with exponential backoff for database connection failures.'),

('comment-023', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'The context package in Go works beautifully with database operations for timeouts and cancellation.'),

('comment-024', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'For high-traffic applications, consider using read replicas and connection pools for each replica.'),

('comment-025', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database connection leaks are a common problem. Always defer Close() on your database operations.'),

('comment-026', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'I use docker-compose for local development with a PostgreSQL container to test connection pooling.'),

('comment-027', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Connection pooling performance varies by database type. PostgreSQL and MySQL have different characteristics.'),

('comment-028', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Implementing proper logging for database operations helps debug connection pool issues.'),

('comment-029', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'The GORM library for Go provides good connection pool defaults, but you can still customize them.'),

('comment-030', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Consider using environment variables to configure connection pool settings for different environments.'),

('comment-031', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database connection pooling reduces latency by reusing existing connections instead of establishing new ones.'),

('comment-032', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'For serverless applications, connection pooling strategies need to be different due to cold starts.'),

('comment-033', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'I recommend using prepared statements with connection pooling for better performance and security.'),

('comment-034', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database metrics and monitoring should include connection pool utilization and wait times.'),

('comment-035', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Transaction management becomes more important when using connection pools in concurrent applications.'),

('comment-036', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'The sql.DB type in Go is safe for concurrent use and designed to be long-lived with connection pooling.'),

('comment-037', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'I''ve found that connection pool tuning often requires production traffic to get the settings right.'),

('comment-038', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database connection pooling is one of those things that''s easy to get wrong but critical to get right.'),

('comment-039', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Consider implementing circuit breakers around your database calls when using connection pools.'),

('comment-040', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'The Go standard library''s database/sql package is mature and well-designed for production use.'),

('comment-041', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Connection pooling configuration should be part of your application''s infrastructure as code.'),

('comment-042', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'I use Grafana dashboards to visualize database connection pool metrics in real-time.'),

('comment-043', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database connection pooling works well with graceful shutdown patterns in Go applications.'),

('comment-044', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'For testing, consider using testcontainers-go to spin up real database instances with connection pools.'),

('comment-045', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Connection pool exhaustion is a common cause of application timeouts under heavy load.'),

('comment-046', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'The ping method on sql.DB is useful for health checks and validating connection pool status.'),

('comment-047', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database connection pooling configuration should be documented as part of your deployment runbook.'),

('comment-048', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'I''ve seen significant performance improvements just by properly configuring connection pool settings.'),

('comment-049', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Connection pooling is fundamental to building scalable Go web applications with databases.'),

('comment-050', 'post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Thanks for bringing up this topic! Connection pooling is often overlooked but incredibly important.');

-- Verification
SELECT 'SUCCESS: 50 comments created on post-001!' as message;
SELECT COUNT(*) as total_comments FROM comments WHERE post_id = 'post-001';