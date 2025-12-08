// views/HomeView.js - Home/Feed Page

import apiClient from '../api/client.js';

export default {
    async render() {
        return `
            <div class="container">
                <div class="home-layout">
                    <!-- Main Content (Posts Feed) -->
                    <div class="main-content">
                        <div id="posts-container" class="posts-container">
                            <div class="loading-container">
                                <div class="loading-spinner"></div>
                                <p>Loading posts...</p>
                            </div>
                        </div>
                    </div>

                    <!-- Sidebar (Categories) -->
                    <aside class="sidebar">
                        <div class="sidebar-card">
                            <h3 class="sidebar-title">Categories</h3>
                            <div id="categories-list" class="categories-list">
                                <div class="loading-text">Loading categories...</div>
                            </div>
                        </div>
                    </aside>
                </div>
            </div>
        `;
    },

    async afterRender() {
        console.log('[HomeView] Rendered');

        // Show hero section on home page
        const heroSection = document.getElementById('hero-section');
        if (heroSection) {
            heroSection.style.display = 'block';
        }

        await Promise.all([
            this.loadPosts(),
            this.loadCategories()
        ]);
    },

    async loadCategories() {
        const container = document.getElementById('categories-list');

        try {
            const response = await apiClient.get('/categories');
            const data = response.data || response;
            const categories = Array.isArray(data) ? data : (data.categories || []);

            if (categories.length === 0) {
                container.innerHTML = '<p class="empty-text">No categories available</p>';
                return;
            }

            container.innerHTML = categories.map(cat => `
                <a href="/category/${cat.category_id}" data-link class="category-item">
                    <span class="category-name">${cat.category_name || cat.name}</span>
                    <span class="category-count">${cat.post_count || 0}</span>
                </a>
            `).join('');

        } catch (error) {
            console.error('[HomeView] Error loading categories:', error);
            container.innerHTML = '<p class="empty-text error">Failed to load categories</p>';
        }
    },

    async loadPosts() {
        const container = document.getElementById('posts-container');

        try {
            const response = await apiClient.get('/posts?limit=20&offset=0');
            const data = response.data || response;
            const posts = data.posts || [];

            if (posts.length === 0) {
                container.innerHTML = `
                    <div class="empty-state">
                        <div class="empty-state-icon">📭</div>
                        <h3 class="empty-state-title">No posts yet</h3>
                        <p class="empty-state-message">Be the first to start a conversation!</p>
                        <a href="/create-post" data-link class="btn btn-primary" style="margin-top: 1rem;">Create Post</a>
                    </div>
                `;
                return;
            }

            container.innerHTML = posts.map(post => this.renderPostCard(post)).join('');

        } catch (error) {
            console.error('[HomeView] Error loading posts:', error);
            container.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">⚠️</div>
                    <h3 class="empty-state-title">Failed to load posts</h3>
                    <p class="empty-state-message">Please try again later.</p>
                </div>
            `;
        }
    },

    renderPostCard(post) {
        const createdDate = new Date(post.created_at).toLocaleDateString();
        const categories = post.categories?.map(cat =>
            `<span class="category-tag">${cat.category_name || cat.name}</span>`
        ).join(' ') || '';

        const netVotes = (post.likes || 0) - (post.dislikes || 0);

        return `
            <article class="post-card" onclick="window.router.navigate('/post/${post.post_id}')">
                <div class="post-vote">
                    <button class="vote-btn upvote" onclick="event.stopPropagation()">▲</button>
                    <span class="vote-count">${netVotes}</span>
                    <button class="vote-btn downvote" onclick="event.stopPropagation()">▼</button>
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

};
