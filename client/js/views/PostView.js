// views/PostView.js - Single Post View with Comments

import apiClient from '../api/client.js';
import state from '../state.js';

export default {
    postId: null,
    post: null,
    comments: [],

    async render(params) {
        this.postId = params.id;

        return `
            <div class="container">
                <div id="post-content">
                    <div class="loading-container">
                        <div class="loading-spinner"></div>
                        <p>Loading post...</p>
                    </div>
                </div>
            </div>

            <!-- Image Lightbox -->
            <div id="image-lightbox">
                <button class="close-btn" onclick="this.parentElement.classList.remove('active')">×</button>
                <img src="" alt="Full size image">
            </div>
        `;
    },

    async afterRender() {
        console.log('[PostView] Rendered for post:', this.postId);

        try {
            await this.loadPost();
            await this.loadComments();
            this.setupEventListeners();
        } catch (error) {
            console.error('[PostView] Error loading post:', error);
            this.showError('Failed to load post');
        }
    },

    async loadPost() {
        const container = document.getElementById('post-content');

        try {
            const response = await apiClient.get(`/posts/view/${this.postId}`);
            this.post = response.data || response;

            container.innerHTML = this.renderPost(this.post);
        } catch (error) {
            console.error('[PostView] Error loading post:', error);
            container.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">⚠️</div>
                    <h3 class="empty-state-title">Failed to load post</h3>
                    <p class="empty-state-message">${error.message || 'Please try again later.'}</p>
                    <a href="/" data-link class="btn btn-primary" style="margin-top: 1rem;">Back to Home</a>
                </div>
            `;
            throw error;
        }
    },

    async loadComments() {
        const container = document.getElementById('comments-list');
        if (!container) return;

        try {
            const response = await apiClient.get(`/comments/for-post/${this.postId}`);
            const data = response.data || response;
            this.comments = data.comments || [];

            if (this.comments.length === 0) {
                container.innerHTML = `
                    <div class="no-comments">
                        <h3>No comments yet</h3>
                        <p>Be the first to share your thoughts!</p>
                    </div>
                `;
                return;
            }

            container.innerHTML = this.comments.map(comment => this.renderComment(comment)).join('');
        } catch (error) {
            console.error('[PostView] Error loading comments:', error);
            container.innerHTML = `
                <div class="empty-state">
                    <p class="empty-state-message">Failed to load comments</p>
                </div>
            `;
        }
    },

    renderPost(post) {
        const currentUser = state.getUser();
        const createdDate = new Date(post.created_at).toLocaleDateString();
        const categories = post.categories?.map(cat =>
            `<a href="/category/${cat.category_id}" data-link class="category-tag">${cat.category_name || cat.name}</a>`
        ).join(' ') || '';

        const isAuthor = currentUser && currentUser.user_id === post.author_id;
        const netVotes = (post.likes || 0) - (post.dislikes || 0);

        // Render images if available
        let imagesHTML = '';
        if (post.images && post.images.length > 0) {
            imagesHTML = post.images.map(img => `
                <img src="${img.image_url}"
                     alt="Post image"
                     class="post-image"
                     onclick="showImageLightbox('${img.image_url}')"
                >
            `).join('');
        }

        return `
            <!-- Back Button -->
            <div style="margin-bottom: var(--space-2xl);">
                <a href="/" data-link class="btn btn-secondary">← Back to Feed</a>
            </div>

            <!-- Main Post Card -->
            <article class="post-card main-post-card">
                <div class="post-header">
                    <div class="post-meta">
                        <span class="author-info">
                            <strong>u/${post.username || 'Anonymous'}</strong>
                        </span>
                        <span>•</span>
                        <span class="post-date">${createdDate}</span>
                    </div>
                    ${categories ? `
                        <div class="post-categories">
                            ${categories}
                        </div>
                    ` : ''}
                </div>

                <div class="post-content">
                    ${post.content}
                </div>

                ${imagesHTML ? `
                    <div class="post-images" style="margin-top: var(--space-2xl);">
                        ${imagesHTML}
                    </div>
                ` : ''}

                <!-- Reactions and Actions -->
                <div style="display: flex; justify-content: space-between; align-items: center; margin-top: var(--space-2xl); flex-wrap: wrap; gap: var(--space-lg);">
                    <div class="reaction-buttons">
                        <button class="reaction-btn like-btn" onclick="handlePostReaction('${post.post_id}', 'like')">
                            👍 Like (${post.likes || 0})
                        </button>
                        <button class="reaction-btn dislike-btn" onclick="handlePostReaction('${post.post_id}', 'dislike')">
                            👎 Dislike (${post.dislikes || 0})
                        </button>
                        <span style="font-weight: var(--font-weight-bold); color: var(--color-text-secondary);">
                            Net: ${netVotes}
                        </span>
                    </div>

                    ${isAuthor ? `
                        <div class="post-actions">
                            <a href="/edit-post/${post.post_id}" data-link class="edit-btn">
                                ✏️ Edit
                            </a>
                            <button class="delete-btn" onclick="handleDeletePost('${post.post_id}')">
                                🗑️ Delete
                            </button>
                        </div>
                    ` : ''}
                </div>
            </article>

            <!-- Comments Section -->
            <div class="comments-section">
                <div class="comments-header">
                    <h3>💬 Comments (${post.comment_count || 0})</h3>
                </div>

                <!-- Add Comment Form -->
                ${currentUser ? `
                    <div class="add-comment">
                        <form id="add-comment-form">
                            <div class="form-group">
                                <textarea
                                    id="comment-content"
                                    class="form-control"
                                    placeholder="Share your thoughts..."
                                    rows="4"
                                    required
                                ></textarea>
                            </div>
                            <div class="form-actions">
                                <button type="submit" class="btn btn-primary">
                                    Post Comment
                                </button>
                            </div>
                        </form>
                    </div>
                ` : `
                    <div class="alert alert-info">
                        <a href="/login" data-link>Login</a> to post a comment
                    </div>
                `}

                <!-- Comments List -->
                <div id="comments-list" class="comments-list">
                    <div class="loading-container">
                        <div class="loading-spinner"></div>
                        <p>Loading comments...</p>
                    </div>
                </div>
            </div>
        `;
    },

    renderComment(comment) {
        const currentUser = state.getUser();
        const createdDate = new Date(comment.created_at).toLocaleDateString();
        const isAuthor = currentUser && currentUser.user_id === comment.author_id;

        return `
            <article class="comment-card">
                <div class="post-meta">
                    <span class="author-info">
                        <strong>u/${comment.username || 'Anonymous'}</strong>
                    </span>
                    <span>•</span>
                    <span class="post-date">${createdDate}</span>
                </div>

                <div class="post-content">
                    ${comment.content}
                </div>

                <div style="display: flex; justify-content: space-between; align-items: center; margin-top: var(--space-lg); flex-wrap: wrap; gap: var(--space-md);">
                    <div class="reaction-buttons">
                        <button class="reaction-btn like-btn btn-sm" onclick="handleCommentReaction('${comment.comment_id}', 'like')">
                            👍 ${comment.likes || 0}
                        </button>
                        <button class="reaction-btn dislike-btn btn-sm" onclick="handleCommentReaction('${comment.comment_id}', 'dislike')">
                            👎 ${comment.dislikes || 0}
                        </button>
                    </div>

                    ${isAuthor ? `
                        <div class="comment-actions">
                            <button class="delete-btn" onclick="handleDeleteComment('${comment.comment_id}')">
                                🗑️ Delete
                            </button>
                        </div>
                    ` : ''}
                </div>
            </article>
        `;
    },

    setupEventListeners() {
        // Add comment form submission
        const form = document.getElementById('add-comment-form');
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                await this.handleAddComment();
            });
        }

        // Make global functions available for onclick handlers
        window.handlePostReaction = this.handlePostReaction.bind(this);
        window.handleCommentReaction = this.handleCommentReaction.bind(this);
        window.handleDeletePost = this.handleDeletePost.bind(this);
        window.handleDeleteComment = this.handleDeleteComment.bind(this);
        window.showImageLightbox = this.showImageLightbox.bind(this);
    },

    async handleAddComment() {
        const content = document.getElementById('comment-content').value.trim();

        if (!content) {
            this.showToast('Please enter a comment', 'warning');
            return;
        }

        try {
            await apiClient.post(`/comments/create-on-post/${this.postId}`, {
                content: content
            });

            // Clear form
            document.getElementById('comment-content').value = '';

            // Reload comments
            await this.loadComments();
            await this.loadPost(); // Reload to update comment count

            this.showToast('Comment posted successfully!', 'success');
        } catch (error) {
            console.error('[PostView] Error adding comment:', error);
            this.showToast(error.message || 'Failed to post comment', 'error');
        }
    },

    async handlePostReaction(postId, reactionType) {
        try {
            // Convert reaction type string to integer (1 = like, 2 = dislike)
            const reactionTypeInt = reactionType === 'like' ? 1 : 2;

            await apiClient.post(`/reactions/posts/toggle`, {
                post_id: postId,
                reaction_type: reactionTypeInt
            });

            // Reload post to update reaction counts
            await this.loadPost();
            this.showToast(`Reaction recorded!`, 'success');
        } catch (error) {
            console.error('[PostView] Error reacting to post:', error);
            this.showToast(error.message || 'Failed to record reaction', 'error');
        }
    },

    async handleCommentReaction(commentId, reactionType) {
        try {
            // Convert reaction type string to integer (1 = like, 2 = dislike)
            const reactionTypeInt = reactionType === 'like' ? 1 : 2;

            await apiClient.post(`/reactions/comments/toggle`, {
                comment_id: commentId,
                reaction_type: reactionTypeInt
            });

            // Reload comments to update reaction counts
            await this.loadComments();
            this.showToast(`Reaction recorded!`, 'success');
        } catch (error) {
            console.error('[PostView] Error reacting to comment:', error);
            this.showToast(error.message || 'Failed to record reaction', 'error');
        }
    },

    async handleDeletePost(postId) {
        if (!confirm('Are you sure you want to delete this post?')) {
            return;
        }

        try {
            await apiClient.delete(`/posts/remove/${postId}`);
            this.showToast('Post deleted successfully!', 'success');

            // Redirect to home after a short delay
            setTimeout(() => {
                window.router.navigate('/');
            }, 1500);
        } catch (error) {
            console.error('[PostView] Error deleting post:', error);
            this.showToast(error.message || 'Failed to delete post', 'error');
        }
    },

    async handleDeleteComment(commentId) {
        if (!confirm('Are you sure you want to delete this comment?')) {
            return;
        }

        try {
            await apiClient.delete(`/comments/remove/${commentId}`);
            this.showToast('Comment deleted successfully!', 'success');

            // Reload comments and post
            await this.loadComments();
            await this.loadPost(); // Reload to update comment count
        } catch (error) {
            console.error('[PostView] Error deleting comment:', error);
            this.showToast(error.message || 'Failed to delete comment', 'error');
        }
    },

    showImageLightbox(imageUrl) {
        const lightbox = document.getElementById('image-lightbox');
        const img = lightbox.querySelector('img');
        img.src = imageUrl;
        lightbox.classList.add('active');
    },

    showToast(message, type = 'info') {
        // Create toast container if it doesn't exist
        let container = document.querySelector('.toast-container');
        if (!container) {
            container = document.createElement('div');
            container.className = 'toast-container';
            document.body.appendChild(container);
        }

        // Create toast
        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        toast.innerHTML = `
            <div class="toast-content">
                <div class="toast-message">${message}</div>
            </div>
            <button class="toast-close" onclick="this.parentElement.remove()">×</button>
        `;

        container.appendChild(toast);

        // Auto-remove after 3 seconds
        setTimeout(() => {
            toast.classList.add('removing');
            setTimeout(() => toast.remove(), 300);
        }, 3000);
    },

    showError(message) {
        const container = document.getElementById('post-content');
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">⚠️</div>
                <h3 class="empty-state-title">Error</h3>
                <p class="empty-state-message">${message}</p>
                <a href="/" data-link class="btn btn-primary" style="margin-top: 1rem;">Back to Home</a>
            </div>
        `;
    }
};
