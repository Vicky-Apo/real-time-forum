# Real-Time Forum: DevOps-Oriented Refactoring Plan

## Executive Summary

This plan transforms the real-time forum into a **production-ready, DevOps portfolio project** while ensuring full compliance with 01 Edu exercise requirements.

**Current Status:**
- ✅ Backend: Well-architected Go API with WebSocket support
- ❌ Frontend: Multi-page Go app (13 HTML files) - **Does not meet requirements**
- ⚠️ DevOps: Basic Docker setup, no CI/CD, testing, or monitoring

**Target Architecture:**
- Backend: Go API (keep existing, enhance with observability)
- Frontend: Single-page vanilla JavaScript application
- Infrastructure: Kubernetes-ready with full CI/CD pipeline
- Observability: Structured logging, metrics, tracing, health checks

---

## Phase 1: Frontend Reconstruction (Week 1-2)
### Meet Exercise Requirements

### 1.1 Single-Page Application Foundation
**Goal:** Convert from 13 HTML files → 1 HTML file with client-side routing

**Tasks:**
- [ ] Create new `client/` directory for SPA frontend
- [ ] Build single `index.html` with mounting point
- [ ] Implement vanilla JS router using History API
- [ ] Create view/component system (no frameworks allowed)
- [ ] Set up ES6 module structure

**Deliverables:**
```
client/
├── index.html              # Single HTML file
├── js/
│   ├── main.js            # Entry point
│   ├── router.js          # Client-side routing
│   ├── api/               # Backend API client
│   │   ├── auth.js
│   │   ├── posts.js
│   │   ├── comments.js
│   │   └── messages.js
│   ├── views/             # Page views
│   │   ├── LoginView.js
│   │   ├── RegisterView.js
│   │   ├── HomeView.js
│   │   ├── PostView.js
│   │   ├── CreatePostView.js
│   │   ├── ProfileView.js
│   │   └── ChatView.js    # NEW: Chat interface
│   ├── components/        # Reusable components
│   │   ├── Navbar.js
│   │   ├── PostCard.js
│   │   ├── CommentList.js
│   │   ├── MessageThread.js    # NEW
│   │   └── UserList.js         # NEW
│   ├── websocket/         # WebSocket client
│   │   ├── WebSocketManager.js
│   │   └── handlers.js
│   └── utils/
│       ├── state.js       # State management
│       └── dom.js         # DOM helpers
└── css/
    ├── global.css         # Reuse existing
    ├── components.css
    └── chat.css           # NEW: Chat styles
```

### 1.2 WebSocket Client Implementation
**Goal:** Real-time messaging and presence

**Tasks:**
- [ ] Create WebSocket connection manager
- [ ] Implement reconnection logic with exponential backoff
- [ ] Handle WebSocket events:
  - `user_online` / `user_offline` → Update user list
  - `typing_start` / `typing_stop` → Show typing indicators
  - `new_message` → Display message in real-time
  - `message_read` → Update read status
- [ ] Integrate with state management

**WebSocket Client Architecture:**
```javascript
class WebSocketManager {
  constructor(url) {
    this.url = url;
    this.ws = null;
    this.reconnectAttempts = 0;
    this.handlers = new Map();
  }

  connect() {
    this.ws = new WebSocket(this.url);
    this.ws.onopen = () => this.onOpen();
    this.ws.onmessage = (e) => this.onMessage(e);
    this.ws.onerror = (e) => this.onError(e);
    this.ws.onclose = () => this.onClose();
  }

  send(type, payload) {
    this.ws.send(JSON.stringify({ type, payload }));
  }

  on(eventType, handler) {
    this.handlers.set(eventType, handler);
  }
}
```

### 1.3 Chat Interface (Private Messaging)
**Goal:** Complete chat UI with all requirements

**Requirements:**
- ✅ User list showing online/offline status
- ✅ Conversations sorted by last message time
- ✅ New users sorted alphabetically
- ✅ Message thread with date + username
- ✅ Load last 10 messages, scroll up for more (throttled)
- ✅ Real-time message delivery
- ✅ Typing indicators
- ✅ Unread message badges

