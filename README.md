# Real-Time Forum 🚀

A modern, full-stack real-time forum application built with Go (backend) and vanilla JavaScript SPA (frontend), featuring WebSocket-based live chat, user authentication, posts, comments, reactions, and OAuth integration.

## ✨ Key Features

### Core Functionality

- **User Authentication**: Secure registration/login with session management + bcrypt password hashing
- **OAuth Integration**: Sign in with GitHub or Google
- **Posts & Comments**: Full CRUD operations with image upload support
- **Reactions System**: Like/dislike for posts and comments
- **Categories**: IT-focused categories (Programming, Web Dev, DevOps, etc.)
- **User Profiles**: View statistics, activity, and created content

### Real-Time Features

- **Live Chat**: WebSocket-based private messaging between users
- **Typing Indicators**: See when someone is typing
- **Real-Time Notifications**: Instant updates for comments, reactions, and messages
- **Online Status**: See who's currently online

### Technical Highlights

- **Single Page Application (SPA)**: Client-side routing, no page reloads
- **RESTful API**: Clean, stateless backend architecture
- **WebSocket Support**: Bidirectional real-time communication
- **Docker Ready**: Full containerization with docker-compose
- **Security**: Rate limiting, CORS, XSS protection, secure session handling

## 🏗️ Architecture

```text
┌─────────────────┐         ┌──────────────────┐         ┌─────────────┐
│   Browser       │         │  Frontend Server │         │   Backend   │
│   (SPA)         │◄────────│  (Go Proxy)      │◄────────│   API       │
│                 │  HTTP   │                  │  HTTP   │   (Go)      │
│  - HTML/CSS/JS  │         │  - Static Files  │         │  - REST API │
│  - Router       │         │  - API Proxy     │         │  - WebSocket│
│  - State Mgmt   │         │  - WS Proxy      │         │  - SQLite   │
└─────────────────┘         └──────────────────┘         └─────────────┘
         │                           │                           │
         │◄──────────────────────────┴───────────────────────────┘
         │                  WebSocket Connection
         └────────────────────────────────────────────────────────►
```

## 📋 Prerequisites

