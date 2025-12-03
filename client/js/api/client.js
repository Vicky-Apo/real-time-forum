// api/client.js - HTTP Client for Backend API

class APIClient {
    constructor(baseURL = '') {
        // Use relative URLs (will work with Nginx proxy)
        this.baseURL = baseURL || '/api';
        console.log('[API] Client initialized with baseURL:', this.baseURL);
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

        console.log(`[API] ${config.method || 'GET'} ${url}`);

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

            console.log(`[API] ${config.method || 'GET'} ${url} - Success`);
            return data;

        } catch (error) {
            console.error(`[API] ${config.method || 'GET'} ${url} - Error:`, error.message);
            throw error;
        }
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