**UI Components:**
```
ChatView
├── UserList (left sidebar)
│   ├── Search/filter users
│   ├── Online users (green dot)
│   ├── Offline users (gray dot)
│   └── Sort: last message → alphabetical
├── MessageThread (main area)
│   ├── Header (recipient name, status)
│   ├── Messages (virtualized scrolling)
│   │   ├── Date dividers
│   │   ├── Message bubbles (sender/receiver)
│   │   └── Read receipts
│   ├── Typing indicator
│   └── Input area
│       ├── Textarea (auto-expand)
│       └── Send button
└── ConversationList (mobile: replaces UserList)
```

**Pagination Strategy:**
```javascript
// Load messages in chunks of 10
const loadMessages = async (userId, offset = 0) => {
  const messages = await api.messages.getConversation(userId, {
    limit: 10,
    offset: offset
  });
  return messages;
};

// Throttled scroll handler
const handleScroll = throttle((e) => {
  if (e.target.scrollTop < 100) {
    loadMoreMessages();
  }
}, 500);
```

### 1.4 Enhanced Registration Form
**Goal:** Capture all required user information

**Required Fields:**
- Nickname (existing)
- Email (existing)
- Password (existing)
- **NEW:** Age
- **NEW:** Gender
- **NEW:** First Name
- **NEW:** Last Name

**Backend Changes:**
- [ ] Update `users` table schema (add columns)
- [ ] Update `UserRegistration` model in [server/internal/models/users_models.go](server/internal/models/users_models.go)
- [ ] Update validation in [server/internal/utils/validate.go](server/internal/utils/validate.go)
- [ ] Update registration handler in [server/internal/handlers/user_handler.go](server/internal/handlers/user_handler.go)

**Frontend Changes:**
- [ ] Add fields to registration view
- [ ] Client-side validation
- [ ] Update API client

### 1.5 Static File Server
**Goal:** Replace Go frontend server with simple static server

**Options:**
1. **Nginx** (production)
2. **Go http.FileServer** (simple)
3. **Caddy** (modern, auto-HTTPS)

**Recommended: Nginx**
```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    # SPA routing: always serve index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy to backend
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # WebSocket proxy
    location /ws {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }
}
```

---

## Phase 2: Backend Enhancement (Week 2-3)
### Production Readiness

### 2.1 Structured Logging
**Goal:** Replace fmt.Printf with production-grade logging

**Implementation:**
- [ ] Add `log/slog` (Go 1.21+ standard library)
- [ ] JSON logging for machine parsing
- [ ] Log levels (DEBUG, INFO, WARN, ERROR)
- [ ] Request ID tracing
- [ ] Contextual logging

**Example:**
```go
// Before
fmt.Printf("User registered: %s\n", username)

// After
slog.Info("user registered",
    "username", username,
    "user_id", userID,
    "request_id", requestID,
    "duration_ms", duration.Milliseconds(),
)
```

**Files to Update:**
- All 12 handler files
- Middleware files
- WebSocket hub
- Database layer

### 2.2 Health & Readiness Endpoints
**Goal:** Kubernetes-compatible health checks

**Endpoints:**
- `GET /health` - Liveness probe (server responding?)
- `GET /ready` - Readiness probe (dependencies healthy?)

**Implementation:**
```go
// internal/handlers/health_handler.go
type HealthHandler struct {
    db *sql.DB
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
    // Simple check: is server alive?
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
    // Check dependencies
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()

    if err := h.db.PingContext(ctx); err != nil {
        http.Error(w, "database unhealthy", http.StatusServiceUnavailable)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "status": "ready",
        "database": "healthy",
    })
}
```

### 2.3 Prometheus Metrics
**Goal:** Observability for monitoring and alerting

**Metrics to Track:**
- HTTP request duration (histogram)
- HTTP request count by status code (counter)
- Active WebSocket connections (gauge)
- Database query duration (histogram)
- Active sessions (gauge)
- Message send rate (counter)

**Implementation:**
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    wsConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "websocket_connections_active",
            Help: "Number of active WebSocket connections",
        },
    )
)

