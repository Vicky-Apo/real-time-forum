# Real-Time Forum: Implementation Plan
## Focus: Meet Exercise Requirements First

This is a focused, step-by-step plan to meet the 01 Edu exercise requirements before adding DevOps layers.

---

## Current Status

### ✅ Backend - Ready
- API endpoints complete (auth, posts, comments, messages)
- WebSocket server implemented ([server/internal/websocket/hub.go](server/internal/websocket/hub.go))
- Database schema mostly ready
- OAuth working (GitHub, Google)

### ❌ Frontend - Needs Complete Rewrite
- **Problem:** 13 HTML templates (multi-page app)
- **Requirement:** 1 HTML file (single-page app)
- **Missing:** WebSocket client, Chat UI

### ⚠️ Backend - Minor Updates Needed
- Add new user registration fields (age, gender, first name, last name)

---

## Implementation Phases

## Phase 1: Project Structure Setup (30 minutes)

### 1.1 Create New Client Directory
```bash
# Remove old frontend (backup first)
mv frontend frontend-old

# Create new SPA structure
mkdir -p client/{js/{api,views,components,websocket,utils},css,assets}
```

### 1.2 Directory Structure
```
client/
├── index.html                 # Single HTML file
├── js/
│   ├── main.js               # Application entry point
│   ├── router.js             # Client-side routing
│   ├── state.js              # Global state management
│   ├── api/                  # Backend API client
│   │   ├── client.js         # Base HTTP client
│   │   ├── auth.js           # Authentication API
│   │   ├── posts.js          # Posts API
│   │   ├── comments.js       # Comments API
│   │   ├── messages.js       # Messages API
│   │   └── notifications.js  # Notifications API
│   ├── views/                # Page views (components)
│   │   ├── LoginView.js
│   │   ├── RegisterView.js
│   │   ├── HomeView.js       # Feed/posts list
│   │   ├── PostView.js       # Single post with comments
│   │   ├── CreatePostView.js
│   │   ├── EditPostView.js
│   │   ├── ProfileView.js
│   │   ├── ChatView.js       # NEW: Chat interface
│   │   └── NotificationsView.js
│   ├── components/           # Reusable UI components
│   │   ├── Navbar.js
│   │   ├── PostCard.js
│   │   ├── CommentForm.js
│   │   ├── CommentList.js
│   │   ├── UserList.js       # NEW: Online users
│   │   ├── MessageThread.js  # NEW: Chat messages
│   │   └── TypingIndicator.js # NEW
│   ├── websocket/            # WebSocket client
│   │   ├── WebSocketManager.js
│   │   └── handlers.js       # WebSocket event handlers
│   └── utils/
│       ├── dom.js            # DOM helpers
│       ├── validation.js     # Form validation
│       └── formatting.js     # Date, text formatting
├── css/
│   ├── global.css            # Copy from old frontend
│   ├── components.css        # Component-specific styles
│   ├── views.css             # View-specific styles
│   └── chat.css              # NEW: Chat interface styles
└── assets/
    └── images/               # Static images
```

### 1.3 Setup Tasks
- [ ] Create directory structure
- [ ] Copy existing CSS files from old frontend
- [ ] Plan component reusability
- [ ] Document module dependencies

---

## Phase 2: Core Infrastructure (2-3 hours)

### 2.1 Single HTML File (index.html)
**Goal:** One HTML file that serves the entire application

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Real-Time Forum</title>
    <link rel="stylesheet" href="/css/global.css">
    <link rel="stylesheet" href="/css/components.css">
    <link rel="stylesheet" href="/css/views.css">
    <link rel="stylesheet" href="/css/chat.css">
</head>
<body>
    <!-- Navbar (always visible) -->
    <div id="navbar"></div>

    <!-- Main content area (dynamically replaced) -->
    <div id="app"></div>

    <!-- Modals/overlays -->
    <div id="modal-root"></div>

    <!-- Toast notifications -->
    <div id="toast-container"></div>

    <!-- Load app as ES6 module -->
    <script type="module" src="/js/main.js"></script>
</body>
</html>
```

**Key Points:**
- Single mounting point: `#app`
- All page changes happen by replacing `#app` content
- Navbar persists across navigation
- No page reloads

### 2.2 Client-Side Router (router.js)
**Goal:** Handle navigation without page refreshes

