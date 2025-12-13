// views/EditPostView.js - Edit Post Page

import apiClient from '../api/client.js';
import { navigate } from '../router.js';

export default {
    post: null,

    async render(params) {
        this.postId = params.id;

        return `
            <div class="container">
                <div class="create-post-container">
                    <div class="create-post-header">
                        <h1>Edit Post</h1>
                        <p>Update your post</p>
                    </div>

                    <div id="loading-state" class="loading-container">
                        <div class="loading-spinner"></div>
                        <p>Loading post...</p>
                    </div>

                    <form id="edit-post-form" class="create-post-form" style="display: none;">
                        <div class="form-group">
                            <label for="post-content">Content</label>
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
                            <label>Current Images</label>
                            <div id="current-images"></div>
                        </div>

                        <div class="form-group">
                            <label for="image-upload">Add New Image (optional)</label>
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
                                Update Post
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        `;
    },

    async afterRender() {

        await this.loadPost();
        await this.loadCategories();
        this.setupEventListeners();
    },

    async loadPost() {
        const loadingState = document.getElementById('loading-state');
        const form = document.getElementById('edit-post-form');

        try {
            const response = await apiClient.get(`/posts/view/${this.postId}`);
            this.post = response.data || response;

            // Populate form
            const contentTextarea = document.getElementById('post-content');
            contentTextarea.value = this.post.content || '';

            // Update character counter
            const charCounter = document.getElementById('char-counter');
            charCounter.textContent = contentTextarea.value.length;

            // Show current images
            const currentImagesDiv = document.getElementById('current-images');
            if (this.post.images && this.post.images.length > 0) {
                currentImagesDiv.innerHTML = this.post.images.map(img => `
                    <div class="current-image">
                        <img src="${img.image_url}" alt="Post image" style="max-width: 200px; border-radius: 8px;">
                    </div>
                `).join('');
            } else {
                currentImagesDiv.innerHTML = '<p class="form-hint">No images</p>';
            }

            loadingState.style.display = 'none';
            form.style.display = 'block';

        } catch (error) {
            console.error('[EditPostView] Error loading post:', error);
            loadingState.innerHTML = `
                <div class="error-message">
                    <p>Failed to load post. ${error.message}</p>
                    <button class="btn btn-secondary" onclick="window.history.back()">Go Back</button>
                </div>
            `;
        }
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

            // Get post's current categories
            const postCategories = (this.post?.categories || []).map(c => c.category_name || c.name);

            container.innerHTML = categories.map(cat => {
                const catName = cat.category_name || cat.name;
                const isChecked = postCategories.includes(catName) ? 'checked' : '';
                return `
                    <label class="category-checkbox">
                        <input type="checkbox" name="categories[]" value="${catName}" ${isChecked}>
                        <span>${catName}</span>
                    </label>
                `;
            }).join('');

        } catch (error) {
            console.error('[EditPostView] Error loading categories:', error);
            container.innerHTML = '<p class="form-hint text-danger">Failed to load categories</p>';
        }
    },

    setupEventListeners() {
        const form = document.getElementById('edit-post-form');
        const contentInput = document.getElementById('post-content');
        const charCounter = document.getElementById('char-counter');
        const cancelBtn = document.getElementById('cancel-btn');

        // Character counter
        contentInput?.addEventListener('input', () => {
            charCounter.textContent = contentInput.value.length;
        });

        // Cancel button
        cancelBtn?.addEventListener('click', () => {
            navigate(`/post/${this.postId}`);
        });

        // Form submission
        form?.addEventListener('submit', async (e) => {
            e.preventDefault();
            await this.handleSubmit(e);
        });
    },

    async handleSubmit(e) {
        const errorDiv = document.getElementById('error-message');
        const submitBtn = e.target.querySelector('button[type="submit"]');

        try {
            errorDiv.textContent = '';
            submitBtn.disabled = true;
            submitBtn.textContent = 'Updating...';

            const formData = new FormData(e.target);

            // Get selected categories
            const selectedCategories = Array.from(
                e.target.querySelectorAll('input[name="categories[]"]:checked')
            ).map(cb => cb.value);

            // Create request body
            const requestData = new FormData();
            requestData.append('content', formData.get('content'));

            if (selectedCategories.length > 0) {
                selectedCategories.forEach(cat => requestData.append('categories', cat));
            }

            const imageFile = formData.get('image');
            if (imageFile && imageFile.size > 0) {
                requestData.append('image', imageFile);
            }

            // Submit
            await apiClient.put(`/posts/edit/${this.postId}`, requestData);

            // Redirect to post
            navigate(`/post/${this.postId}`);

        } catch (error) {
            console.error('[EditPostView] Error updating post:', error);
            errorDiv.textContent = error.message || 'Failed to update post. Please try again.';
            submitBtn.disabled = false;
            submitBtn.textContent = 'Update Post';
        }
    }
};
