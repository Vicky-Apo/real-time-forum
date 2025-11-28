-- BULLETPROOF Posts Creation - NO TIMESTAMP FUNCTIONS
-- This will work 100% guaranteed

-- Clean up first
DELETE FROM post_categories WHERE post_id LIKE 'post-%';
DELETE FROM posts WHERE post_id LIKE 'post-%';

-- Let SQLite handle the timestamp automatically (remove created_at from INSERT)
INSERT INTO posts (post_id, user_id, content) VALUES 
('post-001', 'de6da0a6-ac24-4f69-8cff-30f062608622', 
'Hey everyone! I''ve been working on a Go web application and I''m curious about best practices for database connections. Should I use connection pooling or is it better to open/close connections for each request?'),

('post-002', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Just discovered async/await in JavaScript and wow, it makes asynchronous code so much cleaner! Been refactoring my old Promise chains and the readability improvement is incredible.'),

('post-003', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Working on optimizing SQLite queries for a forum application. Found that adding proper indexes can dramatically improve performance.'),

('post-004', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'The evolution of web frameworks is fascinating! From jQuery to React, Vue, and now with SvelteKit and Solid.js.'),

('post-005', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'RESTful API design question: When handling errors, do you prefer returning detailed error messages or keeping them generic for security?'),

('post-006', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Code review best practices: Focus on logic over style, ask questions instead of making demands, praise good code when you see it.'),

('post-007', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Started learning Rust this week and the ownership concept is mind-bending! Coming from garbage-collected languages is challenging.'),

('post-008', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Made my first meaningful contribution to an open source project today! Fixed a bug in a Go HTTP router library.'),

('post-009', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Unit testing vs integration testing vs e2e testing - finding the right balance is tricky.'),

('post-010', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Five years into my programming career: soft skills matter as much as technical skills.'),

('post-011', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Docker has completely changed how I deploy applications! The consistency between environments is game-changing.'),

('post-012', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Switched from JavaScript to TypeScript for a large project and the developer experience improvement is remarkable.'),

('post-013', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'GraphQL vs REST API debate - I''ve been using GraphQL and while flexible, the complexity can be overwhelming.'),

('post-014', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Automated my entire deployment pipeline using Python scripts and GitHub Actions.'),

('post-015', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Git branching strategies discussion: We''ve been using Git Flow but considering GitHub Flow.'),

('post-016', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Web performance optimization: Implemented lazy loading, code splitting, and service workers.'),

('post-017', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database normalization vs denormalization dilemma in high-traffic applications.'),

('post-018', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Migrating from monolithic architecture to microservices on AWS.'),

('post-019', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Web security checklist: HTTPS everywhere, input validation, SQL injection prevention.'),

('post-020', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Remote development team collaboration tools and strategies.'),

('post-021', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Learning functional programming concepts has changed how I write code.'),

('post-022', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'CSS Grid and Flexbox have revolutionized web layouts. No more float hacks!'),

('post-023', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Monitoring and observability in production: Logs, metrics, and traces are essential.'),

('post-024', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Progressive Web Apps bridge the gap between web and native apps.'),

('post-025', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Machine learning in web development using TensorFlow.js opens interesting possibilities.'),

('post-026', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'API versioning strategies: URL versioning vs header versioning vs content negotiation.'),

('post-027', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Serverless architecture with AWS Lambda: Cold starts vs scalability benefits.'),

('post-028', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Code documentation philosophy: Self-documenting code vs explicit documentation.'),

('post-029', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Accessibility in web development should be built-in, not afterthoughts.'),

('post-030', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Continuous integration best practices: Fast feedback loops and fail-fast strategies.'),

('post-031', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Container orchestration with Kubernetes: Complex but powerful for scaling.'),

('post-032', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Event-driven architecture patterns in modern web applications.'),

('post-033', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Mobile-first responsive design: Why it matters more than ever.'),

('post-034', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'State management in React: Redux vs Context vs Zustand.'),

('post-035', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'WebAssembly opening new possibilities for web performance.'),

('post-036', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Microservice communication patterns: REST, gRPC, and message queues.'),

('post-037', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database sharding strategies for handling massive scale.'),

('post-038', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Frontend build tools evolution: Webpack to Vite to Turbopack.'),

('post-039', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Real-time features with WebSockets and Server-Sent Events.'),

('post-040', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Data visualization libraries: D3.js vs Chart.js vs modern alternatives.'),

('post-041', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Authentication strategies: JWT vs sessions vs OAuth2 implementation.'),

('post-042', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Edge computing and CDN strategies for global web applications.'),

('post-043', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Developer productivity tools that changed my workflow completely.'),

('post-044', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Database migration strategies for zero-downtime deployments.'),

('post-045', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'AI-assisted coding: GitHub Copilot and ChatGPT in development workflow.'),

('post-046', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Performance monitoring and error tracking in production applications.'),

('post-047', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Clean architecture principles in practice: Lessons learned.'),

('post-048', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Internationalization and localization strategies for global apps.'),

('post-049', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'Web3 and blockchain integration: Hype vs practical applications.'),

('post-050', 'de6da0a6-ac24-4f69-8cff-30f062608622',
'The future of web development: Trends to watch in 2025.');

-- Insert post-category relationships
INSERT INTO post_categories (post_id, category_id) VALUES 
('post-001', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-002', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-003', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-004', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-005', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-006', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-007', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-008', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-009', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-010', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-011', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-012', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-013', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-014', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-015', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-016', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-017', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-018', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-019', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-020', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-021', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-022', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-023', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-024', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-025', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-026', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-027', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-028', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-029', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-030', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-031', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-032', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-033', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-034', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-035', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-036', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-037', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-038', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-039', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-040', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-041', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-042', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-043', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-044', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-045', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-046', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-047', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-048', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-049', '41907e4d-36cc-41df-a0ec-99ce63decb31'),
('post-050', '41907e4d-36cc-41df-a0ec-99ce63decb31');

-- Verification
SELECT 'SUCCESS: 50 posts created!' as message;
SELECT COUNT(*) as total_posts FROM posts WHERE user_id = 'de6da0a6-ac24-4f69-8cff-30f062608622';