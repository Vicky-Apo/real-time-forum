// state.js - Global Application State Management

class State {
    constructor() {
        // User state
        this.user = null;

        // WebSocket state
        this.wsConnected = false;

        // Online users
        this.onlineUsers = [];

        // Notifications
        this.unreadCount = 0;

        // Event listeners (pub/sub pattern)
        this.listeners = new Map();

        console.log('[State] State manager initialized');
    }

    // ============= User Management =============

    getUser() {
        if (!this.user) {
            // Try to load from localStorage
            const stored = localStorage.getItem('user');
            if (stored) {
                try {
                    this.user = JSON.parse(stored);
                    console.log('[State] User loaded from localStorage:', this.user.username);
                } catch (e) {
                    console.error('[State] Error parsing stored user:', e);
                    localStorage.removeItem('user');
                }
            }
        }
        return this.user;
    }

    setUser(user) {
        this.user = user;

        if (user) {
            localStorage.setItem('user', JSON.stringify(user));
            console.log('[State] User set:', user.username);
        } else {
            localStorage.removeItem('user');
            console.log('[State] User cleared');
        }

        this.emit('user:changed', user);
    }

    // ============= Online Users =============

    setOnlineUsers(users) {
        this.onlineUsers = users;
        console.log('[State] Online users set:', users.length);
        this.emit('users:online', users);
    }

    addOnlineUser(user) {
        if (!this.onlineUsers.find(u => u.user_id === user.user_id)) {
            this.onlineUsers.push(user);
            console.log('[State] User came online:', user.username);
            this.emit('user:online', user);
            this.emit('users:online', this.onlineUsers);
        }
    }

    removeOnlineUser(userId) {
        const user = this.onlineUsers.find(u => u.user_id === userId);
        this.onlineUsers = this.onlineUsers.filter(u => u.user_id !== userId);

        if (user) {
            console.log('[State] User went offline:', user.username);
            this.emit('user:offline', userId);
            this.emit('users:online', this.onlineUsers);
        }
    }

    getOnlineUsers() {
        return this.onlineUsers;
    }

    // ============= Notifications =============

    setUnreadCount(count) {
        this.unreadCount = count;
        console.log('[State] Unread count set to:', count);
        this.emit('unread:changed', count);
    }

    incrementUnreadCount() {
        this.unreadCount++;
        this.emit('unread:changed', this.unreadCount);
    }

    // ============= WebSocket State =============

    setWsConnected(connected) {
        this.wsConnected = connected;
        console.log('[State] WebSocket connected:', connected);
        this.emit('ws:connection', connected);
    }

    // ============= Event System (Pub/Sub) =============

    on(event, callback) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, []);
        }
        this.listeners.get(event).push(callback);
        // console.log('[State] Listener added for event:', event);
    }

    off(event, callback) {
        const callbacks = this.listeners.get(event);
        if (callbacks) {
            this.listeners.set(
                event,
                callbacks.filter(cb => cb !== callback)
            );
        }
    }

    emit(event, data) {
        const callbacks = this.listeners.get(event);
        if (callbacks) {
            callbacks.forEach(callback => {
                try {
                    callback(data);
                } catch (error) {
                    console.error(`[State] Error in event listener for '${event}':`, error);
                }
            });
        }
    }

    // ============= Utility =============

    clear() {
        this.user = null;
        this.wsConnected = false;
        this.onlineUsers = [];
        this.unreadCount = 0;
        localStorage.removeItem('user');
        console.log('[State] State cleared');
    }
}

// Create singleton instance
const state = new State();

// Export the singleton
export default state;