```javascript
// js/router.js
class Router {
    constructor(routes) {
        this.routes = routes;
        this.currentRoute = null;

        // Handle browser back/forward
        window.addEventListener('popstate', () => this.handleRoute());

        // Intercept link clicks
        document.addEventListener('click', (e) => {
            if (e.target.matches('[data-link]')) {
                e.preventDefault();
                this.navigate(e.target.href);
            }
        });
    }

    navigate(path) {
        window.history.pushState(null, null, path);
        this.handleRoute();
    }

    async handleRoute() {
        const path = window.location.pathname;
        const route = this.matchRoute(path);

        if (!route) {
            // 404 - show home or error
            this.navigate('/');
            return;
        }

        // Check authentication
        if (route.requiresAuth && !state.user) {
            this.navigate('/login');
            return;
        }

        // Render the view
        this.currentRoute = route;
        await this.renderView(route, path);
    }

    matchRoute(path) {
        for (const route of this.routes) {
            const match = this.pathToRegex(route.path).exec(path);
            if (match) {
                // Extract params (e.g., /post/:id)
                const params = this.extractParams(route.path, match);
                return { ...route, params };
            }
        }
        return null;
    }

    pathToRegex(path) {
        // Convert /post/:id to regex: /post/([^/]+)
        return new RegExp('^' + path.replace(/:\w+/g, '([^/]+)') + '$');
    }

    extractParams(path, match) {
        const keys = [...path.matchAll(/:(\w+)/g)].map(m => m[1]);
        return keys.reduce((params, key, i) => {
            params[key] = match[i + 1];
            return params;
        }, {});
    }

    async renderView(route, path) {
        const app = document.getElementById('app');
        app.innerHTML = '<div class="loading">Loading...</div>';

        try {
            // Import and render the view
            const view = await route.component();
            app.innerHTML = await view.render(route.params);

            // Run view's afterRender (event listeners, etc.)
            if (view.afterRender) {
                view.afterRender(route.params);
            }
        } catch (error) {
            console.error('Error rendering view:', error);
            app.innerHTML = '<div class="error">Something went wrong</div>';
        }
    }
}

export default Router;
```

### 2.3 State Management (state.js)
**Goal:** Global application state

```javascript
// js/state.js
class State {
    constructor() {
        this.user = null;
        this.wsConnected = false;
        this.onlineUsers = [];
        this.unreadCount = 0;
        this.listeners = new Map();
    }

    // Get current user
    getUser() {
        if (!this.user) {
            // Try to load from localStorage
            const stored = localStorage.getItem('user');
            if (stored) {
                this.user = JSON.parse(stored);
            }
        }
        return this.user;
    }

    // Set current user
    setUser(user) {
        this.user = user;
        if (user) {
            localStorage.setItem('user', JSON.stringify(user));
        } else {
            localStorage.removeItem('user');
        }
        this.emit('user:changed', user);
    }

    // Online users
    setOnlineUsers(users) {
        this.onlineUsers = users;
        this.emit('users:online', users);
    }

    addOnlineUser(user) {
        if (!this.onlineUsers.find(u => u.id === user.id)) {
            this.onlineUsers.push(user);
            this.emit('user:online', user);
        }
    }

    removeOnlineUser(userId) {
        this.onlineUsers = this.onlineUsers.filter(u => u.id !== userId);
        this.emit('user:offline', userId);
    }

    // Unread messages
    setUnreadCount(count) {
        this.unreadCount = count;
        this.emit('unread:changed', count);
    }

    // Event system (pub/sub)
    on(event, callback) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, []);
        }
        this.listeners.get(event).push(callback);
    }

    off(event, callback) {
        const callbacks = this.listeners.get(event);
        if (callbacks) {
            this.listeners.set(event, callbacks.filter(cb => cb !== callback));
        }
    }

    emit(event, data) {
        const callbacks = this.listeners.get(event);
        if (callbacks) {
            callbacks.forEach(cb => cb(data));
        }
    }
}

const state = new State();
export default state;
```

### 2.4 Base API Client (api/client.js)
**Goal:** Handle HTTP requests to backend

```javascript
// js/api/client.js
class APIClient {
    constructor(baseURL) {
        this.baseURL = baseURL || 'http://localhost:8080/api';
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;

        const config = {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            credentials: 'include', // Include cookies
        };

        try {
            const response = await fetch(url, config);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Request failed');
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    get(endpoint) {
        return this.request(endpoint, { method: 'GET' });
    }

    post(endpoint, body) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(body),
        });
    }

    put(endpoint, body) {
        return this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(body),
        });
    }

    delete(endpoint) {
        return this.request(endpoint, { method: 'DELETE' });
    }
}

const apiClient = new APIClient();
export default apiClient;
```

### 2.5 Tasks
- [ ] Create `index.html`
- [ ] Implement `router.js`
- [ ] Implement `state.js`
- [ ] Implement `api/client.js`
- [ ] Test routing with placeholder views

---

## Phase 3: WebSocket Client (2-3 hours)

### 3.1 WebSocket Manager (websocket/WebSocketManager.js)
**Goal:** Manage WebSocket connection with reconnection logic