// Expose at /metrics
http.Handle("/metrics", promhttp.Handler())
```

### 2.4 Database Migrations
**Goal:** Version-controlled schema changes

**Tool:** [golang-migrate/migrate](https://github.com/golang-migrate/migrate)

**Migration Files:**
```
migrations/
├── 000001_initial_schema.up.sql
├── 000001_initial_schema.down.sql
├── 000002_add_user_profile_fields.up.sql
├── 000002_add_user_profile_fields.down.sql
└── ...
```

**Migration Script:**
```bash
#!/bin/bash
# scripts/migrate.sh
migrate -path ./migrations -database "sqlite3://./forum.db" up
```

**Integration:**
- [ ] Extract current schema to migration 000001
- [ ] Create migration for new user fields
- [ ] Update Dockerfile to run migrations on startup
- [ ] Document rollback procedures

### 2.5 Graceful Shutdown
**Goal:** Handle signals properly for zero-downtime deploys

**Implementation:**
```go
// cmd/main.go
func main() {
    server := &http.Server{
        Addr:    fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort),
        Handler: handler,
    }

    // Start server in goroutine
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    slog.Info("shutting down server...")

    // Graceful shutdown with 30s timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Close WebSocket connections
    hub.Shutdown()

    if err := server.Shutdown(ctx); err != nil {
        slog.Error("server forced to shutdown", "error", err)
    }

    slog.Info("server exited")
}
```

---

## Phase 3: Testing Infrastructure (Week 3-4)
### Critical Gap

### 3.1 Unit Tests
**Goal:** Test business logic in isolation

**Coverage Targets:**
- Handlers: 70%+
- Repositories: 80%+
- Utilities: 90%+
- Middleware: 80%+

**Example Test Structure:**
```go
// internal/repository/user_repository_test.go
func TestUserRepository_CreateUser(t *testing.T) {
    // Setup: in-memory test database
    db := setupTestDB(t)
    defer db.Close()

    repo := repository.NewUserRepository(db)

    // Test
    user := &models.UserRegistration{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
    }

    userID, err := repo.CreateUser(context.Background(), user)

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, userID)

    // Verify in database
    var count int
    db.QueryRow("SELECT COUNT(*) FROM users WHERE user_id = ?", userID).Scan(&count)
    assert.Equal(t, 1, count)
}
```

**Files to Create:**
- `*_test.go` for every Go file
- `testutils/` package for test helpers
- `fixtures/` for test data

### 3.2 Integration Tests
**Goal:** Test API endpoints end-to-end

**Framework:** Native `net/http/httptest`

**Example:**
```go
// internal/handlers/integration_test.go
func TestRegisterLoginFlow(t *testing.T) {
    // Setup test server
    ts := httptest.NewServer(setupRoutes())
    defer ts.Close()

    // Test registration
    resp := testRequest(t, ts, "POST", "/api/auth/register", map[string]string{
        "username": "newuser",
        "email":    "new@example.com",
        "password": "SecurePass123!",
    })
    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    // Test login
    resp = testRequest(t, ts, "POST", "/api/auth/login", map[string]string{
        "identifier": "newuser",
        "password":   "SecurePass123!",
    })
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Extract session cookie
    cookies := resp.Cookies()
    assert.NotEmpty(t, cookies)
}
```

### 3.3 E2E Tests (Frontend)
**Goal:** Test user flows in browser

**Tool:** Playwright (JavaScript)

**Test Scenarios:**
- [ ] Registration → Login → Create Post → Logout
- [ ] Login → Send Message → Receive Reply (WebSocket)
- [ ] Login → React to Post → Verify Notification
- [ ] Navigation: all routes accessible

**Example:**
```javascript
// tests/e2e/chat.spec.js
const { test, expect } = require('@playwright/test');

