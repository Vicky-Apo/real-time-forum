// websocket/WebSocketManager.js - WebSocket Connection Manager

import state from '../state.js';

class WebSocketManager {
    constructor() {
        // Use relative URL - will work with Nginx proxy
        this.url = `ws://${window.location.host}/ws`;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 10;
        this.reconnectDelay = 1000;
        this.isConnecting = false;
        this.shouldReconnect = true;

        console.log('[WS] Manager initialized with URL:', this.url);
    }

    connect() {
        if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
            console.log('[WS] Already connected or connecting');
            return;
        }

        this.isConnecting = true;
        console.log('[WS] Connecting...');

        try {
            this.ws = new WebSocket(this.url);
            this.setupEventHandlers();
        } catch (error) {
            console.error('[WS] Connection error:', error);
            this.isConnecting = false;
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
            console.log('[WS] Message received:', message.type);

            // Emit event based on message type
            state.emit(`ws:${message.type}`, message.payload);

            // Handle specific message types
            this.handleMessage(message);
        } catch (error) {
            console.error('[WS] Failed to parse message:', error);
        }
    }

    handleMessage(message) {
        const { type, payload } = message;

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

            case 'new_message':
                state.emit('message:received', payload);
                this.showBrowserNotification(payload);
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

        console.log(`[WS] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);

        setTimeout(() => {
            this.connect();
        }, delay);
    }

    send(type, payload) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message = JSON.stringify({ type, payload });
            this.ws.send(message);
            console.log('[WS] Sent:', type);
        } else {
            console.warn('[WS] Cannot send, not connected');
        }
    }

    disconnect() {
        console.log('[WS] Disconnecting...');
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