```javascript
// js/websocket/WebSocketManager.js
import state from '../state.js';
import { handleWebSocketMessage } from './handlers.js';

class WebSocketManager {
    constructor(url) {
        this.url = url || 'ws://localhost:8080/ws';
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 10;
        this.reconnectDelay = 1000; // Start with 1s
        this.isConnecting = false;
        this.shouldReconnect = true;
    }

    connect() {
        if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
            return;
        }

        this.isConnecting = true;
        console.log('[WS] Connecting to', this.url);

        try {
            this.ws = new WebSocket(this.url);
            this.setupEventHandlers();
        } catch (error) {
            console.error('[WS] Connection error:', error);
            this.isConnecting = false;
            this.scheduleReconnect();
        }
    }

    setupEventHandlers() {
        this.ws.onopen = () => this.onOpen();
        this.ws.onmessage = (event) => this.onMessage(event);
        this.ws.onerror = (error) => this.onError(error);
        this.ws.onclose = (event) => this.onClose(event);
    }

    onOpen() {
        console.log('[WS] Connected');
        this.isConnecting = false;
        this.reconnectAttempts = 0;
        this.reconnectDelay = 1000;

        state.wsConnected = true;
        state.emit('ws:connected');
    }

    onMessage(event) {
        try {
            const message = JSON.parse(event.data);
            console.log('[WS] Received:', message);

            // Delegate to handler
            handleWebSocketMessage(message);
        } catch (error) {
            console.error('[WS] Failed to parse message:', error);
        }
    }

    onError(error) {
        console.error('[WS] Error:', error);
    }

    onClose(event) {
        console.log('[WS] Disconnected', event.code, event.reason);
        this.isConnecting = false;

        state.wsConnected = false;
        state.emit('ws:disconnected');

        if (this.shouldReconnect) {
            this.scheduleReconnect();
        }
    }

    scheduleReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('[WS] Max reconnection attempts reached');
            state.emit('ws:failed');
            return;
        }

        this.reconnectAttempts++;
        const delay = Math.min(this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1), 30000);

        console.log(`[WS] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);

        setTimeout(() => {
            this.connect();
        }, delay);
    }

    send(type, payload) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message = JSON.stringify({ type, payload });
            this.ws.send(message);
            console.log('[WS] Sent:', { type, payload });
        } else {
            console.warn('[WS] Cannot send, not connected');
        }
    }

    disconnect() {
        this.shouldReconnect = false;
        if (this.ws) {
            this.ws.close();
        }
    }
}

const wsManager = new WebSocketManager();
export default wsManager;
```

### 3.2 WebSocket Event Handlers (websocket/handlers.js)
**Goal:** Handle incoming WebSocket messages

```javascript
// js/websocket/handlers.js
import state from '../state.js';

export function handleWebSocketMessage(message) {
    const { type, payload } = message;

    switch (type) {
        case 'user_online':
            handleUserOnline(payload);
            break;

        case 'user_offline':
            handleUserOffline(payload);
            break;

        case 'typing_start':
            handleTypingStart(payload);
            break;

        case 'typing_stop':
            handleTypingStop(payload);
            break;

        case 'new_message':
            handleNewMessage(payload);
            break;

        case 'message_read':
            handleMessageRead(payload);
            break;

        case 'notification':
            handleNotification(payload);
            break;

        default:
            console.warn('[WS] Unknown message type:', type);
    }
}

function handleUserOnline(user) {
    console.log('[WS] User online:', user);
    state.addOnlineUser(user);
}

function handleUserOffline(userId) {
    console.log('[WS] User offline:', userId);
    state.removeOnlineUser(userId);
}

function handleTypingStart({ userId, userName }) {
    console.log('[WS] User typing:', userName);
    state.emit('typing:start', { userId, userName });
}

function handleTypingStop({ userId }) {
    console.log('[WS] User stopped typing:', userId);
    state.emit('typing:stop', { userId });
}

function handleNewMessage(message) {
    console.log('[WS] New message:', message);
    state.emit('message:received', message);

    // Increment unread count if not on chat view
    if (window.location.pathname !== '/chat') {
        state.setUnreadCount(state.unreadCount + 1);
    }

    // Show browser notification
    showBrowserNotification(message);
}

function handleMessageRead({ messageId }) {
    state.emit('message:read', { messageId });
}

function handleNotification(notification) {
    console.log('[WS] Notification:', notification);
    state.emit('notification:received', notification);
}

function showBrowserNotification(message) {
    if ('Notification' in window && Notification.permission === 'granted') {
        new Notification('New Message', {
            body: `${message.senderName}: ${message.content}`,
            icon: '/assets/logo.png',
        });
    }
}
```

### 3.3 Tasks
- [ ] Implement WebSocket manager
- [ ] Implement message handlers
- [ ] Test connection/reconnection
- [ ] Test event handling
- [ ] Integrate with state management

---

## Phase 4: Authentication Views (3-4 hours)

### 4.1 Update Backend for New User Fields

#### Step 1: Database Migration
```sql
-- migrations/000002_add_user_profile_fields.up.sql
ALTER TABLE users ADD COLUMN age INTEGER;
ALTER TABLE users ADD COLUMN gender TEXT;
ALTER TABLE users ADD COLUMN first_name TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN last_name TEXT NOT NULL DEFAULT '';

