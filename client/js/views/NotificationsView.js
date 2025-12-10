// views/NotificationsView.js - Notifications Page

import apiClient from '../api/client.js';
import { navigate } from '../router.js';
import state from '../state.js';

export default {
    notifications: [],

    async render() {
        return `
            <div class="container">
                <div class="notifications-container">
                    <div class="notifications-header">
                        <h1>📬 Notifications</h1>
                        <p>Stay updated with your forum activity</p>
                    </div>

                    <div id="notifications-list" class="notifications-list">
                        <div class="loading-container">
                            <div class="loading-spinner"></div>
                            <p>Loading notifications...</p>
                        </div>
                    </div>
                </div>
            </div>
        `;
    },

    async afterRender() {
        console.log('[NotificationsView] Rendered');
        await this.loadNotifications();
        this.setupWebSocketListeners();
    },

    async loadNotifications() {
        const container = document.getElementById('notifications-list');

        try {
            const response = await apiClient.get('/notifications');
            const data = response.data || response;
            this.notifications = data.notifications || [];

            console.log('[NotificationsView] Loaded notifications:', this.notifications.length);
            console.log('[NotificationsView] Notifications data:', this.notifications);

            if (this.notifications.length === 0) {
                container.innerHTML = `
                    <div class="empty-state">
                        <div class="empty-state-icon">🔔</div>
                        <h3>No notifications yet</h3>
                        <p>When someone interacts with your posts, you'll see notifications here</p>
                    </div>
                `;
                return;
            }

            container.innerHTML = this.notifications.map(notif => this.renderNotification(notif)).join('');

            // Update unread count based on actual unread notifications
            const unreadCount = this.notifications.filter(n => !n.is_read).length;
            state.setUnreadCount(unreadCount);

        } catch (error) {
            console.error('[NotificationsView] Error loading notifications:', error);
            container.innerHTML = `
                <div class="error-message">
                    <p>Failed to load notifications. ${error.message}</p>
                    <button class="btn btn-secondary" onclick="window.location.reload()">Retry</button>
                </div>
            `;
        }
    },

    renderNotification(notif) {
        const isUnread = !notif.is_read;
        const unreadClass = isUnread ? 'notification-unread' : '';
        const timeAgo = this.getTimeAgo(notif.created_at);

        // Determine icon based on action
        let icon = '📝';
        let actionText = notif.action;

        if (notif.action.includes('liked')) {
            icon = '👍';
            actionText = 'liked your post';
        } else if (notif.action.includes('comment')) {
            icon = '💬';
            actionText = 'commented on your post';
        }

        return `
            <div class="notification-item ${unreadClass}" data-id="${notif.notification_id}" data-post-id="${notif.post_id}">
                <div class="notification-icon">${icon}</div>
                <div class="notification-content">
                    <div class="notification-text">
                        <strong>${notif.trigger_username}</strong> ${actionText}
                    </div>
                    <div class="notification-preview">"${notif.post_content_preview}"</div>
                    <div class="notification-time">${timeAgo}</div>
                </div>
                ${isUnread ? '<div class="notification-badge"></div>' : ''}
            </div>
        `;
    },

    getTimeAgo(timestamp) {
        const now = new Date();
        const past = new Date(timestamp);
        const diffMs = now - past;
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);

        if (diffMins < 1) return 'Just now';
        if (diffMins < 60) return `${diffMins}m ago`;
        if (diffHours < 24) return `${diffHours}h ago`;
        if (diffDays < 7) return `${diffDays}d ago`;

        return past.toLocaleDateString();
    },

    setupWebSocketListeners() {
        // Listen for new notifications
        state.on('notification:received', (payload) => {
            console.log('[NotificationsView] New notification received');
            this.loadNotifications(); // Reload to show new notification
        });

        // Add click handler to notification items
        document.addEventListener('click', async (e) => {
            const notificationItem = e.target.closest('.notification-item');
            if (notificationItem) {
                const postId = notificationItem.dataset.postId;
                const notifId = notificationItem.dataset.id;

                console.log('[NotificationsView] Clicked notification:', { notifId, postId });
                console.log('[NotificationsView] Dataset:', notificationItem.dataset);

                // Mark as read
                try {
                    await apiClient.post(`/notifications/mark-read/${notifId}`);
                    notificationItem.classList.remove('notification-unread');

                    // Decrement unread count
                    const currentUnread = state.unreadCount;
                    if (currentUnread > 0) {
                        state.setUnreadCount(currentUnread - 1);
                    }
                } catch (error) {
                    console.error('[NotificationsView] Error marking as read:', error);
                }

                // Navigate to post
                console.log('[NotificationsView] Navigating to /post/' + postId);
                navigate(`/post/${postId}`);
            }
        });
    },

    cleanup() {
        // Remove event listeners if needed
        state.off('notification:received');
    }
};
