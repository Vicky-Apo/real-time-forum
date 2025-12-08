// views/CreatePostView.js - Create Post Page

import apiClient from '../api/client.js';
import { navigate } from '../router.js';

export default {
    async render() {
        return `
            <div class="container">
                <div class="create-post-container">
                    <div class="create-post-header">
                        <h1>Create Post</h1>
                        <p>Share your thoughts with the community</p>
                    </div>

                    <form id="create-post-form" class="create-post-form">
                        <div class="form-group">
                            <label for="post-content">What's on your mind?</label>
                            <textarea
                                id="post-content"
                                name="content"
                                placeholder="Share something interesting..."
                                required
                                rows="8"
                                maxlength="5000"
                            ></textarea>
                            <div class="char-count">
                                <span id="char-counter">0</span> / 5000
                            </div>
                        </div>

                        <div class="form-group">
                            <label for="categories">Categories (optional)</label>
                            <div id="categories-container" class="categories-checkboxes">
                                <div class="loading-text">Loading categories...</div>
                            </div>
                        </div>

                        <div class="form-group">
                            <label for="image-upload">Add Image (optional)</label>
                            <input
                                type="file"
                                id="image-upload"
                                name="image"
                                accept="image/*"
                            >
                            <small class="form-hint">Supported formats: JPG, PNG, GIF (max 5MB)</small>
                        </div>

                        <div id="error-message" class="error-message"></div>

                        <div class="form-actions">
                            <button type="button" class="btn btn-secondary" id="cancel-btn">
                                Cancel
                            </button>
                            <button type="submit" class="btn btn-primary">
                                Post
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        `;
    },

    async afterRender() {
        console.log('[CreatePostView] Rendered');

        await this.loadCategories();
        this.setupEventListeners();
    },

    async loadCategories() {
        const container = document.getElementById('categories-container');

        try {
            const response = await apiClient.get('/categories');
            const data = response.data || response;
            const categories = Array.isArray(data) ? data : (data.categories || []);

            if (categories.length === 0) {
                container.innerHTML = '<p class="form-hint">No categories available</p>';
                return;
            }

            container.innerHTML = categories.map(cat => `
                <label class="category-checkbox">
                    <input type="checkbox" name="categories[]" value="${cat.category_name || cat.name}">
                    <span>${cat.category_name || cat.name}</span>
                </label>
            `).join('');

        } catch (error) {
            console.error('[CreatePostView] Error loading categories:', error);
            container.innerHTML = '<p class="form-hint text-danger">Failed to load categories</p>';
        }
    },

    setupEventListeners() {
        const form = document.getElementById('create-post-form');
        const textarea = document.getElementById('post-content');
        const charCounter = document.getElementById('char-counter');
        const cancelBtn = document.getElementById('cancel-btn');

        // Character counter
        textarea.addEventListener('input', () => {
            charCounter.textContent = textarea.value.length;
        });

        // Cancel button
        cancelBtn.addEventListener('click', () => {
            navigate('/');
        });

        // Form submission
        form.addEventListener('submit', this.handleSubmit.bind(this));
    },

    async handleSubmit(e) {
        e.preventDefault();

        const form = e.target;
        const content = form.querySelector('#post-content').value.trim();

        // Get selected categories (as category names)
        const selectedCategories = Array.from(
            form.querySelectorAll('input[name="categories[]"]:checked')
        ).map(checkbox => checkbox.value);

        // Validation
        if (!content) {
            this.showError('Please write something!');
            return;
        }

        if (content.length < 10) {
            this.showError('Post must be at least 10 characters long');
            return;
        }

        if (selectedCategories.length === 0) {
            this.showError('Please select at least one category');
            return;
        }

        // Disable submit button
        const submitBtn = form.querySelector('button[type="submit"]');
        const originalText = submitBtn.textContent;
        submitBtn.disabled = true;
        submitBtn.textContent = 'Posting...';

        try {
            // Create FormData for multipart/form-data
            const formData = new FormData();
            formData.append('content', content);

            // Append each category as a separate 'categories' field
            selectedCategories.forEach(category => {
                formData.append('categories', category);
            });

            // Handle image upload if present
            const imageFile = form.querySelector('#image-upload').files[0];
            if (imageFile) {
                formData.append('images', imageFile);
            }

            // Create post
            const response = await apiClient.post('/posts/create', formData);

            console.log('[CreatePostView] Post created:', response);

            // Redirect to home
            navigate('/');

        } catch (error) {
            console.error('[CreatePostView] Error creating post:', error);
            this.showError(error.message || 'Failed to create post. Please try again.');

            // Re-enable button
            submitBtn.disabled = false;
            submitBtn.textContent = originalText;
        }
    },

    showError(message) {
        const errorDiv = document.getElementById('error-message');
        errorDiv.textContent = message;
        errorDiv.classList.add('show');

        // Hide after 5 seconds
        setTimeout(() => {
            errorDiv.classList.remove('show');
        }, 5000);
    }
};
