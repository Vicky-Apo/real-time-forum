// components/Navbar.js - Navigation Bar Component

import state from '../state.js';
import apiClient from '../api/client.js';
import { navigate } from '../router.js';
import { getInitials } from '../utils/helpers.js';

export function renderNavbar() {
    const navbar = document.getElementById('navbar');
    const user = state.getUser();

    if (!user) {
        // Not logged in - show minimal navbar
        navbar.innerHTML = `
            <div class="navbar-container">
                <div class="navbar-brand">
                    <h1>
                        <i class="fas fa-robot robot-icon"></i>
                        Real-Time Forum
                    </h1>
                </div>
            </div>
        `;
        return;
    }

    // Logged in - show full navbar
    // Notification badge
    const unreadCount = state.unreadCount;
    const notificationBadge = unreadCount > 0 ? `<span class="badge">${unreadCount}</span>` : '';

    // Message badge
    const unreadMessageCount = state.unreadMessageCount;
    const messageBadge = unreadMessageCount > 0 ? `<span class="badge">${unreadMessageCount}</span>` : '';

    // Use id or user_id depending on what's available
    const userId = user.id || user.user_id;

    navbar.innerHTML = `
        <div class="navbar-container">
            <div class="navbar-brand">
                <h1>
                    <i class="fas fa-robot robot-icon"></i>
                    <a href="/" data-link>Real-Time Forum</a>
                </h1>
            </div>

            <nav class="navbar-menu">
                <a href="/" data-link class="nav-link">
                    Home
                </a>
                <a href="/chat" data-link class="nav-link" style="position: relative;">
                    Chat
                    ${messageBadge}
                </a>
                <a href="/notifications" data-link class="nav-link" style="position: relative;">
                    Notifications
                    ${notificationBadge}
                </a>
                <a href="/create-post" data-link class="nav-link btn-primary">
                    New Post
                </a>
            </nav>

            <div class="navbar-user">
                <div class="user-menu">
                   <button class="user-menu-btn" 
                        id="user-menu-btn" 
                        aria-label="User menu"
                        aria-expanded="false"
                        aria-haspopup="true">
                        <div class="user-avatar">
                            ${getInitials(user.username)}
                        </div>
                        <span class="user-name">${user.username}</span>
                        <span class="dropdown-icon">▼</span>
                    </button>
                    <div class="user-dropdown" id="user-dropdown">
                        <div class="user-dropdown-header">
                            <div class="user-avatar-large">
                                ${getInitials(user.username)}
                            </div>
                            <div class="user-info">
                                <div class="user-display-name">${user.username}</div>
                                <div class="user-email">${user.email || 'User Account'}</div>
                            </div>
                        </div>
                        <div class="user-dropdown-divider"></div>
                        <a href="/profile/${userId}" data-link class="dropdown-item">
                            <span class="dropdown-icon-left">👤</span>
                            <span>My Profile</span>
                        </a>
                        <a href="#" id="logout-btn" class="dropdown-item logout-item">
                            <span class="dropdown-icon-left">🚪</span>
                            <span>Logout</span>
                        </a>
                    </div>
                </div>
            </div>
        </div>
    `;

    // Add event listeners
    const userMenuBtn = document.getElementById('user-menu-btn');
    const userDropdown = document.getElementById('user-dropdown');
    const logoutBtn = document.getElementById('logout-btn');

    // Toggle dropdown
    userMenuBtn?.addEventListener('click', (e) => {
        e.stopPropagation();
        const isExpanded = userDropdown.classList.toggle('show');
        userMenuBtn.setAttribute('aria-expanded', isExpanded);
    });

    // Add keyboard support
userMenuBtn?.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        userDropdown?.classList.remove('show');
        userMenuBtn.setAttribute('aria-expanded', 'false');
    }
    if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        userMenuBtn.click();
    }
});

// Close dropdown when clicking outside
document.addEventListener('click', () => {
    const wasShown = userDropdown?.classList.contains('show');
    userDropdown?.classList.remove('show');
    if (wasShown) {
        userMenuBtn.setAttribute('aria-expanded', 'false');
    }
});

    // Logout
    logoutBtn?.addEventListener('click', async (e) => {
        e.preventDefault();
        await handleLogout();
    });
}

async function handleLogout() {
    try {
        await apiClient.post('/auth/logout');
        state.clear();
        navigate('/login');
    } catch (error) {
        console.error('[Navbar] Logout error:', error);
        // Force logout anyway
        state.clear();
        navigate('/login');
    }
}

// Update navbar when unread notification count changes
state.on('unread:changed', () => {
    renderNavbar();
});

// Update navbar when unread message count changes
state.on('unread-messages:changed', () => {
    renderNavbar();
});

// Update navbar when user changes
state.on('user:changed', () => {
    renderNavbar();
});