- **Go**: Version 1.21+ (backend) and 1.24+ (client proxy)
- **Docker & Docker Compose**: For containerized deployment
- **Modern Browser**: Chrome, Firefox, Safari, or Edge (ES6+ support)

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd real-time-forum
   ```

2. **Start all services**

   ```bash
   docker-compose up --build
   ```

3. **Access the application**
   - Frontend: <http://localhost:3000>
   - Backend API: <http://localhost:8080>

4. **Stop services**

   ```bash
   docker-compose down
   ```

### Local Development (Without Docker)

1. **Start backend**

   ```bash
   make backend
   # or
   cd server && go run ./cmd
   ```

2. **Start frontend** (in a new terminal)

   ```bash
   make frontend
   # or
   cd client && go run main.go
   ```

3. **Run both simultaneously**

   ```bash
   make dev
   ```

## 📁 Project Structure

```text
real-time-forum/
├── server/                      # Backend API
│   ├── cmd/
│   │   └── main.go              # Entry point
│   ├── config/                  # Configuration management
│   ├── database/                # Database initialization & migrations
│   ├── internal/
│   │   ├── handlers/            # HTTP & WebSocket handlers
│   │   │   ├── user_handler.go
│   │   │   ├── post_handler.go
│   │   │   ├── comment_handler.go
│   │   │   ├── message_handler.go
│   │   │   ├── notification_handler.go
│   │   │   ├── oauth_handler.go
│   │   │   └── websocket_handler.go
│   │   ├── middleware/          # Auth, CORS, rate limiting, security headers
│   │   ├── models/              # Data structures
│   │   ├── repository/          # Data access layer (SQLite)
│   │   ├── routes/              # API route definitions
│   │   ├── utils/               # Helpers (validation, cookies, tokens, images)
│   │   └── websocket/           # WebSocket hub & client management
│   ├── Dockerfile               # Backend container config
│   └── go.mod
│
├── client/                      # Frontend SPA + Proxy Server
│   ├── main.go                  # Go proxy server (serves SPA + proxies API/WS)
│   ├── index.html               # SPA entry point
│   ├── css/                     # Stylesheets
│   │   ├── global.css
│   │   ├── components.css
│   │   ├── chat.css
│   │   ├── post.css
│   │   └── ...
│   ├── js/
│   │   ├── main.js              # App initialization
│   │   ├── router.js            # Client-side routing
│   │   ├── state.js             # Global state management
│   │   ├── config.js            # Configuration
│   │   ├── api/
│   │   │   └── client.js        # API client wrapper
│   │   ├── components/
│   │   │   ├── Navbar.js
│   │   │   └── Footer.js
│   │   ├── views/               # SPA pages
│   │   │   ├── HomeView.js
│   │   │   ├── LoginView.js
│   │   │   ├── RegisterView.js
│   │   │   ├── PostView.js
│   │   │   ├── CreatePostView.js
│   │   │   ├── ChatView.js
│   │   │   ├── ProfileView.js
│   │   │   └── NotificationsView.js
│   │   ├── websocket/
│   │   │   └── WebSocketManager.js  # WebSocket connection handler
│   │   └── utils/
│   │       ├── logger.js        # Development logging
│   │       ├── constants.js     # App constants
│   │       ├── sanitize.js      # XSS protection
│   │       └── helpers.js
│   ├── images/                  # Static assets
│   ├── Dockerfile               # Frontend container config
│   └── go.mod
│
├── docker-compose.yml           # Multi-container orchestration
├── Makefile                     # Development commands
└── README.md                    # This file
```

## ⚙️ Configuration

### Environment Variables

#### Backend (docker-compose.yml or server/.env)

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_PATH=./DBPath/forum.db
DB_MAX_CONNECTIONS=10

# URLs (use localhost for OAuth callbacks, internal for Docker networking)
FRONTEND_BASE_URL=http://localhost:3000
BACKEND_BASE_URL=http://localhost:8080

# Security
SESSION_DURATION=24h
BCRYPT_COST=10
SESSION_NAME=forum_session
ENVIRONMENT=development

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://frontend:3000
ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With

# Content Limits
MAX_POST_CONTENT_LENGTH=500
MIN_POST_CONTENT_LENGTH=10
MAX_COMMENT_LENGTH=150
MIN_COMMENT_LENGTH=5
MAX_USERNAME_LENGTH=15
MIN_USERNAME_LENGTH=5
MAX_PASSWORD_LENGTH=15
MIN_PASSWORD_LENGTH=3

# OAuth (Optional - leave empty to disable)
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URI=http://localhost:8080/api/auth/github/callback

GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URI=http://localhost:8080/api/auth/google/callback

# Rate Limiting
RATE_LIMIT_REQUESTS=100000
RATE_LIMIT_WINDOW=60

# Pagination
DEFAULT_PAGE_SIZE=20
MAX_PAGE_SIZE=50

# Images
UPLOAD_DIR=./uploads/
MAX_IMAGES_PER_POST=5
MAX_IMAGES_PER_MESSAGE=3
MAX_MESSAGE_IMAGE_SIZE=5242880  # 5MB
```

#### Frontend (docker-compose.yml or client env)

```env
PORT=:3000
BACKEND_URL=http://backend:8080  # Internal Docker networking
```

### OAuth Setup

To enable GitHub/Google OAuth:

1. **GitHub OAuth**:

   - Go to GitHub Settings → Developer settings → OAuth Apps
   - Create new OAuth App
   - Set Authorization callback URL: `http://localhost:8080/api/auth/github/callback`
   - Copy Client ID and Client Secret to backend `.env`

2. **Google OAuth**:

   - Go to <https://console.cloud.google.com/>
   - Create a project → APIs & Services → Credentials
   - Create OAuth 2.0 Client ID
   - Add Authorized redirect URI: `http://localhost:8080/api/auth/google/callback`
   - Copy Client ID and Client Secret to backend `.env`