-- Add constraints
CREATE INDEX idx_users_first_name ON users(first_name);
CREATE INDEX idx_users_last_name ON users(last_name);
```

**OR** (if not using migrations, update existing schema):

```go
// server/database/sql_statements.go
const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    age INTEGER,
    gender TEXT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`
```

#### Step 2: Update Models
```go
// server/internal/models/users_models.go
type UserRegistration struct {
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  string `json:"password"`
    Age       int    `json:"age"`        // NEW
    Gender    string `json:"gender"`     // NEW
    FirstName string `json:"first_name"` // NEW
    LastName  string `json:"last_name"`  // NEW
}

type User struct {
    UserID    string    `json:"user_id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Age       int       `json:"age"`        // NEW
    Gender    string    `json:"gender"`     // NEW
    FirstName string    `json:"first_name"` // NEW
    LastName  string    `json:"last_name"`  // NEW
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### Step 3: Update Validation
```go
// server/internal/utils/validate.go
func ValidateUserRegistration(user *models.UserRegistration) error {
    // Existing validations...

    // Age validation
    if user.Age < 13 || user.Age > 120 {
        return errors.New("age must be between 13 and 120")
    }

    // Gender validation
    validGenders := []string{"male", "female", "other", "prefer-not-to-say"}
    if !contains(validGenders, strings.ToLower(user.Gender)) {
        return errors.New("invalid gender")
    }

    // Name validations
    if len(user.FirstName) < 1 || len(user.FirstName) > 50 {
        return errors.New("first name must be between 1 and 50 characters")
    }
    if len(user.LastName) < 1 || len(user.LastName) > 50 {
        return errors.New("last name must be between 1 and 50 characters")
    }

    return nil
}
```

#### Step 4: Update Repository
```go
// server/internal/repository/user_repository.go
func (r *UserRepository) CreateUser(ctx context.Context, user *models.UserRegistration) (string, error) {
    query := `
        INSERT INTO users (user_id, username, email, password, age, gender, first_name, last_name)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

    userID := uuid.New().String()

    _, err := r.db.ExecContext(ctx, query,
        userID,
        user.Username,
        user.Email,
        user.Password,
        user.Age,
        user.Gender,
        user.FirstName,
        user.LastName,
    )

    return userID, err
}
```

### 4.2 Registration View (Frontend)

```javascript
// js/views/RegisterView.js
import apiClient from '../api/client.js';
import { navigate } from '../router.js';

export default {
    async render() {
        return `
            <div class="auth-container">
                <div class="auth-card">
                    <h1>Create Account</h1>
                    <form id="register-form" class="auth-form">
                        <div class="form-row">
                            <div class="form-group">
                                <label for="first-name">First Name *</label>
                                <input type="text" id="first-name" name="first_name" required>
                            </div>
                            <div class="form-group">
                                <label for="last-name">Last Name *</label>
                                <input type="text" id="last-name" name="last_name" required>
                            </div>
                        </div>

                        <div class="form-group">
                            <label for="username">Username *</label>
                            <input type="text" id="username" name="username" required
                                   minlength="3" maxlength="20">
                        </div>

                        <div class="form-group">
                            <label for="email">Email *</label>
                            <input type="email" id="email" name="email" required>
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label for="age">Age *</label>
                                <input type="number" id="age" name="age" required
                                       min="13" max="120">
                            </div>
                            <div class="form-group">
                                <label for="gender">Gender *</label>
                                <select id="gender" name="gender" required>
                                    <option value="">Select...</option>
                                    <option value="male">Male</option>
                                    <option value="female">Female</option>
                                    <option value="other">Other</option>
                                    <option value="prefer-not-to-say">Prefer not to say</option>
                                </select>
                            </div>
                        </div>

                        <div class="form-group">
                            <label for="password">Password *</label>
                            <input type="password" id="password" name="password" required
                                   minlength="8">
                        </div>

                        <div class="form-group">
                            <label for="confirm-password">Confirm Password *</label>
                            <input type="password" id="confirm-password" name="confirm_password" required>
                        </div>

                        <div id="error-message" class="error-message"></div>

                        <button type="submit" class="btn btn-primary">Register</button>
                    </form>

                    <p class="auth-footer">
                        Already have an account?
                        <a href="/login" data-link>Login</a>
                    </p>
                </div>
            </div>
        `;
    },

    afterRender() {
        const form = document.getElementById('register-form');
        form.addEventListener('submit', this.handleSubmit);
    },

    async handleSubmit(e) {
        e.preventDefault();

        const formData = new FormData(e.target);
        const data = Object.fromEntries(formData);

        // Validate passwords match
        if (data.password !== data.confirm_password) {
            showError('Passwords do not match');
            return;
        }

        // Convert age to number
        data.age = parseInt(data.age);

        try {
            await apiClient.post('/auth/register', data);

            // Auto-login after registration
            await apiClient.post('/auth/login', {
                identifier: data.username,
                password: data.password,
            });

            // Redirect to home
            navigate('/');
        } catch (error) {
            showError(error.message);
        }
    }
};