test('users can send messages in real-time', async ({ page, context }) => {
  // User 1: Login
  await page.goto('http://localhost:3000');
  await page.fill('input[name="username"]', 'user1');
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');

  // Navigate to chat
  await page.click('a[href="/chat"]');

  // Open new page for User 2
  const page2 = await context.newPage();
  await page2.goto('http://localhost:3000');
  // ... User 2 login

  // User 1 sends message
  await page.fill('textarea[name="message"]', 'Hello from User 1');
  await page.click('button#send');

  // User 2 receives message (real-time via WebSocket)
  await expect(page2.locator('text=Hello from User 1')).toBeVisible({ timeout: 2000 });
});
```

### 3.4 Load Testing
**Goal:** Identify performance bottlenecks

**Tool:** k6 (Go-based, easy to use)

**Scenarios:**
- 100 concurrent users browsing posts
- 50 users sending messages simultaneously
- WebSocket connection stress test (500+ connections)

**Example:**
```javascript
// tests/load/post_browsing.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 }, // Ramp up to 100 users
    { duration: '5m', target: 100 }, // Stay at 100 users
    { duration: '2m', target: 0 },   // Ramp down
  ],
};

export default function () {
  let res = http.get('http://localhost:8080/api/posts');
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
  sleep(1);
}
```

---

## Phase 4: CI/CD Pipeline (Week 4)
### DevOps Core

### 4.1 GitHub Actions Workflow
**Goal:** Automated testing, building, and deployment

**Workflows:**

#### `.github/workflows/ci.yml` - Continuous Integration
```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Cache Go modules
        uses: actions/cache@v4
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

      - name: Run tests
        run: |
          cd server
          go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./server/coverage.out

  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install Playwright
        run: |
          cd client
          npm ci
          npx playwright install --with-deps

      - name: Run E2E tests
        run: |
          cd client
          npm run test:e2e

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          working-directory: server

  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'

      - name: Upload to GitHub Security
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'

  build-docker:
    runs-on: ubuntu-latest
    needs: [test-backend, test-frontend, lint]
    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build backend image
        uses: docker/build-push-action@v5
        with:
          context: ./server
          push: false
          tags: forum-backend:${{ github.sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Build frontend image
        uses: docker/build-push-action@v5
        with:
          context: ./client
          push: false
          tags: forum-frontend:${{ github.sha }}
```

#### `.github/workflows/cd.yml` - Continuous Deployment
```yaml
name: CD

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v2

      - name: Build, tag, and push images
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_REGISTRY/forum-backend:$IMAGE_TAG ./server
          docker build -t $ECR_REGISTRY/forum-frontend:$IMAGE_TAG ./client
          docker push $ECR_REGISTRY/forum-backend:$IMAGE_TAG
          docker push $ECR_REGISTRY/forum-frontend:$IMAGE_TAG

      - name: Deploy to Kubernetes
        run: |
          aws eks update-kubeconfig --name forum-cluster --region us-east-1
          kubectl set image deployment/backend backend=$ECR_REGISTRY/forum-backend:$IMAGE_TAG
          kubectl set image deployment/frontend frontend=$ECR_REGISTRY/forum-frontend:$IMAGE_TAG
          kubectl rollout status deployment/backend
          kubectl rollout status deployment/frontend
```

### 4.2 Code Quality Tools

#### golangci-lint Configuration
```yaml
# .golangci.yml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - revive
    - gosec  # Security checks

linters-settings:
  gosec:
    severity: medium
  revive:
    rules:
      - name: exported
        severity: warning
```

### 4.3 Pre-commit Hooks
**Goal:** Catch issues before they reach CI

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json

  - repo: https://github.com/golangci/golangci-lint
    rev: v1.55.2
    hooks:
      - id: golangci-lint
        args: [--fix]

  - repo: local
    hooks:
      - id: go-test
        name: go test
        entry: bash -c 'cd server && go test ./...'
        language: system
        pass_filenames: false
```

---

## Phase 5: Kubernetes Deployment (Week 5)
### Container Orchestration

### 5.1 Kubernetes Manifests

#### Backend Deployment
```yaml
# k8s/backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: forum-backend
  namespace: forum
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: forum-backend
  template:
    metadata:
      labels:
        app: forum-backend
        version: v1
    spec:
      containers:
      - name: backend
        image: your-registry/forum-backend:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: SERVER_HOST
          value: "0.0.0.0"
        - name: SERVER_PORT
          value: "8080"
        - name: DB_PATH
          value: "/data/forum.db"
        - name: ENVIRONMENT
          value: "production"
        envFrom:
        - secretRef:
            name: forum-secrets
        - configMapRef:
            name: forum-config
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        volumeMounts:
        - name: data
          mountPath: /data
        - name: uploads
          mountPath: /app/uploads
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: forum-db-pvc
      - name: uploads
        persistentVolumeClaim:
          claimName: forum-uploads-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: forum-backend
  namespace: forum
spec:
  selector:
    app: forum-backend
  ports:
  - protocol: TCP
    port: 8080
    targetPort: 8080
    name: http
  type: ClusterIP
```

#### Frontend Deployment (Nginx)
```yaml
# k8s/frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: forum-frontend
  namespace: forum
spec:
  replicas: 2
  selector:
    matchLabels:
      app: forum-frontend
  template:
    metadata:
      labels:
        app: forum-frontend
    spec:
      containers:
      - name: nginx
        image: your-registry/forum-frontend:latest
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 10
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: forum-frontend
  namespace: forum
spec:
  selector:
    app: forum-frontend
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
  type: ClusterIP
```

#### Ingress
```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: forum-ingress
  namespace: forum
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/websocket-services: "forum-backend"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - forum.example.com
    secretName: forum-tls
  rules:
  - host: forum.example.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: forum-backend
            port:
              number: 8080
      - path: /ws
        pathType: Prefix
        backend:
          service:
            name: forum-backend
            port:
              number: 8080
      - path: /
        pathType: Prefix
        backend:
          service:
            name: forum-frontend
            port:
              number: 80
```

#### Persistent Storage
```yaml
# k8s/storage.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: forum-db-pvc
  namespace: forum
spec:
  accessModes:
  - ReadWriteOnce
  storageClassName: gp3
  resources:
    requests:
      storage: 10Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: forum-uploads-pvc
  namespace: forum
spec:
  accessModes:
  - ReadWriteMany  # Multiple pods need access
  storageClassName: efs
  resources:
    requests:
      storage: 50Gi
```

#### ConfigMap & Secrets
```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: forum-config
  namespace: forum
data:
  FRONTEND_BASE_URL: "https://forum.example.com"
  BACKEND_BASE_URL: "https://forum.example.com/api"
  ALLOWED_ORIGINS: "https://forum.example.com"
  DB_MAX_CONNECTIONS: "25"
  RATE_LIMIT_REQUESTS: "1000"
  RATE_LIMIT_WINDOW: "60"
---
apiVersion: v1
kind: Secret
metadata:
  name: forum-secrets
  namespace: forum
type: Opaque
stringData:
  GITHUB_CLIENT_ID: "your-github-client-id"
  GITHUB_CLIENT_SECRET: "your-github-client-secret"
  GOOGLE_CLIENT_ID: "your-google-client-id"
  GOOGLE_CLIENT_SECRET: "your-google-client-secret"
```

### 5.2 Helm Chart (Optional but Recommended)
**Goal:** Templated, reusable Kubernetes configs

```
helm/
├── Chart.yaml
├── values.yaml
├── values-prod.yaml
├── values-staging.yaml
└── templates/
    ├── backend-deployment.yaml
    ├── frontend-deployment.yaml
    ├── ingress.yaml
    ├── configmap.yaml
    └── secrets.yaml
```

**Deployment:**
```bash
helm install forum ./helm -f values-prod.yaml
helm upgrade forum ./helm -f values-prod.yaml
```

---

## Phase 6: Monitoring & Observability (Week 5-6)
### Production Operations

### 6.1 Prometheus + Grafana Stack

#### Deploy Monitoring Stack
```bash
# Using kube-prometheus-stack Helm chart
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install monitoring prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace
```

#### ServiceMonitor for Backend
```yaml
# k8s/servicemonitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: forum-backend
  namespace: forum
spec:
  selector:
    matchLabels:
      app: forum-backend
  endpoints:
  - port: http
    path: /metrics
    interval: 15s
```

#### Grafana Dashboard
**Panels:**
- HTTP request rate (requests/sec)
- HTTP error rate (4xx, 5xx)
- Request duration (p50, p95, p99)
- Active WebSocket connections
- Database query duration
- Active user sessions
- Message send rate

**Import from:** [Grafana Dashboard JSON](https://grafana.com/grafana/dashboards/)

### 6.2 Alerting Rules

```yaml
# k8s/prometheusrule.yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: forum-alerts
  namespace: forum
spec:
  groups:
  - name: forum
    interval: 30s
    rules:
    - alert: HighErrorRate
      expr: |
        sum(rate(http_requests_total{status=~"5.."}[5m]))
        /
        sum(rate(http_requests_total[5m])) > 0.05
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "High error rate detected"
        description: "Error rate is {{ $value | humanizePercentage }}"

    - alert: HighRequestLatency
      expr: |
        histogram_quantile(0.95,
          sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
        ) > 1
      for: 10m
      labels:
        severity: warning
      annotations:
        summary: "High request latency"
        description: "95th percentile latency is {{ $value }}s"

    - alert: DatabaseDown
      expr: up{job="forum-backend"} == 0
      for: 1m
      labels:
        severity: critical
      annotations:
        summary: "Backend is down"

    - alert: LowWebSocketConnections
      expr: websocket_connections_active < 1
      for: 30m
      labels:
        severity: info
      annotations:
        summary: "No active WebSocket connections"
```

### 6.3 Centralized Logging

#### Deploy EFK Stack (Elasticsearch, Fluentd, Kibana)
**Alternative:** Loki (lightweight, cost-effective)

```yaml
# k8s/fluentd-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
  namespace: logging
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/forum-*.log
      pos_file /var/log/fluentd-containers.log.pos
      tag kubernetes.*
      <parse>
        @type json
        time_key timestamp
        time_format %Y-%m-%dT%H:%M:%S.%NZ
      </parse>
    </source>

    <filter kubernetes.**>
      @type kubernetes_metadata
    </filter>

    <match kubernetes.**>
      @type elasticsearch
      host elasticsearch.logging.svc.cluster.local
      port 9200
      index_name forum-logs
      <buffer>
        @type memory
        flush_interval 5s
      </buffer>
    </match>
```

### 6.4 Distributed Tracing (Optional)

**Tool:** Jaeger or Tempo

**Implementation:**
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
    ctx, span := otel.Tracer("forum").Start(r.Context(), "GetPosts")
    defer span.End()

    posts, err := h.postRepo.GetAllPosts(ctx)
    if err != nil {
        span.RecordError(err)
        // ...
    }
    // ...
}
```

---

## Phase 7: Security Hardening (Ongoing)

### 7.1 Secrets Management
**Tool:** External Secrets Operator + AWS Secrets Manager

```yaml
# k8s/externalsecret.yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: forum-secrets
  namespace: forum
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: SecretStore
  target:
    name: forum-secrets
  data:
  - secretKey: GITHUB_CLIENT_SECRET
    remoteRef:
      key: forum/github-oauth
      property: client_secret
```

### 7.2 Network Policies
**Goal:** Restrict pod-to-pod communication

```yaml
# k8s/networkpolicy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: backend-netpol
  namespace: forum
spec:
  podSelector:
    matchLabels:
      app: forum-backend
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: forum-frontend
    - podSelector:
        matchLabels:
          app: nginx-ingress
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector: {}
    ports:
    - protocol: TCP
      port: 53  # DNS
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
```

### 7.3 Pod Security Standards
```yaml
# k8s/podsecuritypolicy.yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: restricted
spec:
  privileged: false
  allowPrivilegeEscalation: false
  requiredDropCapabilities:
  - ALL
  runAsUser:
    rule: MustRunAsNonRoot
  seLinux:
    rule: RunAsAny
  fsGroup:
    rule: RunAsAny
  volumes:
  - 'configMap'
  - 'emptyDir'
  - 'projected'
  - 'secret'
  - 'persistentVolumeClaim'
```

### 7.4 Image Scanning
**CI Integration:**
```yaml
# In GitHub Actions
- name: Scan image with Trivy
  uses: aquasecurity/trivy-action@master
  with:
    image-ref: 'forum-backend:${{ github.sha }}'
    format: 'table'
    exit-code: '1'
    severity: 'CRITICAL,HIGH'
```

---

## Phase 8: Documentation (Week 6)

### 8.1 API Documentation
**Tool:** Generate OpenAPI spec from code

```yaml
# openapi.yaml (excerpt)
openapi: 3.0.0
info:
  title: Real-Time Forum API
  version: 1.0.0
paths:
  /api/auth/register:
    post:
      summary: Register new user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UserRegistration'
      responses:
        '201':
          description: User created successfully
        '400':
          description: Invalid input
```

**Generate docs:** Use Swagger UI or Redoc

### 8.2 Architecture Diagrams
**Tools:** draw.io, Mermaid, PlantUML

**Diagrams to Create:**
1. System architecture (frontend, backend, database)
2. Network topology (Kubernetes services, ingress)
3. Data flow (user request → response)
4. WebSocket message flow
5. CI/CD pipeline flow
6. Deployment architecture

### 8.3 Runbooks
**Topics:**
- Local development setup
- Deploying to Kubernetes
- Rolling back a deployment
- Scaling the application
- Database backup and restore
- Incident response procedures
- Monitoring and alerting guide

### 8.4 README Updates
**Sections:**
- Project overview
- Architecture
- Technology stack
- Prerequisites
- Quick start (Docker Compose)
- Development workflow
- Testing strategy
- Deployment (Kubernetes)
- Monitoring and observability
- Contributing guidelines
- License

---

## Timeline Summary

| Phase | Duration | Key Deliverables |
|-------|----------|------------------|
| Phase 1: Frontend Reconstruction | 2 weeks | SPA, WebSocket client, Chat UI |
| Phase 2: Backend Enhancement | 1 week | Logging, metrics, health checks |
| Phase 3: Testing Infrastructure | 1 week | Unit, integration, E2E, load tests |
| Phase 4: CI/CD Pipeline | 1 week | GitHub Actions, automated testing |
| Phase 5: Kubernetes Deployment | 1 week | K8s manifests, Helm charts |
| Phase 6: Monitoring & Observability | 1 week | Prometheus, Grafana, alerts |
| Phase 7: Security Hardening | Ongoing | Secrets, network policies, scanning |
| Phase 8: Documentation | 1 week | API docs, runbooks, diagrams |

**Total Estimated Time:** 6-8 weeks full-time

---

## DevOps Skills Demonstrated

This project showcases entry-level DevOps/SysAdmin skills:

### Core DevOps Practices
- [x] Version control (Git)
- [x] CI/CD pipelines (GitHub Actions)
- [x] Infrastructure as Code (Kubernetes manifests, Helm)
- [x] Automated testing (unit, integration, E2E)
- [x] Containerization (Docker, multi-stage builds)
- [x] Container orchestration (Kubernetes)

### Observability
- [x] Structured logging (slog/JSON)
- [x] Metrics collection (Prometheus)
- [x] Visualization (Grafana dashboards)
- [x] Alerting (Prometheus Alertmanager)
- [x] Centralized logging (EFK/Loki)
- [x] Distributed tracing (optional: Jaeger)

### Security
- [x] Secrets management (External Secrets Operator)
- [x] Vulnerability scanning (Trivy)
- [x] Network policies
- [x] Pod security policies
- [x] Non-root containers
- [x] HTTPS/TLS (cert-manager)

### Reliability
- [x] Health checks (liveness, readiness)
- [x] Graceful shutdown
- [x] Rolling updates
- [x] Auto-scaling (HPA)
- [x] Database backups
- [x] Disaster recovery procedures

### Toolchain
- Go, JavaScript, SQL, Bash
- Docker, Kubernetes, Helm
- GitHub Actions, Prometheus, Grafana
- Trivy, golangci-lint, Playwright
- Nginx, SQLite, WebSockets

---

## Next Steps

1. **Review this plan** - Discuss priorities and timeline
2. **Set up project structure** - Create directories for `client/`, `k8s/`, etc.
3. **Start Phase 1** - Frontend reconstruction (highest priority for compliance)
4. **Iterate incrementally** - Complete each phase before moving to next

**Question for you:**
- Do you want to start with Phase 1 (frontend SPA) to meet exercise requirements first?
- Or should we tackle some quick DevOps wins (CI/CD, testing) in parallel?
- What's your timeline for completing this project?