## 📚 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/api/auth/register` | Register new user | No |
| `POST` | `/api/auth/login` | User login | No |
| `POST` | `/api/auth/logout` | User logout | Yes |
| `GET` | `/api/auth/me` | Get current user | Yes |
| `GET` | `/api/auth/github/login` | Initiate GitHub OAuth | No |
| `GET` | `/api/auth/github/callback` | GitHub OAuth callback | No |
| `GET` | `/api/auth/google/login` | Initiate Google OAuth | No |
| `GET` | `/api/auth/google/callback` | Google OAuth callback | No |

### Posts Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/posts` | Get all posts (paginated) | Yes |
| `GET` | `/api/posts/view/{id}` | Get single post | Yes |
| `POST` | `/api/posts/create` | Create new post | Yes |
| `PUT` | `/api/posts/edit/{id}` | Update post | Yes (owner) |
| `DELETE` | `/api/posts/remove/{id}` | Delete post | Yes (owner) |
| `GET` | `/api/posts/my-posts` | Get user's posts | Yes |
| `GET` | `/api/posts/liked-posts` | Get liked posts | Yes |
| `GET` | `/api/posts/category/{name}` | Get posts by category | Yes |

### Comments Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/comments/for-post/{id}` | Get comments for post | Yes |
| `POST` | `/api/comments/create-on-post/{id}` | Create comment | Yes |
| `PUT` | `/api/comments/edit/{id}` | Update comment | Yes (owner) |
| `DELETE` | `/api/comments/remove/{id}` | Delete comment | Yes (owner) |

### Reactions Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/api/reactions/posts/toggle` | Toggle post like/dislike | Yes |
| `POST` | `/api/reactions/comments/toggle` | Toggle comment like/dislike | Yes |

### Messages Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/messages/conversations` | Get all conversations | Yes |
| `GET` | `/api/messages/conversation/{userId}` | Get messages with user | Yes |
| `POST` | `/api/messages/send` | Send message | Yes |

### Notifications Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/notifications` | Get all notifications | Yes |
| `GET` | `/api/notifications/unread-count` | Get unread count | Yes |
| `PUT` | `/api/notifications/mark-read/{id}` | Mark as read | Yes |
| `PUT` | `/api/notifications/mark-all-read` | Mark all as read | Yes |

### Categories & Users

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/categories` | Get all categories | Yes |
| `GET` | `/api/users/online` | Get online users | Yes |
| `GET` | `/api/users/profile/{id}` | Get user profile | Yes |

### WebSocket

| Endpoint | Description | Auth Required |
|----------|-------------|---------------|
| `WS /ws` | WebSocket connection for real-time features | Yes (via session cookie) |

**WebSocket Message Types:**

- `typing_start` / `typing_stop` - Typing indicators
- `message` - New chat message
- `notification` - Real-time notification
- `user_online` / `user_offline` - Online status updates

## 🌐 Frontend Routes

| Route | View | Description |
|-------|------|-------------|
| `/` | HomeView | Latest posts, sorting, pagination |
| `/login` | LoginView | User authentication |
| `/register` | RegisterView | New user registration |
| `/post/:id` | PostView | Single post with comments |
| `/create-post` | CreatePostView | Create new post |
| `/edit-post/:id` | EditPostView | Edit existing post |
| `/profile/:id` | ProfileView | User profile & statistics |
| `/category/:name` | CategoryView | Posts filtered by category |
| `/chat` or `/chat/:userId` | ChatView | Real-time private messaging |
| `/notifications` | NotificationsView | Notification center |

## 🐳 Docker Deployment

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes (database reset)
docker-compose down -v

# Rebuild and restart
docker-compose up --build -d
```

### Architecture

- **3 containers**: `forum-db`, `forum-backend`, `forum-frontend`
- **1 network**: `forum-net` (internal communication)
- **1 volume**: `forum_db` (SQLite database persistence)

### Container Details