function showError(message) {
    const errorDiv = document.getElementById('error-message');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
}
```

### 4.3 Login View
```javascript
// js/views/LoginView.js
import apiClient from '../api/client.js';
import state from '../state.js';
import { navigate } from '../router.js';
import wsManager from '../websocket/WebSocketManager.js';

export default {
    async render() {
        return `
            <div class="auth-container">
                <div class="auth-card">
                    <h1>Welcome Back</h1>
                    <form id="login-form" class="auth-form">
                        <div class="form-group">
                            <label for="identifier">Username or Email</label>
                            <input type="text" id="identifier" name="identifier" required>
                        </div>

                        <div class="form-group">
                            <label for="password">Password</label>
                            <input type="password" id="password" name="password" required>
                        </div>

                        <div id="error-message" class="error-message"></div>

                        <button type="submit" class="btn btn-primary">Login</button>
                    </form>

                    <div class="auth-divider">OR</div>

                    <div class="oauth-buttons">
                        <a href="/api/auth/github/login" class="btn btn-github">
                            Login with GitHub
                        </a>
                        <a href="/api/auth/google/login" class="btn btn-google">
                            Login with Google
                        </a>
                    </div>

                    <p class="auth-footer">
                        Don't have an account?
                        <a href="/register" data-link>Register</a>
                    </p>
                </div>
            </div>
        `;
    },

    afterRender() {
        const form = document.getElementById('login-form');
        form.addEventListener('submit', this.handleSubmit);
    },

    async handleSubmit(e) {
        e.preventDefault();

        const formData = new FormData(e.target);
        const data = Object.fromEntries(formData);

        try {
            const response = await apiClient.post('/auth/login', data);

            // Store user in state
            state.setUser(response.user);

            // Connect WebSocket
            wsManager.connect();

            // Redirect to home
            navigate('/');
        } catch (error) {
            showError(error.message);
        }
    }
};

function showError(message) {
    const errorDiv = document.getElementById('error-message');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
}
```

### 4.4 Tasks
- [ ] Update backend schema for new user fields
- [ ] Update models, validation, repository
- [ ] Create `RegisterView.js`
- [ ] Create `LoginView.js`
- [ ] Test registration flow
- [ ] Test login flow

---

## Phase 5: Chat Interface (6-8 hours)

This is the most complex part. The chat UI has three main components:

### 5.1 Chat View Layout
```
┌─────────────────────────────────────────┐
│           Navbar (always visible)       │
├───────────┬─────────────────────────────┤
│           │  Message Thread             │
│   User    │  ┌───────────────────────┐  │
│   List    │  │ Message 1              │  │
│           │  │ Message 2              │  │
│   Online  │  │ ...                    │  │
│   Users   │  └───────────────────────┘  │
│           │                             │
│   Filter  │  Typing: User is typing...  │
│   Search  │                             │
│           │  ┌───────────────────────┐  │
│           │  │ Message input         │  │
│           │  │ [Send]                │  │
│           │  └───────────────────────┘  │
└───────────┴─────────────────────────────┘
```

### 5.2 User List Component
```javascript
// js/components/UserList.js
import state from '../state.js';
import apiClient from '../api/client.js';

export default class UserList {
    constructor(onUserSelect) {
        this.onUserSelect = onUserSelect;
        this.users = [];
        this.conversations = [];
        this.filter = '';

        // Listen to state changes
        state.on('users:online', (users) => this.updateOnlineStatus(users));
        state.on('user:online', (user) => this.handleUserOnline(user));
        state.on('user:offline', (userId) => this.handleUserOffline(userId));
    }

    async init() {
        // Load conversations (users with message history)
        this.conversations = await apiClient.get('/conversations');

        // Load all online users
        const onlineUsers = await apiClient.get('/users/online');
        state.setOnlineUsers(onlineUsers);

        this.render();
    }

    render() {
        const container = document.getElementById('user-list');

        // Combine conversations and online users
        const allUsers = this.mergeUsers();

        // Sort: conversations by last message, then alphabetically
        const sorted = this.sortUsers(allUsers);

        // Filter
        const filtered = this.filter
            ? sorted.filter(u => u.username.toLowerCase().includes(this.filter.toLowerCase()))
            : sorted;

        container.innerHTML = `
            <div class="user-list-header">
                <h2>Messages</h2>
                <input
                    type="text"
                    id="user-search"
                    placeholder="Search users..."
                    value="${this.filter}"
                >
            </div>
            <div class="user-list-items">
                ${filtered.map(user => this.renderUserItem(user)).join('')}
            </div>
        `;

        // Event listeners
        document.getElementById('user-search').addEventListener('input', (e) => {
            this.filter = e.target.value;
            this.render();
        });

        document.querySelectorAll('.user-item').forEach(item => {
            item.addEventListener('click', () => {
                const userId = item.dataset.userId;
                this.onUserSelect(userId);
            });
        });
    }

