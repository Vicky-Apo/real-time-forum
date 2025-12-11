// components/Navbar.js - Navigation Bar Component

import state from '../state.js';
import apiClient from '../api/client.js';
import { navigate } from '../router.js';

export function renderNavbar() {
    const navbar = document.getElementById('navbar');
    const user = state.getUser();

    if (!user) {
        // Not logged in - show minimal navbar
        navbar.innerHTML = `
            <div class="navbar-container">
                <div class="navbar-brand">
                    <h1>Real-Time Forum</h1>
                </div>
            </div>
        `;
        return;
    }

    // Logged in - show full navbar
    const unreadCount = state.unreadCount;
    console.log('[Navbar] Unread count:', unreadCount);
    console.log('[Navbar] Unread count type:', typeof unreadCount);
    const badgeHTML = `<span class="badge">${unreadCount}</span>`;
    console.log('[Navbar] Badge HTML:', badgeHTML);
    const unreadBadge = unreadCount > 0 ? badgeHTML : '';

    navbar.innerHTML = `
        <div class="navbar-container">
            <div class="navbar-brand">
                <h1><a href="/" data-link>Real-Time Forum</a></h1>
            </div>

            <nav class="navbar-menu">
                <a href="/" data-link class="nav-link">
                    Home
                </a>
                <a href="/chat" data-link class="nav-link">
                    Chat
                </a>
                <a href="/notifications" data-link class="nav-link" style="position: relative;">
                    Notifications
                    ${unreadBadge}
                </a>
                <a href="/create-post" data-link class="nav-link btn-primary">
                    New Post
                </a>
            </nav>

            <div class="navbar-user">
                <div class="user-menu">
                    <button class="user-menu-btn" id="user-menu-btn">
                        <div class="user-avatar">
                            ${user.username.slice(0, 2).toUpperCase()}
                        </div>
                        <span class="user-name">${user.username}</span>
                        <span class="dropdown-icon">▼</span>
                    </button>
                    <div class="user-dropdown" id="user-dropdown">
                        <div class="user-dropdown-header">
                            <div class="user-avatar-large">
                                ${user.username.slice(0, 2).toUpperCase()}
                            </div>
                            <div class="user-info">
                                <div class="user-display-name">${user.username}</div>
                                <div class="user-email">${user.email || 'User Account'}</div>
                            </div>
                        </div>
                        <div class="user-dropdown-divider"></div>
                        <a href="/profile/${user.user_id}" data-link class="dropdown-item">
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
        userDropdown.classList.toggle('show');
    });

    // Close dropdown when clicking outside
    document.addEventListener('click', () => {
        userDropdown?.classList.remove('show');
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

// Update navbar when unread count changes
state.on('unread:changed', () => {
    renderNavbar();
});

// Update navbar when user changes
state.on('user:changed', () => {
    renderNavbar();
});
