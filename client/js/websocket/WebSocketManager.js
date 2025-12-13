// websocket/WebSocketManager.js - WebSocket Connection Manager

import state from '../state.js';

class WebSocketManager {
    constructor() {
        // Connect to WebSocket through the same origin (will be proxied to backend)
        // This ensures cookies and CORS work correctly
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.host; // includes port
        this.url = `${protocol}//${host}/ws`;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 10;
        this.reconnectDelay = 1000;
        this.isConnecting = false;
        this.shouldReconnect = true;
    }

    connect() {
        if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
            return;
        }

        // Don't block if WebSocket is not supported
        if (typeof WebSocket === 'undefined') {
            console.warn('[WS] WebSocket not supported in this browser');
            return;
        }

        this.isConnecting = true;

        try {
            this.ws = new WebSocket(this.url);
            this.setupEventHandlers();
        } catch (error) {
            console.error('[WS] Connection error:', error);
            this.isConnecting = false;
            // Don't schedule reconnect if it's a fundamental error (like invalid URL)
            if (error.message && error.message.includes('Invalid')) {
                console.error('[WS] Invalid WebSocket URL, not attempting reconnect');
                return;
            }
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
        console.log('[WS] Connected successfully');
        this.isConnecting = false;
        this.reconnectAttempts = 0;
        this.reconnectDelay = 1000;

        state.setWsConnected(true);
        state.emit('ws:connected');
    }

    onMessage(event) {
        try {
            const message = JSON.parse(event.data);
            // Backend sends 'event' field, not 'type'
            const eventType = message.event || message.type;

            // Emit event based on message type
            state.emit(`ws:${eventType}`, message.payload);

            // Handle specific message types
            this.handleMessage(eventType, message.payload);
        } catch (error) {
            console.error('[WS] Failed to parse message:', error);
        }
    }

    handleMessage(type, payload) {
        switch (type) {
            case 'user_online':
                state.addOnlineUser(payload);
                break;

            case 'user_offline':
                state.removeOnlineUser(payload.user_id);
                break;

            case 'typing_start':
                state.emit('typing:start', payload);
                break;

            case 'typing_stop':
                state.emit('typing:stop', payload);
                break;

            case 'receive_message':
            case 'new_message':
                state.emit('message:received', payload);
                this.showBrowserNotification(payload);

                // Increment unread message count
                // Note: ChatView will decrement this when user opens the conversation
                state.incrementUnreadMessageCount();
                break;

            case 'message_read':
                state.emit('message:read', payload);
                break;

            case 'notification':
                state.emit('notification:received', payload);
                state.incrementUnreadCount();
                break;

            default:
                console.warn('[WS] Unknown message type:', type);
        }
    }

    onError(error) {
        console.error('[WS] Error:', error);
        // Don't let WebSocket errors block the application
        // The connection will be handled by onClose
    }

    onClose(event) {
        console.log('[WS] Disconnected -', event.code, event.reason);
        this.isConnecting = false;
        this.ws = null;

        state.setWsConnected(false);
        state.emit('ws:disconnected');

        if (this.shouldReconnect && state.getUser()) {
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
        const delay = Math.min(
            this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1),
            30000
        );

        setTimeout(() => {
            this.connect();
        }, delay);
    }

    send(type, payload) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            // Backend expects 'event' field, not 'type'
            const message = JSON.stringify({ event: type, payload });
            this.ws.send(message);
        } else {
            console.warn('[WS] Cannot send, not connected');
        }
    }

    disconnect() {
        this.shouldReconnect = false;

        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }

        state.setWsConnected(false);
    }

    showBrowserNotification(message) {
        if ('Notification' in window && Notification.permission === 'granted') {
            new Notification('New Message', {
                body: `${message.sender_name}: ${message.content}`,
                icon: '/assets/logo.png',
                tag: 'message-notification'
            });
        }
    }
}

// Create and export singleton
const wsManager = new WebSocketManager();
export default wsManager;