    renderUserItem(user) {
        const isOnline = state.onlineUsers.some(u => u.user_id === user.user_id);
        const statusClass = isOnline ? 'online' : 'offline';
        const unreadBadge = user.unread_count > 0
            ? `<span class="unread-badge">${user.unread_count}</span>`
            : '';

        return `
            <div class="user-item ${statusClass}" data-user-id="${user.user_id}">
                <div class="user-avatar">
                    <span class="status-dot"></span>
                </div>
                <div class="user-info">
                    <div class="user-name">${user.username}</div>
                    <div class="last-message">${user.last_message || 'Start a conversation'}</div>
                </div>
                ${unreadBadge}
            </div>
        `;
    }

    mergeUsers() {
        // Merge conversations with online users
        const userMap = new Map();

        // Add conversations
        this.conversations.forEach(conv => {
            userMap.set(conv.user_id, conv);
        });

        // Add online users not in conversations
        state.onlineUsers.forEach(user => {
            if (!userMap.has(user.user_id)) {
                userMap.set(user.user_id, {
                    ...user,
                    last_message: null,
                    last_message_time: null,
                    unread_count: 0,
                });
            }
        });

        return Array.from(userMap.values());
    }

    sortUsers(users) {
        return users.sort((a, b) => {
            // Users with conversations come first
            if (a.last_message_time && !b.last_message_time) return -1;
            if (!a.last_message_time && b.last_message_time) return 1;

            // Sort by last message time
            if (a.last_message_time && b.last_message_time) {
                return new Date(b.last_message_time) - new Date(a.last_message_time);
            }

            // Alphabetical for users without conversations
            return a.username.localeCompare(b.username);
        });
    }

    updateOnlineStatus(users) {
        this.render();
    }

    handleUserOnline(user) {
        this.render();
    }

    handleUserOffline(userId) {
        this.render();
    }
}
```

### 5.3 Message Thread Component
```javascript
// js/components/MessageThread.js
import state from '../state.js';
import apiClient from '../api/client.js';
import wsManager from '../websocket/WebSocketManager.js';
import { formatDate, formatTime } from '../utils/formatting.js';
import { throttle } from '../utils/helpers.js';

export default class MessageThread {
    constructor(userId) {
        this.userId = userId;
        this.messages = [];
        this.offset = 0;
        this.limit = 10;
        this.hasMore = true;
        this.isLoading = false;
        this.typingTimeout = null;

        // Listen to WebSocket events
        state.on('message:received', (msg) => this.handleNewMessage(msg));
        state.on('typing:start', (data) => this.handleTypingStart(data));
        state.on('typing:stop', (data) => this.handleTypingStop(data));
    }

    async init() {
        await this.loadMessages();
        this.render();
        this.scrollToBottom();
    }

    async loadMessages(prepend = false) {
        if (this.isLoading || !this.hasMore) return;

        this.isLoading = true;

        try {
            const messages = await apiClient.get(`/messages/${this.userId}?limit=${this.limit}&offset=${this.offset}`);

            if (messages.length < this.limit) {
                this.hasMore = false;
            }

            if (prepend) {
                this.messages = [...messages.reverse(), ...this.messages];
            } else {
                this.messages = messages.reverse();
            }

            this.offset += messages.length;
            this.isLoading = false;

            return messages;
        } catch (error) {
            console.error('Failed to load messages:', error);
            this.isLoading = false;
        }
    }

    render() {
        const container = document.getElementById('message-thread');

        const recipient = state.onlineUsers.find(u => u.user_id === this.userId);
        const recipientName = recipient?.username || 'User';
        const isOnline = !!recipient;

        container.innerHTML = `
            <div class="thread-header">
                <div class="recipient-info">
                    <h2>${recipientName}</h2>
                    <span class="status">${isOnline ? 'Online' : 'Offline'}</span>
                </div>
            </div>

            <div id="messages-container" class="messages-container">
                ${this.messages.map(msg => this.renderMessage(msg)).join('')}
            </div>

            <div id="typing-indicator" class="typing-indicator" style="display: none;">
                <span class="typing-dots"></span>
                <span id="typing-text"></span>
            </div>

            <div class="message-input-container">
                <textarea
                    id="message-input"
                    placeholder="Type a message..."
                    rows="1"
                ></textarea>
                <button id="send-btn" class="btn btn-primary">Send</button>
            </div>
        `;

        this.attachEventListeners();
    }

