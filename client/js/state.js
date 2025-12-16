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

        // Messages
        this.unreadMessageCount = 0;

        // Event listeners (pub/sub pattern)
        this.listeners = new Map();
    }

    // ============= User Management =============

    getUser() {
        if (!this.user) {
            // Try to load from localStorage
            const stored = localStorage.getItem('user');
            if (stored) {
                try {
                    this.user = JSON.parse(stored);
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
        } else {
            localStorage.removeItem('user');
        }

        this.emit('user:changed', user);
    }

    // ============= Online Users =============

    setOnlineUsers(users) {
        this.onlineUsers = users;
        this.emit('users:online', users);
    }

    addOnlineUser(user) {
        if (!this.onlineUsers.find(u => u.user_id === user.user_id)) {
            this.onlineUsers.push(user);
            this.emit('user:online', user);
            this.emit('users:online', this.onlineUsers);
        }
    }

    removeOnlineUser(userId) {
        const user = this.onlineUsers.find(u => u.user_id === userId);
        this.onlineUsers = this.onlineUsers.filter(u => u.user_id !== userId);

        if (user) {
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
        this.emit('unread:changed', count);
    }

    incrementUnreadCount() {
        this.unreadCount++;
        this.emit('unread:changed', this.unreadCount);
    }

    // ============= Messages =============

    setUnreadMessageCount(count) {
        this.unreadMessageCount = count;
        this.emit('unread-messages:changed', count);
    }

    incrementUnreadMessageCount() {
        this.unreadMessageCount++;
        this.emit('unread-messages:changed', this.unreadMessageCount);
    }

    // ============= WebSocket State =============

    setWsConnected(connected) {
        this.wsConnected = connected;
        this.emit('ws:connection', connected);
    }

    // ============= Event System (Pub/Sub) =============

    on(event, callback) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, []);
        }
        this.listeners.get(event).push(callback);
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
        this.setUser(null);  // This will emit 'user:changed' event
        this.wsConnected = false;
        this.onlineUsers = [];
        this.unreadCount = 0;
        this.unreadMessageCount = 0;
    }
}

// Create singleton instance
const state = new State();

// Export the singleton
export default state;
