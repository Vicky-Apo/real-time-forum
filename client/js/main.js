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
import EditPostView from './views/EditPostView.js';
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
        path: '/edit-post/:id',
        component: async () => EditPostView,
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
        let currentUser = null;
        try {
            currentUser = await checkAuthStatus();
        } catch (error) {
            console.error('[App] Error checking auth status:', error);
            // Continue without user - they can still access login/register pages
        }

        if (currentUser) {
            state.setUser(currentUser);

            // Load initial unread notification count (non-blocking)
            setTimeout(async () => {
                try {
                    const notifResponse = await apiClient.get('/notifications');
                    const data = notifResponse.data || notifResponse;
                    const notifications = data.notifications || [];
                    const unreadCount = notifications.filter(n => !n.is_read).length;
                    state.setUnreadCount(unreadCount);
                } catch (error) {
                    console.error('[App] Error loading notifications count:', error);
                }
            }, 0);

            // Load initial unread message count (non-blocking)
            setTimeout(async () => {
                try {
                    const msgResponse = await apiClient.get('/messages/unread-count');
                    const msgData = msgResponse.data || msgResponse;
                    const unreadMessageCount = msgData.unread_count || 0;
                    state.setUnreadMessageCount(unreadMessageCount);
                } catch (error) {
                    console.error('[App] Error loading message count:', error);
                }
            }, 0);
        }

        // Initialize router
        const router = new Router(routes);
        window.router = router; // Make router globally accessible

        // Render navbar AFTER setting unread count
        try {
            renderNavbar();
        } catch (error) {
            console.error('[App] Error rendering navbar:', error);
            // Continue even if navbar fails
        }

        // Connect WebSocket if we have an authenticated user (from API or localStorage)
        // Use setTimeout to ensure this doesn't block page rendering
        const user = state.getUser();
        if (user) {
            console.log('[App] User authenticated, connecting WebSocket...');
            // Connect WebSocket asynchronously to not block page load
            setTimeout(() => {
                try {
                    wsManager.connect();
                } catch (error) {
                    console.error('[App] WebSocket connection failed, but continuing:', error);
                    // Don't let WebSocket errors prevent the app from loading
                }
            }, 100);
        }

        // Handle initial route
        try {
            await router.handleRoute();
        } catch (error) {
            console.error('[App] Error handling initial route:', error);
            // Show error but don't crash
            const app = document.getElementById('app');
            if (app) {
                app.innerHTML = `
                    <div class="error-container">
                        <h2>Error Loading Page</h2>
                        <p>${error.message || 'Unknown error occurred'}</p>
                        <button onclick="window.location.reload()" class="btn btn-primary">
                            Reload Page
                        </button>
                    </div>
                `;
            }
        }

        console.log('[App] Application initialized successfully');
    } catch (error) {
        console.error('[App] Initialization failed:', error);
        // Try to show error, but don't fail completely
        try {
            showError('Failed to initialize application: ' + (error.message || 'Unknown error'));
        } catch (e) {
            console.error('[App] Could not show error message:', e);
        }
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
        // User logged in, connect WebSocket (non-blocking)
        setTimeout(() => {
            try {
                wsManager.connect();
            } catch (error) {
                console.error('[App] WebSocket connection failed:', error);
            }
        }, 100);
    } else if (!user && state.wsConnected) {
        // User logged out, disconnect WebSocket
        try {
            wsManager.disconnect();
        } catch (error) {
            console.error('[App] WebSocket disconnect error:', error);
        }
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
