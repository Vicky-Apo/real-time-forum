// views/ProfileView.js - User Profile Page

import apiClient from '../api/client.js';
import state from '../state.js';
import { getInitials } from '../utils/helpers.js';

export default {
    profile: null,
    userId: null,

    async render(params) {
        this.userId = params.id;

        return `
            <div class="container">
                <div class="profile-container">
                    <div id="profile-loading" class="loading-container">
                        <div class="loading-spinner"></div>
                        <p>Loading profile...</p>
                    </div>

                    <div id="profile-content" style="display: none;">
                        <!-- Profile header -->
                        <div class="profile-header">
                            <div class="profile-avatar-large" id="profile-avatar"></div>
                            <div class="profile-info">
                                <h1 id="profile-username"></h1>
                                <p id="profile-email" class="profile-email"></p>
                                <p id="profile-joined" class="profile-meta"></p>
                            </div>
                        </div>

                        <!-- Profile stats -->
                        <div class="profile-stats">
                            <div class="stats-grid">
                            <div class="stat-card">
                                    <div class="stat-icon"><i class="fas fa-file-alt" style="color: var(--color-primary);"></i></div>
                                    <div class="stat-details">
                                        <div class="stat-number" id="stat-posts">0</div>
                                <div class="stat-label">Posts</div>
                                    </div>
                            </div>
                            <div class="stat-card">
                                    <div class="stat-icon"><i class="fas fa-comments" style="color: var(--color-primary);"></i></div>
                                    <div class="stat-details">
                                        <div class="stat-number" id="stat-comments">0</div>
                                <div class="stat-label">Comments</div>
                                    </div>
                            </div>
                            <div class="stat-card">
                                    <div class="stat-icon"><i class="fas fa-thumbs-up" style="color: var(--color-primary);"></i></div>
                                    <div class="stat-details">
                                        <div class="stat-number" id="stat-likes-received">0</div>
                                <div class="stat-label">Likes Received</div>
                                    </div>
                            </div>
                            <div class="stat-card">
                                    <div class="stat-icon"><i class="fas fa-heart" style="color: var(--color-primary);"></i></div>
                                    <div class="stat-details">
                                        <div class="stat-number" id="stat-posts-liked">0</div>
                                <div class="stat-label">Posts Liked</div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- Activity tabs -->
                        <div class="profile-tabs">
                            <button class="tab-btn active" data-tab="posts">
                                My Posts
                            </button>
                            <button class="tab-btn" data-tab="liked">
                                Liked Posts
                            </button>
                            <button class="tab-btn" data-tab="commented">
                                Commented Posts
                            </button>
                        </div>

                        <!-- Tab content -->
                        <div class="profile-tab-content">
                            <div id="tab-posts" class="tab-pane active">
                                <div class="loading-text">Loading posts...</div>
                            </div>
                            <div id="tab-liked" class="tab-pane">
                                <div class="loading-text">Loading liked posts...</div>
                            </div>
                            <div id="tab-commented" class="tab-pane">
                                <div class="loading-text">Loading commented posts...</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
    },

    async afterRender() {
        await this.loadProfile();
        this.setupEventListeners();
    },

    async loadProfile() {
        const loading = document.getElementById('profile-loading');
        const content = document.getElementById('profile-content');

        try {
            const response = await apiClient.get(`/users/profile/${this.userId}`);
            this.profile = response.data || response;

            // Populate profile data
            this.renderProfileData();

            // Load initial tab (posts)
            await this.loadUserPosts();

            loading.style.display = 'none';
            content.style.display = 'block';

        } catch (error) {
            console.error('[ProfileView] Error loading profile:', error);
            loading.innerHTML = `
                <div class="error-message">
                    <p>Failed to load profile. ${error.message}</p>
                    <button class="btn btn-secondary" onclick="window.history.back()">Go Back</button>
                </div>
            `;
        }
    },

    renderProfileData() {
        // Avatar
        const avatar = document.getElementById('profile-avatar');
        avatar.textContent = getInitials(this.profile.username);

        // Basic info
        document.getElementById('profile-username').textContent = this.profile.username;
        document.getElementById('profile-email').textContent = this.profile.email;

        const joinedDate = new Date(this.profile.created_at).toLocaleDateString('en-US', {
            month: 'long',
            year: 'numeric'
        });
        document.getElementById('profile-joined').textContent = `Joined ${joinedDate}`;

        // Stats
        const stats = this.profile.stats || {};
        document.getElementById('stat-posts').textContent = stats.total_posts || 0;
        document.getElementById('stat-comments').textContent = stats.total_comments || 0;
        document.getElementById('stat-likes-received').textContent = stats.likes_received || 0;
        document.getElementById('stat-posts-liked').textContent = stats.posts_liked || 0;
    },

    setupEventListeners() {
        // Tab switching
        const tabBtns = document.querySelectorAll('.tab-btn');
        tabBtns.forEach(btn => {
            btn.addEventListener('click', async () => {
                const tabName = btn.dataset.tab;

                // Update active states
                tabBtns.forEach(b => b.classList.remove('active'));
                btn.classList.add('active');

                document.querySelectorAll('.tab-pane').forEach(pane => {
                    pane.classList.remove('active');
                });
                document.getElementById(`tab-${tabName}`).classList.add('active');

                // Load tab content
                if (tabName === 'posts') {
                    await this.loadUserPosts();
                } else if (tabName === 'liked') {
                    await this.loadLikedPosts();
                } else if (tabName === 'commented') {
                    await this.loadCommentedPosts();
                }
            });
        });
    },

    async loadUserPosts() {
        const container = document.getElementById('tab-posts');
        try {
            const response = await apiClient.get(`/users/posts/${this.userId}`);
            const data = response.data || response;
            const posts = data.posts || [];

            if (posts.length === 0) {
                container.innerHTML = `
                    <div class="empty-state">
                        <p>No posts yet</p>
                    </div>
                `;
                return;
            }

            container.innerHTML = posts.map(post => this.renderPostCard(post)).join('');
        } catch (error) {
            console.error('[ProfileView] Error loading posts:', error);
            container.innerHTML = `<div class="error-message">Failed to load posts</div>`;
        }
    },

    async loadLikedPosts() {
        const container = document.getElementById('tab-liked');
        try {
            const response = await apiClient.get(`/users/liked-posts/${this.userId}`);
            const data = response.data || response;
            const posts = data.posts || [];

            if (posts.length === 0) {
                container.innerHTML = `
                    <div class="empty-state">
                        <p>No liked posts yet</p>
                    </div>
                `;
                return;
            }

            container.innerHTML = posts.map(post => this.renderPostCard(post)).join('');
        } catch (error) {
            console.error('[ProfileView] Error loading liked posts:', error);
            container.innerHTML = `<div class="error-message">Failed to load liked posts</div>`;
        }
    },

    async loadCommentedPosts() {
        const container = document.getElementById('tab-commented');
        try {
            const response = await apiClient.get(`/users/commented-posts/${this.userId}`);
            const data = response.data || response;
            const posts = data.posts || [];

            if (posts.length === 0) {
                container.innerHTML = `
                    <div class="empty-state">
                        <p>No commented posts yet</p>
                    </div>
                `;
                return;
            }

            container.innerHTML = posts.map(post => this.renderPostCard(post)).join('');
        } catch (error) {
            console.error('[ProfileView] Error loading commented posts:', error);
            container.innerHTML = `<div class="error-message">Failed to load commented posts</div>`;
        }
    },

    renderPostCard(post) {
        const date = new Date(post.created_at).toLocaleDateString();
        const categories = (post.categories || []).map(cat =>
            `<a href="/category/${cat.category_id}" data-link class="category-tag">${cat.category_name || cat.name}</a>`
        ).join('');

        return `
            <article class="profile-post-card" onclick="window.router.navigate('/post/${post.id}')">
                <div class="profile-post-header">
                    <div class="profile-post-meta">
                        <span class="profile-post-author">${post.username || 'Anonymous'}</span>
                        <span class="meta-separator">•</span>
                        <span class="profile-post-date">${date}</span>
                    </div>
                </div>
                <div class="profile-post-content">
                    ${post.content.substring(0, 200)}${post.content.length > 200 ? '...' : ''}
                </div>
                ${categories ? `
                    <div class="profile-post-categories">
                        ${categories}
                    </div>
                ` : ''}
                <div class="profile-post-stats">
                    <span class="stat-item">
                        <i class="fas fa-thumbs-up"></i>
                        <span>${post.like_count || 0}</span>
                    </span>
                    <span class="stat-item">
                        <i class="fas fa-thumbs-down"></i>
                        <span>${post.dislike_count || 0}</span>
                    </span>
                    <span class="stat-item">
                        <i class="fas fa-comments"></i>
                        <span>${post.comment_count || 0}</span>
                    </span>
                </div>
            </article>
        `;
    }
};
