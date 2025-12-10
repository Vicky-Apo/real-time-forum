// views/ChatView.js - Chat/Messaging Page with Real-time Support

import apiClient from '../api/client.js';
import state from '../state.js';
import wsManager from '../websocket/WebSocketManager.js';

export default {
    conversations: [],
    currentConversation: null,
    currentUser: null,
    messages: [],
    onlineUsers: new Set(),
    typingTimeout: null,

    async render() {
        return `
            <link rel="stylesheet" href="/css/chat.css">
            <div class="container">
                <div class="chat-container">
                    <!-- Conversations Sidebar -->
                    <div class="conversations-sidebar">
                        <div class="conversations-header">
                            <h2>💬 Messages</h2>
                        </div>
                        <div id="conversations-list" class="conversations-list">
                            <div class="loading-container">
                                <div class="loading-spinner"></div>
                                <p>Loading conversations...</p>
                            </div>
                        </div>
                    </div>

                    <!-- Chat Window -->
                    <div class="chat-window">
                        <div id="chat-content">
                            <div class="chat-empty-state">
                                <div class="chat-empty-state-icon">💬</div>
                                <h3>Select a conversation</h3>
                                <p>Choose a user from the list to start messaging</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
    },

    async afterRender() {
        console.log('[ChatView] Rendered');
        this.currentUser = state.getUser();

        try {
            // Load conversations (includes all users with online status)
            await this.loadConversations();

            // Setup WebSocket listeners
            this.setupWebSocketListeners();

            // Setup event listeners
            this.setupEventListeners();
        } catch (error) {
            console.error('[ChatView] Error initializing:', error);
        }
    },

    async loadConversations() {
        const container = document.getElementById('conversations-list');

        try {
            // Get conversations (now includes all users with online status)
            const response = await apiClient.get('/conversations');
            console.log('[ChatView] Raw API response:', response);

            // API client already normalizes the response
            this.conversations = response.conversations || [];

            // Build online users set from conversations
            this.onlineUsers = new Set(
                this.conversations
                    .filter(conv => conv.is_online)
                    .map(conv => conv.user_id)
            );

            console.log('[ChatView] Loaded conversations:', this.conversations.length);
            console.log('[ChatView] Online users:', this.onlineUsers.size);
            console.log('[ChatView] Container element:', container);

            if (this.conversations.length === 0) {
                container.innerHTML = `
                    <div class="empty-state" style="padding: var(--space-3xl) var(--space-2xl);">
                        <p class="empty-state-message">No users available</p>
                        <p style="font-size: var(--font-size-sm); margin-top: var(--space-md);">
                            No other users found in the system.
                        </p>
                    </div>
                `;
                return;
            }

            const html = this.conversations.map(conv => this.renderConversation(conv)).join('');
            console.log('[ChatView] Generated HTML length:', html.length);
            console.log('[ChatView] First 200 chars of HTML:', html.substring(0, 200));
            container.innerHTML = html;
        } catch (error) {
            console.error('[ChatView] Error loading conversations:', error);
            container.innerHTML = `
                <div class="empty-state" style="padding: var(--space-3xl) var(--space-2xl);">
                    <p class="empty-state-message">No conversations yet</p>
                    <p style="font-size: var(--font-size-sm); margin-top: var(--space-md);">
                        Start messaging other users to see them here!
                    </p>
                </div>
            `;
        }
    },

    renderConversation(conversation) {
        const isOnline = this.onlineUsers.has(conversation.user_id);
        const unreadBadge = conversation.unread_count > 0
            ? `<span class="conversation-unread">${conversation.unread_count}</span>`
            : '';

        const initials = conversation.username?.slice(0, 2).toUpperCase() || 'U';

        // Format last message properly
        let lastMessageText = 'No messages yet';
        if (conversation.last_message && conversation.last_message.content) {
            const prefix = conversation.last_message.is_from_me ? 'You: ' : '';
            lastMessageText = prefix + conversation.last_message.content;
            // Truncate if too long
            if (lastMessageText.length > 50) {
                lastMessageText = lastMessageText.substring(0, 50) + '...';
            }
        }

        return `
            <div class="conversation-item" data-user-id="${conversation.user_id}" onclick="selectConversation('${conversation.user_id}')">
                <div class="conversation-avatar">${initials}</div>
                <div class="conversation-info">
                    <div class="conversation-name">
                        ${conversation.username}
                        <span class="${isOnline ? 'online-indicator' : 'offline-indicator'}"></span>
                    </div>
                    <div class="conversation-preview">
                        ${lastMessageText}
                    </div>
                </div>
                ${unreadBadge}
            </div>
        `;
    },

    async selectConversation(userId) {
        console.log('[ChatView] Selecting conversation:', userId);

        // Find conversation
        this.currentConversation = this.conversations.find(c => c.user_id === userId);

        if (!this.currentConversation) {
            console.error('[ChatView] Conversation not found for user:', userId);
            return;
        }

        // Update UI
        document.querySelectorAll('.conversation-item').forEach(item => {
            item.classList.remove('active');
        });
        document.querySelector(`[data-user-id="${userId}"]`)?.classList.add('active');

        // Load messages
        await this.loadMessages(userId);
    },

    async loadMessages(userId) {
        const chatContent = document.getElementById('chat-content');

        try {
            const response = await apiClient.get(`/messages/${userId}`);
            const data = response.data || response;
            // Reverse messages so oldest appears first (backend sends newest first)
            this.messages = (data.messages || []).reverse();

            const isOnline = this.onlineUsers.has(userId);

            chatContent.innerHTML = `
                <div class="chat-header">
                    <div class="chat-header-info">
                        <div class="chat-header-avatar">
                            ${this.currentConversation.username?.slice(0, 2).toUpperCase() || 'U'}
                        </div>
                        <div class="chat-header-details">
                            <h3>${this.currentConversation.username}</h3>
                            <div class="chat-header-status">
                                <span class="${isOnline ? 'online-indicator' : 'offline-indicator'}"></span>
                                ${isOnline ? 'Online' : 'Offline'}
                            </div>
                        </div>
                    </div>
                </div>

                <div class="messages-area" id="messages-area">
                    ${this.messages.length > 0
                        ? this.messages.map(msg => this.renderMessage(msg)).join('')
                        : '<div class="chat-empty-state"><p>No messages yet. Start the conversation!</p></div>'
                    }
                    <div id="typing-indicator"></div>
                </div>

                <div class="message-input-area">
                    <form class="message-input-form" id="message-form">
                        <div class="message-input-wrapper">
                            <textarea
                                id="message-input"
                                class="message-input"
                                placeholder="Type a message..."
                                rows="1"
                            ></textarea>
                        </div>
                        <div class="message-actions">
                            <label class="btn-icon" for="image-upload" title="Send image">
                                📷
                                <input type="file" id="image-upload" accept="image/*" style="display: none;">
                            </label>
                            <button type="submit" class="btn btn-send">Send</button>
                        </div>
                    </form>
                </div>
            `;

            // Scroll to bottom
            this.scrollToBottom();

            // Setup message form
            this.setupMessageForm();
            this.setupTypingIndicator();

        } catch (error) {
            console.error('[ChatView] Error loading messages:', error);

            // Even if loading messages fails, show the chat interface so user can send first message
            const isOnline = this.onlineUsers.has(userId);

            chatContent.innerHTML = `
                <div class="chat-header">
                    <div class="chat-header-info">
                        <div class="chat-header-avatar">
                            ${this.currentConversation.username?.slice(0, 2).toUpperCase() || 'U'}
                        </div>
                        <div class="chat-header-details">
                            <h3>${this.currentConversation.username}</h3>
                            <div class="chat-header-status">
                                <span class="${isOnline ? 'online-indicator' : 'offline-indicator'}"></span>
                                ${isOnline ? 'Online' : 'Offline'}
                            </div>
                        </div>
                    </div>
                </div>

                <div class="messages-area" id="messages-area">
                    <div class="chat-empty-state">
                        <p>Start a new conversation with ${this.currentConversation.username}!</p>
                    </div>
                    <div id="typing-indicator"></div>
                </div>

                <div class="message-input-area">
                    <form class="message-input-form" id="message-form">
                        <div class="message-input-wrapper">
                            <textarea
                                id="message-input"
                                class="message-input"
                                placeholder="Type a message..."
                                rows="1"
                            ></textarea>
                        </div>
                        <div class="message-actions">
                            <label class="btn-icon" for="image-upload" title="Send image">
                                📷
                                <input type="file" id="image-upload" accept="image/*" style="display: none;">
                            </label>
                            <button type="submit" class="btn btn-send">Send</button>
                        </div>
                    </form>
                </div>
            `;

            // Setup message form even on error
            this.setupMessageForm();
            this.setupTypingIndicator();
        }
    },

    renderMessage(message) {
        const isSent = message.sender_id === this.currentUser.user_id;
        const time = new Date(message.created_at).toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit'
        });
        const initials = message.sender_username?.slice(0, 2).toUpperCase() || 'U';

        let imageHTML = '';
        if (message.images && message.images.length > 0) {
            imageHTML = message.images.map(img => `
                <img src="${img.image_url}"
                     alt="Message image"
                     class="message-image"
                     onclick="showImageLightbox('${img.image_url}')"
                     style="max-width: 200px; border-radius: 8px; margin-top: 8px; cursor: pointer;"
                >
            `).join('');
        }

        return `
            <div class="message ${isSent ? 'sent' : 'received'}">
                <div class="message-avatar">${initials}</div>
                <div class="message-content-wrapper">
                    <div class="message-content">
                        ${message.content}
                        ${imageHTML}
                    </div>
                    <div class="message-time">${time}</div>
                </div>
            </div>
        `;
    },

    setupMessageForm() {
        const form = document.getElementById('message-form');
        const input = document.getElementById('message-input');
        const imageInput = document.getElementById('image-upload');

        if (!form || !input) return;

        // Auto-resize textarea
        input.addEventListener('input', () => {
            input.style.height = 'auto';
            input.style.height = input.scrollHeight + 'px';
        });

        // Handle form submission
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            await this.sendMessage();
        });

        // Handle image upload
        imageInput.addEventListener('change', async () => {
            if (imageInput.files.length > 0) {
                await this.sendMessage(imageInput.files[0]);
                imageInput.value = '';
            }
        });

        // Focus input
        input.focus();
    },

    setupTypingIndicator() {
        const input = document.getElementById('message-input');
        if (!input) return;

        input.addEventListener('input', () => {
            // Send typing indicator
            if (this.currentConversation) {
                wsManager.send('typing_start', {
                    recipient_id: this.currentConversation.user_id
                });

                // Clear previous timeout
                if (this.typingTimeout) {
                    clearTimeout(this.typingTimeout);
                }

                // Set timeout to stop typing
                this.typingTimeout = setTimeout(() => {
                    wsManager.send('typing_stop', {
                        recipient_id: this.currentConversation.user_id
                    });
                }, 2000);
            }
        });
    },

    async sendMessage(imageFile = null) {
        const input = document.getElementById('message-input');
        const content = input.value.trim();

        if (!content && !imageFile) {
            return;
        }

        if (!this.currentConversation) {
            this.showToast('Please select a conversation first', 'warning');
            return;
        }

        try {
            const formData = new FormData();
            formData.append('recipient_id', this.currentConversation.user_id);

            if (content) {
                formData.append('content', content);
            }

            if (imageFile) {
                formData.append('images', imageFile);
            }

            await apiClient.post('/messages/send', formData);

            // Clear input
            input.value = '';
            input.style.height = 'auto';

            // Reload messages
            await this.loadMessages(this.currentConversation.user_id);

        } catch (error) {
            console.error('[ChatView] Error sending message:', error);
            this.showToast(error.message || 'Failed to send message', 'error');
        }
    },

    setupWebSocketListeners() {
        // Listen for new messages
        state.on('message:received', (message) => {
            console.log('[ChatView] New message received:', message);

            // Normalize WebSocket payload to match database format
            const normalizedMessage = {
                sender_id: message.sender_id,
                sender_username: message.sender_name,  // WebSocket uses sender_name
                content: message.content,
                created_at: message.sent_at,  // WebSocket uses sent_at
                images: message.images || []
            };

            // If message is for current conversation, add it
            if (this.currentConversation &&
                (message.sender_id === this.currentConversation.user_id)) {
                this.addMessageToUI(normalizedMessage);
            }

            // Reload conversations to update preview
            this.loadConversations();
        });

        // Listen for typing indicators
        state.on('typing:start', (data) => {
            if (this.currentConversation && data.user_id === this.currentConversation.user_id) {
                this.showTypingIndicator();
            }
        });

        state.on('typing:stop', (data) => {
            if (this.currentConversation && data.user_id === this.currentConversation.user_id) {
                this.hideTypingIndicator();
            }
        });

        // Listen for online/offline status
        state.on('ws:user_online', (data) => {
            this.onlineUsers.add(data.user_id);
            this.updateOnlineStatus();
        });

        state.on('ws:user_offline', (data) => {
            this.onlineUsers.delete(data.user_id);
            this.updateOnlineStatus();
        });
    },

    addMessageToUI(message) {
        const messagesArea = document.getElementById('messages-area');
        if (!messagesArea) return;

        // Remove empty state if exists
        const emptyState = messagesArea.querySelector('.chat-empty-state');
        if (emptyState) {
            emptyState.remove();
        }

        // Add message
        const typingIndicator = document.getElementById('typing-indicator');
        const messageHTML = this.renderMessage(message);
        typingIndicator.insertAdjacentHTML('beforebegin', messageHTML);

        // Scroll to bottom
        this.scrollToBottom();
    },

    showTypingIndicator() {
        const indicator = document.getElementById('typing-indicator');
        if (indicator && !indicator.querySelector('.typing-indicator')) {
            indicator.innerHTML = `
                <div class="typing-indicator">
                    <div class="typing-dots">
                        <div class="typing-dot"></div>
                        <div class="typing-dot"></div>
                        <div class="typing-dot"></div>
                    </div>
                </div>
            `;
            this.scrollToBottom();
        }
    },

    hideTypingIndicator() {
        const indicator = document.getElementById('typing-indicator');
        if (indicator) {
            indicator.innerHTML = '';
        }
    },

    updateOnlineStatus() {
        // Update conversation list
        this.loadConversations();

        // Update chat header if conversation is open
        if (this.currentConversation) {
            const isOnline = this.onlineUsers.has(this.currentConversation.user_id);
            const statusElement = document.querySelector('.chat-header-status');
            if (statusElement) {
                statusElement.innerHTML = `
                    <span class="${isOnline ? 'online-indicator' : 'offline-indicator'}"></span>
                    ${isOnline ? 'Online' : 'Offline'}
                `;
            }
        }
    },

    scrollToBottom() {
        setTimeout(() => {
            const messagesArea = document.getElementById('messages-area');
            if (messagesArea) {
                messagesArea.scrollTop = messagesArea.scrollHeight;
            }
        }, 100);
    },

    setupEventListeners() {
        // Make selectConversation available globally
        window.selectConversation = this.selectConversation.bind(this);
        window.showImageLightbox = this.showImageLightbox.bind(this);
    },

    showImageLightbox(imageUrl) {
        let lightbox = document.getElementById('image-lightbox');
        if (!lightbox) {
            lightbox = document.createElement('div');
            lightbox.id = 'image-lightbox';
            lightbox.innerHTML = `
                <button class="close-btn" onclick="this.parentElement.classList.remove('active')">×</button>
                <img src="" alt="Full size image">
            `;
            document.body.appendChild(lightbox);
        }

        const img = lightbox.querySelector('img');
        img.src = imageUrl;
        lightbox.classList.add('active');
    },

    showToast(message, type = 'info') {
        let container = document.querySelector('.toast-container');
        if (!container) {
            container = document.createElement('div');
            container.className = 'toast-container';
            document.body.appendChild(container);
        }

        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        toast.innerHTML = `
            <div class="toast-content">
                <div class="toast-message">${message}</div>
            </div>
            <button class="toast-close" onclick="this.parentElement.remove()">×</button>
        `;

        container.appendChild(toast);

        setTimeout(() => {
            toast.classList.add('removing');
            setTimeout(() => toast.remove(), 300);
        }, 3000);
    }
};
