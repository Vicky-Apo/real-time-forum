/**
 * HTTP Client for Backend API
 */

import config from '../config.js';

class APIClient {
    constructor(baseURL = '') {
        // Use relative URLs (will work with Nginx proxy)
        this.baseURL = baseURL || config.apiBaseURL;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;

        const config = {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            credentials: 'include', // Include cookies for session
        };

        // Remove Content-Type for FormData (browser sets it with boundary)
        if (options.body instanceof FormData) {
            delete config.headers['Content-Type'];
        }

        try {
            const response = await fetch(url, config);

            // Try to parse JSON response
            let data;
            const contentType = response.headers.get('content-type');

            if (contentType && contentType.includes('application/json')) {
                data = await response.json();
            } else {
                data = await response.text();
            }

            if (!response.ok) {
                const error = new Error(data.error || data || 'Request failed');
                error.status = response.status;
                error.data = data;
                throw error;
            }


            // Normalize data structure from backend to frontend format
            return this.normalizeResponse(data);

        } catch (error) {
            console.error(`[API] ${config.method || 'GET'} ${url} - Error:`, error.message);
            throw error;
        }
    }

    // Normalize backend response to frontend format
    normalizeResponse(data) {
        // If response has success and data fields, extract data
        if (data && typeof data === 'object' && 'success' in data && 'data' in data) {
            data = data.data;
        }

        // Normalize posts
        if (Array.isArray(data)) {
            return data.map(item => this.normalizeItem(item));
        } else if (data && typeof data === 'object') {
            // Check if it's a paginated response with posts/comments
            if (data.posts && Array.isArray(data.posts)) {
                data.posts = data.posts.map(post => this.normalizePost(post));
            }
            if (data.comments && Array.isArray(data.comments)) {
                data.comments = data.comments.map(comment => this.normalizeComment(comment));
            }
            // Single item response
            if (data.post_id || data.post_content) {
                return this.normalizePost(data);
            }
            if (data.comment_id || data.comment_content) {
                return this.normalizeComment(data);
            }
        }

        return data;
    }

    normalizeItem(item) {
        if (!item || typeof item !== 'object') return item;

        if (item.post_id || item.post_content) {
            return this.normalizePost(item);
        }
        if (item.comment_id || item.comment_content) {
            return this.normalizeComment(item);
        }
        return item;
    }

    normalizePost(post) {
        if (!post) return post;

        return {
            ...post,
            id: post.post_id || post.id,
            content: post.post_content || post.content,
            // Normalize images array
            images: post.images || [],
        };
    }

    normalizeComment(comment) {
        if (!comment) return comment;

        return {
            ...comment,
            id: comment.comment_id || comment.id,
            content: comment.comment_content || comment.content,
        };
    }

    // GET request
    get(endpoint, options = {}) {
        return this.request(endpoint, {
            method: 'GET',
            ...options,
        });
    }

    // POST request
    post(endpoint, body, options = {}) {
        return this.request(endpoint, {
            method: 'POST',
            body: body instanceof FormData ? body : JSON.stringify(body),
            ...options,
        });
    }

    // PUT request
    put(endpoint, body, options = {}) {
        return this.request(endpoint, {
            method: 'PUT',
            body: body instanceof FormData ? body : JSON.stringify(body),
            ...options,
        });
    }

    // DELETE request
    delete(endpoint, options = {}) {
        return this.request(endpoint, {
            method: 'DELETE',
            ...options,
        });
    }
}

// Create and export singleton instance
const apiClient = new APIClient();
export default apiClient;
