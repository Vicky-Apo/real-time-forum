// Global instance
let notificationManager = null;

/**
 * Notification System - Frontend JavaScript
 * Handles notification polling, dropdown, and browser notifications
 */
class NotificationManager {
    constructor() {
        this.lastNotificationCount = 0;
        this.pollingInterval = null;
        this.pollingFrequency = 2000; // 2 seconds
        this.isPolling = false;
        this.isDropdownOpen = false;
        this.notifications = [];
        
        // DOM elements
        this.badge = null;
        this.bell = null;
        this.dropdown = null;
        this.notificationList = null;
        
        this.init();
    }

    init() {
        // Get DOM elements
        this.badge = document.getElementById('notification-badge');
        this.bell = document.getElementById('notification-bell');
        this.dropdown = document.getElementById('notification-dropdown');
        this.notificationList = document.getElementById('notification-list');
        
        if (!this.badge || !this.bell || !this.dropdown) {
            console.warn('Notification elements not found in DOM');
            return;
        }

        // Request browser notification permission
        this.requestNotificationPermission();
        
        // Start polling
        this.startPolling();
        
        // Add click outside to close dropdown
        document.addEventListener('click', (e) => {
            if (!e.target.closest('.notification-container')) {
                this.closeDropdown();
            }
        });
        
        // Add visibility change listener to pause/resume polling
        document.addEventListener('visibilitychange', () => {
            if (document.hidden) {
                this.pausePolling();
            } else {
                this.resumePolling();
            }
        });
        
        console.log('✅ NotificationManager initialized successfully');
    }

    requestNotificationPermission() {
        if ('Notification' in window && Notification.permission === 'default') {
            Notification.requestPermission().then(permission => {
                console.log('Notification permission:', permission);
            });
        }
    }

    startPolling() {
        if (this.isPolling) {
            return;
        }
        
        this.isPolling = true;
        
        // Initial call
        this.pollNotifications();
        
        // Set up interval
        this.pollingInterval = setInterval(() => {
            this.pollNotifications();
        }, this.pollingFrequency);
        
        console.log('🔔 Notification polling started');
    }

    stopPolling() {
        if (this.pollingInterval) {
            clearInterval(this.pollingInterval);
            this.pollingInterval = null;
        }
        this.isPolling = false;
        console.log('🔕 Notification polling stopped');
    }

    pausePolling() {
        this.stopPolling();
    }

    resumePolling() {
        if (!this.isPolling) {
            this.startPolling();
        }
    }