    renderMessage(message) {
        const isMine = message.sender_id === state.getUser().user_id;
        const messageClass = isMine ? 'message-mine' : 'message-theirs';

        return `
            <div class="message ${messageClass}" data-message-id="${message.message_id}">
                <div class="message-content">${message.content}</div>
                <div class="message-meta">
                    <span class="message-time">${formatTime(message.created_at)}</span>
                    ${isMine && message.is_read ? '<span class="read-receipt">✓✓</span>' : ''}
                </div>
            </div>
        `;
    }

    attachEventListeners() {
        const input = document.getElementById('message-input');
        const sendBtn = document.getElementById('send-btn');
        const messagesContainer = document.getElementById('messages-container');

        // Send message
        sendBtn.addEventListener('click', () => this.sendMessage());
        input.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });

        // Typing indicator
        input.addEventListener('input', () => this.handleTyping());

        // Auto-expand textarea
        input.addEventListener('input', () => {
            input.style.height = 'auto';
            input.style.height = input.scrollHeight + 'px';
        });

        // Load more messages on scroll up
        messagesContainer.addEventListener('scroll', throttle((e) => {
            if (e.target.scrollTop < 100 && this.hasMore) {
                this.loadMoreMessages();
            }
        }, 500));
    }

    async sendMessage() {
        const input = document.getElementById('message-input');
        const content = input.value.trim();

        if (!content) return;

        try {
            const message = await apiClient.post('/messages/send', {
                recipient_id: this.userId,
                content: content,
            });

            // Add to messages
            this.messages.push(message);

            // Clear input
            input.value = '';
            input.style.height = 'auto';

            // Re-render and scroll
            this.render();
            this.scrollToBottom();

        } catch (error) {
            console.error('Failed to send message:', error);
            alert('Failed to send message');
        }
    }

    handleTyping() {
        // Send typing start
        wsManager.send('typing_start', { recipient_id: this.userId });

        // Clear existing timeout
        if (this.typingTimeout) {
            clearTimeout(this.typingTimeout);
        }

        // Send typing stop after 2 seconds of inactivity
        this.typingTimeout = setTimeout(() => {
            wsManager.send('typing_stop', { recipient_id: this.userId });
        }, 2000);
    }

    handleNewMessage(message) {
        // Only handle messages for this conversation
        if (message.sender_id === this.userId || message.recipient_id === this.userId) {
            this.messages.push(message);
            this.render();
            this.scrollToBottom();

            // Mark as read
            if (message.sender_id === this.userId) {
                apiClient.post(`/messages/mark-read/${message.message_id}`);
            }
        }
    }

    handleTypingStart({ userId, userName }) {
        if (userId === this.userId) {
            const indicator = document.getElementById('typing-indicator');
            const text = document.getElementById('typing-text');
            text.textContent = `${userName} is typing...`;
            indicator.style.display = 'block';
        }
    }

    handleTypingStop({ userId }) {
        if (userId === this.userId) {
            const indicator = document.getElementById('typing-indicator');
            indicator.style.display = 'none';
        }
    }

    async loadMoreMessages() {
        const container = document.getElementById('messages-container');
        const oldScrollHeight = container.scrollHeight;

        await this.loadMessages(true);
        this.render();

        // Maintain scroll position
        container.scrollTop = container.scrollHeight - oldScrollHeight;
    }

    scrollToBottom() {
        const container = document.getElementById('messages-container');
        container.scrollTop = container.scrollHeight;
    }
}
```

### 5.4 Chat View (Main)
```javascript
// js/views/ChatView.js
import UserList from '../components/UserList.js';
import MessageThread from '../components/MessageThread.js';

export default {
    currentThread: null,
    userList: null,

    async render() {
        return `
            <div class="chat-container">
                <div id="user-list" class="user-list-panel"></div>
                <div id="message-thread" class="message-thread-panel">
                    <div class="no-conversation">
                        <p>Select a user to start chatting</p>
                    </div>
                </div>
            </div>
        `;
    },

    async afterRender() {
        // Initialize user list
        this.userList = new UserList((userId) => this.openConversation(userId));
        await this.userList.init();
    },

    async openConversation(userId) {
        // Clean up previous thread
        if (this.currentThread) {
            // Remove event listeners
        }

        // Create new thread
        this.currentThread = new MessageThread(userId);
        await this.currentThread.init();
    }
};
```

### 5.5 Tasks
- [ ] Create `UserList.js` component
- [ ] Create `MessageThread.js` component
- [ ] Create `ChatView.js`
- [ ] Implement message pagination (scroll up for more)
- [ ] Implement throttle/debounce utilities
- [ ] Style chat interface
- [ ] Test sending/receiving messages
- [ ] Test typing indicators
- [ ] Test online/offline status

---

## Phase 6: Posts & Comments Views (4-6 hours)

### 6.1 Views to Create
- `HomeView.js` - Feed with posts list
- `PostView.js` - Single post with comments
- `CreatePostView.js` - Create new post
- `EditPostView.js` - Edit existing post
- `ProfileView.js` - User profile

### 6.2 Components to Create
- `PostCard.js` - Reusable post display
- `CommentList.js` - List of comments
- `CommentForm.js` - Create/edit comment

These are simpler than chat, so I won't detail them fully here. Follow the same pattern as the chat components.

### 6.3 Tasks
- [ ] Create all view files
- [ ] Create all component files
- [ ] Implement post creation with image upload
- [ ] Implement commenting
- [ ] Implement reactions (like/dislike)
- [ ] Test all functionality

---

## Phase 7: Static File Server (30 minutes)

### 7.1 Simple Nginx Configuration
```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    # SPA routing: serve index.html for all routes
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Proxy API requests to backend
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Proxy WebSocket
    location /ws {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }
}
```

### 7.2 Frontend Dockerfile
```dockerfile
FROM nginx:alpine

