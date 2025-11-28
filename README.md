# Forum Application

A full-stack forum application built with Go, featuring a REST API backend and a web frontend with user authentication, posts, comments, and reactions.

## 🚀 Features

- **User Authentication**: Registration and login with session management
- **Posts Management**: Create, read, update, and delete forum posts
- **Comments System**: Add comments to posts with full CRUD operations
- **Reactions**: Like/dislike system for posts and comments
- **Categories**: Organize posts with category support
- **User Profiles**: View user statistics and activity
- **OAuth Integration**: GitHub and Google social login

## 📋 Prerequisites

- **Go**: Version 1.24 or higher
- **Docker**: For containerized deployment (optional)

## 🏗️ Project Structure

```
forum-project/
├── api/                           # Backend API Server
│   ├── main.go                    # API entry point
│   ├── config/                    # Configuration management
│   ├── database/                  # Database initialization
│   ├── internal/
│   │   ├── handlers/              # HTTP request handlers
│   │   ├── middleware/            # Authentication, CORS, rate limiting
│   │   ├── models/                # Data models
│   │   ├── repository/            # Data access layer
│   │   ├── routes/                # API route definitions
│   │   └── utils/                 # Utility functions
│   ├── Dockerfile                 # API container configuration
│   └── .env                       # API environment variables
├── frontend/                      # Frontend Web Server
│   ├── main.go                    # Frontend entry point
│   ├── config/                    # Frontend configuration
│   ├── internal/
│   │   ├── handlers/              # Web page handlers
│   │   ├── services/              # API client services
│   │   └── routes/                # Frontend route definitions
│   ├── web/
│   │   ├── templates/             # HTML templates
│   │   └── static/                # CSS, JS, images
│   ├── Dockerfile                 # Frontend container configuration
│   └── .env                       # Frontend environment variables
├── makefile                       # Development commands
└── README.md                      # This file
```

## 🛠️ Installation & Setup

### Option 1: Local Development

#### Start Both Servers (Recommended)
```bash
make dev
```

#### Start Servers Individually
```bash
# Backend API
make backend

# Frontend Web Server
make frontend

# Or run both 
make dev
```




**Access URLs:**
- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`

### Option 2: Docker Deployment

#### Backend API Container

1. **Create Docker network**
   ```bash
   docker network create forum-network
   ```

2. **Navigate to API directory**
   ```bash
   cd api
   ```

3. **Build API image**
   ```bash
   docker build -t api-image:latest .
   ```

4. **Run API container**
   ```bash
   docker run -d -p 8080:8080 -v api_db_data:/app/DBPath --name api-container --network forum-network api-image:latest
   ```

#### Frontend Container

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Build frontend image**
   ```bash
   docker build -t frontend-image:latest .
   ```

3. **Run frontend container**
   ```bash
   docker run -d -p 3000:3000 --name frontend-container --network forum-network frontend-image:latest
   ```

#### Docker Management Commands

**Check status:**
```bash
docker images
docker ps
```

**Stop and remove containers:**
```bash
# Stop containers
docker stop api-container frontend-container

# Remove containers
docker rm api-container frontend-container

# Remove images
docker rmi api-image:latest frontend-image:latest

# Complete cleanup
docker rm -f api-container frontend-container
docker rmi -f api-image:latest frontend-image:latest
```

**Network management:**
```bash
# Check network
docker network inspect forum-network

# Remove network
docker network rm forum-network
```

## ⚙️ Configuration

### Backend API (.env)
```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database Configuration
DB_PATH=./DBPath/forum.db

# Security Configuration
SESSION_DURATION=24h
BCRYPT_COST=16
SESSION_NAME=forum_session

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000
ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With,Cookie

# Content Limits
MIN_POST_CONTENT_LENGTH=10
MAX_POST_CONTENT_LENGTH=500
MIN_COMMENT_LENGTH=5
MAX_COMMENT_LENGTH=150
```

### Frontend (.env)
```env
# Frontend Server Configuration
FRONTEND_PORT=3000

# Backend API Configuration
API_BASE_URL=http://localhost:8080/api

# File Paths
TEMPLATES_DIR=./web/templates
STATIC_DIR=./web/static

# Session Configuration
SESSION_NAME=forum_session

# Environment
ENVIRONMENT=development
```

## 📚 API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `POST /api/auth/me` - Get current user

### Posts
- `GET /api/posts` - Get all posts
- `GET /api/posts/view/{id}` - Get single post
- `POST /api/posts/create` - Create new post
- `PUT /api/posts/edit/{id}` - Update post
- `DELETE /api/posts/remove/{id}` - Delete post

### Comments
- `GET /api/comments/for-post/{id}` - Get comments for post
- `POST /api/comments/create-on-post/{id}` - Create comment
- `PUT /api/comments/edit/{id}` - Update comment
- `DELETE /api/comments/remove/{id}` - Delete comment

### Reactions
- `POST /api/reactions/posts/toggle` - Toggle post reaction
- `POST /api/reactions/comments/toggle` - Toggle comment reaction

### Categories
- `GET /api/categories` - Get all categories

## 🌐 Frontend Pages

- **Home** (`/`) - Latest posts with sorting and pagination
- **Login** (`/login`) - User authentication
- **Register** (`/register`) - New user registration
- **Post View** (`/post/{id}`) - Individual post with comments
- **Create Post** (`/create-post`) - New post creation
- **Edit Post** (`/edit-post/{id}`) - Post editing
- **User Profile** (`/profile`) - User statistics and activity
- **Category View** (`/category/{id}`) - Posts by category

## 🧪 Testing

### API Testing with curl

```bash
# Register a new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "TestPass123!",
    "confirm_password": "TestPass123!"
  }'

# Login and save session
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123!"
  }'

# Create a post
curl -X POST http://localhost:8080/api/posts/create \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "content": "This is my first post!",
    "category_names": ["General discussion"]
  }'
```

### Frontend Testing

1. Navigate to `http://localhost:3000`
2. Register a new account
3. Login with your credentials
4. Create posts and comments
5. Test reactions and user profiles

## 🗄️ Database

- **Database**: SQLite with automatic schema creation
- **Location**: `./DBPath/forum.db` (API directory)
- **Tables**: users, sessions, posts, comments, categories, reactions


## 🔧 Development Commands

```bash
# Start both servers
make dev

# Start only backend
make backend

# Start only frontend  
make frontend
```

## 🐳 Docker Features

- **Multi-stage builds** for optimized image sizes
- **Non-root user** execution for security
- **Volume persistence** for database data
- **Alpine Linux** base for minimal attack surface

---

**Built with Go, SQLite, and Docker**