| Container | Port | Image | Description |
|-----------|------|-------|-------------|
| `forum-db` | - | `keinos/sqlite3:latest` | SQLite database container |
| `forum-backend` | 8080 | `real-time-forum_backend` | Go API server |
| `forum-frontend` | 3000 | `real-time-forum_frontend` | Go proxy + SPA |

## 🧪 Testing

### Manual Testing

1. **Register a new user**
   - Go to <http://localhost:3000>
   - Click "Register" and create an account

2. **Test OAuth**
   - Click "Sign in with GitHub" or "Sign in with Google"
   - Authorize the application

3. **Create content**
   - Create a post with categories
   - Add comments to posts
   - Like/dislike posts and comments

4. **Test real-time chat**
   - Open chat in two different browsers
   - Send messages and see live updates
   - Test typing indicators

5. **Test notifications**
   - Have another user comment on your post
   - Check the notification bell icon
   - Mark notifications as read

### API Testing with curl

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Create post
curl -X POST http://localhost:8080/api/posts/create \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "content": "My first post!",
    "category_names": ["Programming"]
  }'

# Get all posts
curl -X GET "http://localhost:8080/api/posts?page=1&per_page=10&sort=latest" \
  -b cookies.txt
```

## 🗄️ Database

- **Type**: SQLite3
- **Location**: `server/DBPath/forum.db` (persisted in Docker volume)
- **Auto-initialization**: Tables created automatically on first run
- **Categories**: Pre-populated with IT-focused categories

### Database Schema

- `users` - User accounts
- `sessions` - Active sessions
- `posts` - Forum posts
- `comments` - Post comments
- `categories` - Post categories
- `post_categories` - Many-to-many relationship
- `post_reactions` - Post likes/dislikes
- `comment_reactions` - Comment likes/dislikes
- `oauth_accounts` - OAuth provider linkage
- `oauth_states` - CSRF protection for OAuth
- `messages` - Private messages
- `notifications` - User notifications
- `post_images` - Uploaded post images
- `message_images` - Uploaded message images

## 🔒 Security Features

- **Password Hashing**: bcrypt with configurable cost
- **Session Management**: Secure HTTP-only cookies
- **CSRF Protection**: State validation for OAuth flows
- **Rate Limiting**: Configurable request throttling
- **XSS Protection**: HTML escaping on client-side
- **SQL Injection Prevention**: Prepared statements
- **CORS**: Configurable allowed origins
- **Security Headers**: Content-Type-Options, X-Frame-Options
- **Input Validation**: Server-side validation for all inputs

## 🛠️ Development

### Available Commands

```bash
make dev        # Start both backend and frontend
make backend    # Start only backend server
make frontend   # Start only frontend server
```

### Adding New Features

1. **Backend**:
   - Add handler in `server/internal/handlers/`
   - Create repository methods in `server/internal/repository/`
   - Add route in `server/internal/routes/routes.go`

2. **Frontend**:
   - Create view in `client/js/views/`
   - Add route in `client/js/router.js`
   - Add API call in `client/js/api/client.js`

### Code Quality

- **Go formatting**: `go fmt ./...`
- **Go linting**: `golangci-lint run` (if installed)
- **Browser logging**: Disabled in production, enabled in development

## 👥 Contributors

This project was developed by:

- **Kostas Apostolou** - [@kapostol]
- **Vicky Apostolou** - [@vapostol]
- **Dilhan Aslamaci** - [@daslamac]

## 🎓 About

This project was built as part of the Zone01 Athens curriculum.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

This project was created as part of the Zone01 Athens curriculum.

## 🎓 Learning Outcomes

This project demonstrates:

- Full-stack web development with Go
- RESTful API design
- WebSocket real-time communication
- SPA architecture with vanilla JavaScript
- Docker containerization
- OAuth 2.0 integration
- Database design and SQLite usage
- Security best practices
- Session-based authentication

---

**Built with ❤️ using Go, SQLite, WebSockets, CSS and Vanilla JavaScript**