# Copy static files
COPY client/ /usr/share/nginx/html/

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### 7.3 Update docker-compose.yml
```yaml
version: "3.9"
services:
  backend:
    build: ./server
    ports:
      - "8080:8080"
    environment:
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8080
      - ALLOWED_ORIGINS=http://localhost:3000
    volumes:
      - ./forum.db:/app/forum.db
      - ./uploads:/app/uploads

  frontend:
    build: .
    dockerfile: Dockerfile.frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
```

### 7.4 Tasks
- [ ] Create nginx.conf
- [ ] Create Dockerfile.frontend
- [ ] Update docker-compose.yml
- [ ] Test: `docker-compose up`

---

## Phase 8: Testing & Polish (2-3 hours)

### 8.1 Manual Testing Checklist
- [ ] Registration with all fields
- [ ] Login with username
- [ ] Login with email
- [ ] OAuth login (GitHub, Google)
- [ ] Create post
- [ ] Comment on post
- [ ] Like/dislike post and comments
- [ ] Send message to online user
- [ ] Receive message (open two browsers)
- [ ] See typing indicator
- [ ] See online/offline status
- [ ] Load more messages (scroll up)
- [ ] Navigate between all pages
- [ ] Browser back/forward buttons work
- [ ] Logout
- [ ] Refresh page (should maintain state)

### 8.2 Browser Testing
- [ ] Chrome
- [ ] Firefox
- [ ] Safari
- [ ] Mobile responsive

### 8.3 Fix Bugs
- [ ] Document and fix all bugs found

---

## Completion Checklist

### Exercise Requirements
- [ ] ✅ Single HTML file (`index.html`)
- [ ] ✅ JavaScript handles all frontend logic
- [ ] ✅ Client-side routing (no page refreshes)
- [ ] ✅ Registration with all fields (nickname, age, gender, first name, last name, email, password)
- [ ] ✅ Login with username OR email + password
- [ ] ✅ Logout from any page
- [ ] ✅ Create posts with categories
- [ ] ✅ Create comments on posts
- [ ] ✅ View posts in feed
- [ ] ✅ View comments on post detail page
- [ ] ✅ Private messaging with real-time delivery
- [ ] ✅ User list showing online/offline status
- [ ] ✅ Conversations sorted by last message time
- [ ] ✅ New users sorted alphabetically
- [ ] ✅ Chat displays previous messages
- [ ] ✅ Load 10 messages at a time (scroll up for more)
- [ ] ✅ Messages show date and username
- [ ] ✅ WebSocket for real-time messages
- [ ] ✅ No frontend frameworks (vanilla JS)
- [ ] ✅ Only allowed packages used

### Technical Quality
- [ ] ✅ Code is clean and organized
- [ ] ✅ Good separation of concerns
- [ ] ✅ Error handling implemented
- [ ] ✅ No console errors
- [ ] ✅ Responsive design
- [ ] ✅ Accessible UI

---

## Time Estimate

| Phase | Estimated Time |
|-------|---------------|
| 1. Project Setup | 30 min |
| 2. Core Infrastructure | 2-3 hours |
| 3. WebSocket Client | 2-3 hours |
| 4. Authentication Views | 3-4 hours |
| 5. Chat Interface | 6-8 hours |
| 6. Posts & Comments | 4-6 hours |
| 7. Static Server | 30 min |
| 8. Testing & Polish | 2-3 hours |
| **TOTAL** | **20-28 hours** |

With focused work, this could be done in **3-5 days full-time** or **1-2 weeks part-time**.

---

## Next Steps

1. **Review this plan** - Does it make sense?
2. **Backup old frontend** - `mv frontend frontend-old`
3. **Start Phase 1** - Create directory structure
4. **Work incrementally** - Complete each phase before moving on

Ready to start? Let me know which phase to begin with!
