// views/CategoryView.js - Category Posts View

import apiClient from '../api/client.js';

export default {
    categoryId: null,
    category: null,

    init() {
        // Setup global handler for vote buttons
        window.handleCategoryVote = (postId, reactionType) => {
            this.handleVote(postId, reactionType);
        };
    },

    async render(params) {
        this.init();
        this.categoryId = params.id;

        return `
            <div class="container">
                <div id="category-content">
                    <div class="loading-container">
                        <div class="loading-spinner"></div>
                        <p>Loading posts...</p>
                    </div>
                </div>
            </div>
        `;
    },

    async afterRender() {
        console.log('[CategoryView] Rendered for category:', this.categoryId);

        try {
            await this.loadCategoryPosts();
        } catch (error) {
            console.error('[CategoryView] Error loading category:', error);
            this.showError('Failed to load category posts');
        }
    },

    async loadCategoryPosts() {
        const container = document.getElementById('category-content');

        try {
            const response = await apiClient.get(`/posts/by-category/${this.categoryId}`);
            const data = response.data || response;
            const posts = data.posts || [];
            this.category = data.category || { category_name: 'Category' };

            container.innerHTML = `
                <!-- Back Button -->
                <div style="margin-bottom: var(--space-2xl);">
                    <a href="/" data-link class="btn btn-secondary">← Back to Feed</a>
                </div>

                <!-- Category Header -->
                <div class="page-banner category-page">
                    <div class="banner-content">
                        <h2>📁 ${this.category.category_name || 'Category'}</h2>
                        <p class="banner-description">Browse posts in this category</p>
                    </div>
                    <div class="post-count">${posts.length} posts</div>
                </div>

                <!-- Posts List -->
                <div id="posts-container" class="posts-container">
                    ${posts.length > 0 ? posts.map(post => this.renderPostCard(post)).join('') : this.renderEmptyState()}
                </div>
            `;

        } catch (error) {
            console.error('[CategoryView] Error loading category posts:', error);
            container.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">⚠️</div>
                    <h3 class="empty-state-title">Failed to load category posts</h3>
                    <p class="empty-state-message">${error.message || 'Please try again later.'}</p>
                    <a href="/" data-link class="btn btn-primary" style="margin-top: 1rem;">Back to Home</a>
                </div>
            `;
            throw error;
        }
    },

    renderPostCard(post) {
        const createdDate = new Date(post.created_at).toLocaleDateString();
        const categories = post.categories?.map(cat =>
            `<span class="category-tag">${cat.category_name || cat.name}</span>`
        ).join(' ') || '';

        const netVotes = (post.likes || 0) - (post.dislikes || 0);

        return `
            <article class="post-card" data-post-id="${post.post_id}" onclick="window.router.navigate('/post/${post.post_id}')">
                <div class="post-vote">
                    <button class="vote-btn upvote" onclick="event.stopPropagation(); window.handleCategoryVote('${post.post_id}', 1)" title="Like this post">▲</button>
                    <span class="vote-count" title="Net votes (likes - dislikes)">${netVotes}</span>
                    <button class="vote-btn downvote" onclick="event.stopPropagation(); window.handleCategoryVote('${post.post_id}', 2)" title="Dislike this post">▼</button>
                </div>
                <div class="post-main">
                    <div class="post-meta">
                        <span><strong>u/${post.username || 'Anonymous'}</strong></span>
                        <span>•</span>
                        <span>${createdDate}</span>
                        ${categories ? `<span>•</span>${categories}` : ''}
                    </div>
                    <div class="post-content">
                        ${post.content}
                    </div>
                    <div class="post-stats">
                        <span>💬 ${post.comment_count || 0} comments</span>
                    </div>
                </div>
            </article>
        `;
    },

    renderEmptyState() {
        return `
            <div class="empty-state">
                <div class="empty-state-icon">📭</div>
                <h3 class="empty-state-title">No posts in this category yet</h3>
                <p class="empty-state-message">Be the first to post in this category!</p>
                <a href="/create-post" data-link class="btn btn-primary" style="margin-top: 1rem;">Create Post</a>
            </div>
        `;
    },

    showError(message) {
        const container = document.getElementById('category-content');
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">⚠️</div>
                <h3 class="empty-state-title">Error</h3>
                <p class="empty-state-message">${message}</p>
                <a href="/" data-link class="btn btn-primary" style="margin-top: 1rem;">Back to Home</a>
            </div>
        `;
    },

    async handleVote(postId, reactionType) {
        try {
            // Make API call to toggle reaction
            await apiClient.post('/reactions/posts/toggle', {
                post_id: postId,
                reaction_type: parseInt(reactionType)
            });

            // Reload just this post's data to update vote counts
            const response = await apiClient.get(`/posts/view/${postId}`);
            const updatedPost = response.data || response;

            // Find the post card and update vote count
            const postCard = document.querySelector(`[data-post-id="${postId}"]`);
            if (postCard) {
                const netVotes = (updatedPost.likes || 0) - (updatedPost.dislikes || 0);
                const voteCountElement = postCard.querySelector('.vote-count');
                if (voteCountElement) {
                    voteCountElement.textContent = netVotes;
                }
            }

            console.log('[CategoryView] Vote recorded successfully');
        } catch (error) {
            console.error('[CategoryView] Error voting on post:', error);
            alert(error.message || 'Failed to record vote');
        }
    }
};
