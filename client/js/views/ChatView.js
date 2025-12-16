// views/ChatView.js - Chat/Messaging Page with Real-time Support

import apiClient from '../api/client.js';
import state from '../state.js';
import wsManager from '../websocket/WebSocketManager.js';
import { getInitials, showImageLightbox, showToast, formatTime, throttle } from '../utils/helpers.js';

export default {
    conversations: [],
    currentConversation: null,
    currentUser: null,
    messages: [],
    onlineUsers: new Set(),
    typingTimeout: null,
    isLoadingOlderMessages: false,
    hasMoreMessages: true,
    oldestMessageTimestamp: null,

    async render() {
        return `
            <link rel="stylesheet" href="/css/chat.css">
            <div class="container">
                <div class="chat-container">
                    <!-- Conversations Sidebar -->
                    <div class="conversations-sidebar">
                        <div class="conversations-header">
                            <h2><i class="fas fa-comments"></i> Messages</h2>
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
                                <div class="chat-empty-state-icon">
                                    <i class="fas fa-comments"></i>
                                </div>
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
        this.currentUser = state.getUser();

        // Normalize user object to always have user_id field
        if (this.currentUser && !this.currentUser.user_id && this.currentUser.id) {
            this.currentUser.user_id = this.currentUser.id;
        }

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

        // If container doesn't exist, user has navigated away - skip update
        if (!container) {
            return;
        }

        try {
            // Get conversations (now includes all users with online status)
            const response = await apiClient.get('/conversations');

            // API client already normalizes the response
            this.conversations = response.conversations || [];

            // Build online users set from conversations
            this.onlineUsers = new Set(
                this.conversations
                    .filter(conv => conv.is_online)
                    .map(conv => conv.user_id)
            );

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

        const initials = getInitials(conversation.username);

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
            // Load more messages initially (50 instead of default 10)
            const response = await apiClient.get(`/messages/${userId}?limit=50`);
            const data = response.data || response;
            // Reverse messages so oldest appears first (backend sends newest first)
            this.messages = (data.messages || []).reverse();

            const isOnline = this.onlineUsers.has(userId);

            chatContent.innerHTML = `
                <div class="chat-header">
                    <div class="chat-header-info">
                        <div class="chat-header-avatar">
                            ${getInitials(this.currentConversation.username)}
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
                        <div class="message-input-row">
                        <div class="message-input-wrapper">
                            <textarea
                                id="message-input"
                                class="message-input"
                                placeholder="Type a message..."
                                rows="1"
                            ></textarea>
                            </div>
                            <button type="submit" class="btn btn-send">Send</button>
                        </div>
                        <div class="message-actions">
                            <label class="btn-upload-image" for="image-upload">
                                <span class="upload-icon"><i class="fas fa-camera"></i></span>
                                Upload Image
                            </label>
                            <input type="file" id="image-upload" accept="image/*" style="display: none;">
                        </div>
                    </form>
                </div>
            `;

            // Reset infinite scroll state for new conversation
            this.hasMoreMessages = true;
            this.oldestMessageTimestamp = this.messages[0]?.created_at || null;

            // Scroll to bottom
            this.scrollToBottom();

            // Setup message form
            this.setupMessageForm();
            this.setupTypingIndicator();
            this.setupInfiniteScroll();

            // Fetch updated unread count after marking messages as read
            // (backend marks messages as read when we fetch them)
            await this.updateUnreadCount();

        } catch (error) {
            console.error('[ChatView] Error loading messages:', error);

            // Even if loading messages fails, show the chat interface so user can send first message
            const isOnline = this.onlineUsers.has(userId);

            chatContent.innerHTML = `
                <div class="chat-header">
                    <div class="chat-header-info">
                        <div class="chat-header-avatar">
                            ${getInitials(this.currentConversation.username)}
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
                        <div class="message-input-row">
                        <div class="message-input-wrapper">
                            <textarea
                                id="message-input"
                                class="message-input"
                                placeholder="Type a message..."
                                rows="1"
                            ></textarea>
                            </div>
                            <button type="submit" class="btn btn-send">Send</button>
                        </div>
                        <div class="message-actions">
                            <label class="btn-upload-image" for="image-upload">
                                <span class="upload-icon"><i class="fas fa-camera"></i></span>
                                Upload Image
                            </label>
                            <input type="file" id="image-upload" accept="image/*" style="display: none;">
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
        // Debug logging for every message to understand the issue
        const isSent = message.sender_id === this.currentUser.user_id;
        const time = formatTime(message.created_at);

        // Get initials - use sender_username if available, otherwise use conversation username for received messages
        let initials;
        if (message.sender_username) {
            initials = getInitials(message.sender_username);
        } else if (!isSent && this.currentConversation) {
            // For received messages without sender_username, use the conversation username
            initials = getInitials(this.currentConversation.username);
        } else if (isSent && this.currentUser) {
            // For sent messages without sender_username, use current user's username
            initials = getInitials(this.currentUser.username);
        } else {
            initials = 'U';
        }

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

        // Handle Enter key to send (Shift+Enter for new line)
        input.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                form.dispatchEvent(new Event('submit'));
            }
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

    setupInfiniteScroll() {
        const messagesArea = document.getElementById('messages-area');
        if (!messagesArea) {
            return;
        }

        // Throttle the scroll event handler to prevent spamming (required by project specs)
        // This ensures loadOlderMessages is called at most once every 500ms
        const throttledScrollHandler = throttle(async () => {
            // Check if scrolled near top (within 50px threshold for easier triggering)
            const isNearTop = messagesArea.scrollTop < 50;

            if (isNearTop && !this.isLoadingOlderMessages && this.hasMoreMessages) {
                await this.loadOlderMessages();
            }
        }, 500); // 500ms throttle delay to prevent scroll event spam

        messagesArea.addEventListener('scroll', throttledScrollHandler);
    },

    async loadOlderMessages() {
        if (!this.currentConversation || this.isLoadingOlderMessages || !this.hasMoreMessages) {
            return;
        }

        this.isLoadingOlderMessages = true;
        const messagesArea = document.getElementById('messages-area');

        // Store current scroll height to restore scroll position
        const previousScrollHeight = messagesArea.scrollHeight;

        try {
            // Make API call with before parameter for pagination (using timestamp)
            const response = await apiClient.get(`/messages/${this.currentConversation.user_id}?before=${encodeURIComponent(this.oldestMessageTimestamp)}&limit=10`);
            const data = response.data || response;
            const olderMessages = (data.messages || []).reverse();

            if (olderMessages.length === 0) {
                this.hasMoreMessages = false;
                return;
            }

            // Update oldest message timestamp
            this.oldestMessageTimestamp = olderMessages[0]?.created_at;

            // Prepend messages to current messages array
            this.messages = [...olderMessages, ...this.messages];

            // Render older messages at the top
            const olderMessagesHTML = olderMessages.map(msg => this.renderMessage(msg)).join('');
            const typingIndicator = document.getElementById('typing-indicator');

            // Insert before the first message (after any typing indicator parent)
            const firstMessage = messagesArea.querySelector('.message');
            if (firstMessage) {
                firstMessage.insertAdjacentHTML('beforebegin', olderMessagesHTML);
            } else if (typingIndicator) {
                typingIndicator.insertAdjacentHTML('beforebegin', olderMessagesHTML);
            } else {
                messagesArea.insertAdjacentHTML('afterbegin', olderMessagesHTML);
            }

            // Restore scroll position
            const newScrollHeight = messagesArea.scrollHeight;
            messagesArea.scrollTop = newScrollHeight - previousScrollHeight;

        } catch (error) {
            console.error('[ChatView] Error loading older messages:', error);
        } finally {
            this.isLoadingOlderMessages = false;
        }
    },

    async sendMessage(imageFile = null) {
        const input = document.getElementById('message-input');
        const content = input.value.trim();

        if (!content && !imageFile) {
            return;
        }

        if (!this.currentConversation) {
            showToast('Please select a conversation first', 'warning');
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

            const response = await apiClient.post('/messages/send', formData);

            // Clear input
            input.value = '';
            input.style.height = 'auto';

            // Add sent message to UI instead of reloading all messages
            const sentMessage = response.data || response;
            if (sentMessage && sentMessage.message_id) {
                // Normalize the sent message
                const normalizedMessage = {
                    sender_id: this.currentUser.user_id,
                    sender_username: this.currentUser.username,
                    content: sentMessage.content || content,
                    created_at: sentMessage.created_at || new Date().toISOString(),
                    message_id: sentMessage.message_id,
                    images: sentMessage.images || []
                };
                this.messages.push(normalizedMessage);
                this.addMessageToUI(normalizedMessage);
            }

        } catch (error) {
            console.error('[ChatView] Error sending message:', error);
            showToast(error.message || 'Failed to send message', 'error');
        }
    },

    setupWebSocketListeners() {
        // Listen for new messages
        state.on('message:received', (message) => {

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

    async updateUnreadCount() {
        try {
            const response = await apiClient.get('/messages/unread-count');
            const data = response.data || response;
            const unreadCount = data.unread_count || 0;
            state.setUnreadMessageCount(unreadCount);

            // Also reload conversations to update badges in sidebar
            await this.loadConversations();
        } catch (error) {
            console.error('[ChatView] Error updating unread message count:', error);
        }
    },

    setupEventListeners() {
        // Make selectConversation and showImageLightbox available globally for inline onclick handlers
        window.selectConversation = this.selectConversation.bind(this);
        window.showImageLightbox = showImageLightbox;
    }
};