    async pollNotifications() {
        try {
            const response = await fetch('/api/notifications', {
                method: 'GET',
                credentials: 'include'
            });

            if (!response.ok) {
                if (response.status === 401) {
                    this.stopPolling();
                    return;
                }
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            const data = await response.json();
            this.notifications = data.notifications || [];
            this.updateNotificationBadge(data.unread_count || 0);
            
            // Update dropdown if it's open
            if (this.isDropdownOpen) {
                this.updateDropdownContent();
            }
            
        } catch (error) {
            console.error('Error polling notifications:', error);
        }
    }

    updateNotificationBadge(unreadCount) {
        if (!this.badge) {
            return;
        }

        if (unreadCount > 0) {
            this.badge.textContent = unreadCount > 99 ? '99+' : unreadCount;
            this.badge.style.display = 'flex';
            
            if (unreadCount > this.lastNotificationCount && this.lastNotificationCount > 0) {
                this.showBrowserNotification(
                    `You have ${unreadCount} unread notification${unreadCount > 1 ? 's' : ''}`
                );
            }
            
            this.updatePageTitle(unreadCount);
            
        } else {
            this.badge.style.display = 'none';
            this.updatePageTitle(0);
        }

        this.lastNotificationCount = unreadCount;
    }

    toggleDropdown() {
        if (this.isDropdownOpen) {
            this.closeDropdown();
        } else {
            this.openDropdown();
        }
    }

    openDropdown() {
        if (!this.dropdown) return;
        
        this.dropdown.classList.add('show');
        this.isDropdownOpen = true;
        
        // Load notifications immediately when opened
        this.updateDropdownContent();
        console.log('📂 Dropdown opened');
    }

    closeDropdown() {
        if (!this.dropdown) return;
        
        this.dropdown.classList.remove('show');
        this.isDropdownOpen = false;
        console.log('📁 Dropdown closed');
    }

    updateDropdownContent() {
        if (!this.notificationList) return;

        // Show loading initially
        const loadingElement = document.getElementById('notification-loading');
        const emptyElement = document.getElementById('notification-empty');
        
        if (loadingElement) loadingElement.style.display = 'block';
        if (emptyElement) emptyElement.style.display = 'none';

        // Clear existing notifications
        const existingNotifications = this.notificationList.querySelectorAll('.dropdown-notification');
        existingNotifications.forEach(notification => notification.remove());

        if (this.notifications.length === 0) {
            if (loadingElement) loadingElement.style.display = 'none';
            if (emptyElement) emptyElement.style.display = 'block';
            return;
        }

        // Hide loading
        if (loadingElement) loadingElement.style.display = 'none';

        // Show recent notifications (limit to 10 for dropdown)
        const recentNotifications = this.notifications.slice(0, 10);
        
        recentNotifications.forEach(notification => {
            const notificationElement = this.createNotificationElement(notification);
            this.notificationList.appendChild(notificationElement);
        });
    }

    createNotificationElement(notification) {
        const div = document.createElement('div');
        div.className = `dropdown-notification ${notification.is_read ? 'read' : 'unread'}`;
        div.onclick = () => {
            window.location.href = `/post/${notification.post_id}`;
        };

        const icon = this.getNotificationIcon(notification.action);
        const timeAgo = this.formatTimeAgo(new Date(notification.created_at));
        
        // ✅ FIXED: Just show the action directly from backend, no extra formatting
        const username = this.escapeHtml(notification.trigger_username);
        const preview = this.escapeHtml(notification.post_content_preview);
        const action = this.escapeHtml(notification.action);

        div.innerHTML = `
            <div class="notification-icon">${icon}</div>
            <div class="notification-content">
                <div class="notification-message">
                    <strong>${username}</strong> ${action} "<em>${preview}</em>"
                </div>
                <div class="notification-time">${timeAgo}</div>
            </div>
            ${!notification.is_read ? `
                <div class="notification-actions-btn">
                    <button class="mark-read-btn-small" onclick="event.stopPropagation(); markAsReadDropdown('${notification.notification_id}', this)">
                        ✓
                    </button>
                </div>
            ` : ''}
        `;

        return div;
    }

    getNotificationIcon(action) {
        if (action && action.includes('like') && !action.includes('dislike')) {
            return '👍';
        }
        if (action && action.includes('dislike')) {
            return '👎';
        }
        if (action && action.includes('comment')) {
            return '💬';
        }
        return '🔔';
    }

    formatTimeAgo(date) {
        const now = new Date();
        const diff = now - date;
        const seconds = Math.floor(diff / 1000);
        const minutes = Math.floor(seconds / 60);
        const hours = Math.floor(minutes / 60);
        const days = Math.floor(hours / 24);

        if (seconds < 60) return 'just now';
        if (minutes < 60) return `${minutes}m ago`;
        if (hours < 24) return `${hours}h ago`;
        if (days < 7) return `${days}d ago`;
        return date.toLocaleDateString();
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    showBrowserNotification(message) {
        if ('Notification' in window && Notification.permission === 'granted') {
            const notification = new Notification('Forum Notification', {
                body: message,
                icon: '/static/favicon.ico',
                tag: 'forum-notification',
                requireInteraction: false
            });

            setTimeout(() => notification.close(), 5000);

            notification.onclick = () => {
                window.focus();
                this.openDropdown();
                notification.close();
            };
        }
    }

    updatePageTitle(count) {
        const originalTitle = document.title.replace(/^\(\d+\)\s/, '');
        
        if (count > 0) {
            document.title = `(${count}) ${originalTitle}`;
        } else {
            document.title = originalTitle;
        }
    }

    async markAsRead(notificationId) {
        try {
            const response = await fetch(`/notifications/mark-read/${notificationId}`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            // Immediately poll to update badge and dropdown
            this.pollNotifications();
            
            return true;
        } catch (error) {
            console.error('Error marking notification as read:', error);
            return false;
        }
    }
}

// Initialize notifications when DOM is ready
function initNotifications() {
    // Only initialize if user is logged in (check for notification elements)
    if (document.getElementById('notification-badge') && document.getElementById('notification-bell')) {
        if (!notificationManager) {
            notificationManager = new NotificationManager();
        }
        return true;
    }
    return false;
}

// Global functions called from HTML templates
function toggleNotificationDropdown(event) {
    event.preventDefault();
    event.stopPropagation();
    
    // Initialize if not already done
    if (!notificationManager) {
        const initialized = initNotifications();
        if (!initialized) {
            console.error('Cannot initialize notifications - missing elements');
            return;
        }
    }
    
    if (notificationManager) {
        notificationManager.toggleDropdown();
    }
}

function closeNotificationDropdown() {
    if (notificationManager) {
        notificationManager.closeDropdown();
    }
}

async function markAsReadDropdown(notificationId, button) {
    if (!notificationManager) {
        console.error('Notification manager not initialized');
        return;
    }

    button.disabled = true;
    button.textContent = '...';

    try {
        const success = await notificationManager.markAsRead(notificationId);
        
        if (success) {
            const notification = button.closest('.dropdown-notification');
            if (notification) {
                notification.classList.remove('unread');
                notification.classList.add('read');
                
                const actionsDiv = button.closest('.notification-actions-btn');
                if (actionsDiv) {
                    actionsDiv.remove();
                }
            }
        } else {
            throw new Error('Failed to mark as read');
        }
    } catch (error) {
        console.error('Error marking notification as read:', error);
        button.disabled = false;
        button.textContent = '✓';
    }
}

async function markAsRead(notificationId, button) {
    if (!notificationManager) {
        console.error('Notification manager not initialized');
        return;
    }

    const originalText = button.textContent;
    button.disabled = true;
    button.textContent = 'Marking...';

    try {
        const success = await notificationManager.markAsRead(notificationId);
        
        if (success) {
            const notification = button.closest('.notification');
            if (notification) {
                notification.classList.remove('unread');
                notification.classList.add('read');
                
                const actionsDiv = button.closest('.notification-actions');
                if (actionsDiv) {
                    actionsDiv.remove();
                }
            }
        } else {
            throw new Error('Failed to mark as read');
        }
    } catch (error) {
        console.error('Error marking notification as read:', error);
        alert('Failed to mark notification as read. Please try again.');
        
        button.disabled = false;
        button.textContent = originalText;
    }
}

// Initialize when DOM is loaded
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initNotifications);
} else {
    initNotifications();
}