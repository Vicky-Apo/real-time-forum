// main.js - Application Entry Point

import Router from './router.js';
import state from './state.js';
import wsManager from './websocket/WebSocketManager.js';
import { renderNavbar } from './components/Navbar.js';
import apiClient from './api/client.js';

// Import views (lazy loaded by router)
import LoginView from './views/LoginView.js';
import RegisterView from './views/RegisterView.js';
import HomeView from './views/HomeView.js';
import PostView from './views/PostView.js';
import CreatePostView from './views/CreatePostView.js';
import CategoryView from './views/CategoryView.js';
import ProfileView from './views/ProfileView.js';
import ChatView from './views/ChatView.js';
import NotificationsView from './views/NotificationsView.js';

// Define routes
const routes = [
    {
        path: '/',
        component: async () => HomeView,
        requiresAuth: true
    },
    {
        path: '/login',
        component: async () => LoginView,
        requiresAuth: false
    },
    {
        path: '/register',
        component: async () => RegisterView,
        requiresAuth: false
    },
    {
        path: '/post/:id',
        component: async () => PostView,
        requiresAuth: true
    },
    {
        path: '/create-post',
        component: async () => CreatePostView,
        requiresAuth: true
    },
    {
        path: '/category/:id',
        component: async () => CategoryView,
        requiresAuth: true
    },
    {
        path: '/profile/:id',
        component: async () => ProfileView,
        requiresAuth: true
    },
    {
        path: '/chat',
        component: async () => ChatView,
        requiresAuth: true
    },
    {
        path: '/notifications',
        component: async () => NotificationsView,
        requiresAuth: true
    }
];

// Initialize application
async function initApp() {
    console.log('[App] Initializing application...');

    try {
        // Check if user is already logged in
        const currentUser = await checkAuthStatus();

        if (currentUser) {
            state.setUser(currentUser);
        }

        // Connect WebSocket if we have an authenticated user (from API or localStorage)
        const user = state.getUser();
        if (user) {
            console.log('[App] User authenticated, connecting WebSocket...');
            wsManager.connect();
        }

        // Initialize router
        const router = new Router(routes);
        window.router = router; // Make router globally accessible

        // Render navbar
        renderNavbar();

        // Handle initial route
        await router.handleRoute();

        console.log('[App] Application initialized successfully');
    } catch (error) {
        console.error('[App] Initialization failed:', error);
        showError('Failed to initialize application');
    }
}

// Check authentication status
async function checkAuthStatus() {
    try {
        const response = await apiClient.post('/auth/me');
        return response.user;
    } catch (error) {
        // Not authenticated
        return null;
    }
}

// Global error handler
function showError(message) {
    const app = document.getElementById('app');
    app.innerHTML = `
        <div class="error-container">
            <h2>Oops! Something went wrong</h2>
            <p>${message}</p>
            <button onclick="window.location.reload()" class="btn btn-primary">
                Reload Page
            </button>
        </div>
    `;
}

// Listen for state changes
state.on('user:changed', (user) => {
    console.log('[App] User state changed:', user ? user.username : 'logged out');
    renderNavbar();

    if (user && !state.wsConnected) {
        // User logged in, connect WebSocket
        wsManager.connect();
    } else if (!user && state.wsConnected) {
        // User logged out, disconnect WebSocket
        wsManager.disconnect();
    }
});

// Handle browser notification permissions
if ('Notification' in window && Notification.permission === 'default') {
    Notification.requestPermission().then(permission => {
        console.log('[App] Notification permission:', permission);
    });
}

// Start the application when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initApp);
} else {
    initApp();
}

// Export for debugging
window.state = state;
window.wsManager = wsManager;
