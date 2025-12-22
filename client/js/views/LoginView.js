// views/LoginView.js - Login Page

import apiClient from '../api/client.js';
import state from '../state.js';
import { navigate } from '../router.js';
import wsManager from '../websocket/WebSocketManager.js';

export default {
    async render() {
        return `
            <div class="auth-container">
                <div class="auth-card">
                    <h1>Welcome Back</h1>
                    <p class="auth-subtitle">Sign in to continue to Real-Time Forum</p>

                    <form id="login-form" class="auth-form">
                        <div class="form-group">
                            <label for="identifier">Username or Email</label>
                            <input
                                type="text"
                                id="identifier"
                                name="identifier"
                                placeholder="Enter your username or email"
                                required
                                autofocus
                            >
                        </div>

                        <div class="form-group">
                            <label for="password">Password</label>
                            <input
                                type="password"
                                id="password"
                                name="password"
                                placeholder="Enter your password"
                                required
                            >
                        </div>

                        <div id="error-message" class="error-message"></div>

                        <button type="submit" class="btn btn-primary btn-block">
                            Login
                        </button>
                    </form>

                    <p class="auth-footer">
                        Don't have an account?
                        <a href="/register" data-link class="auth-link">Create one</a>
                    </p>
                </div>
            </div>
        `;
    },

    afterRender() {

        const form = document.getElementById('login-form');
        form.addEventListener('submit', this.handleSubmit.bind(this));
    },

    async handleSubmit(e) {
        e.preventDefault();

        const form = e.target;
        const formData = new FormData(form);
        const data = {
            identifier: formData.get('identifier').trim(),
            password: formData.get('password')
        };

        // Validation
        if (!data.identifier || !data.password) {
            this.showError('Please fill in all fields');
            return;
        }

        // Disable submit button
        const submitBtn = form.querySelector('button[type="submit"]');
        const originalText = submitBtn.textContent;
        submitBtn.disabled = true;
        submitBtn.textContent = 'Logging in...';

        try {
            // Call login API
            const response = await apiClient.post('/auth/login', data);

            // Backend wraps response in { success: true, data: {...} }
            const loginData = response.data || response;

            // Store user in state
            if (loginData.user) {
                state.setUser(loginData.user);
            }

            // Connect WebSocket
            wsManager.connect();

            // Redirect to home
            navigate('/');

        } catch (error) {
            console.error('[LoginView] Login failed:', error);
            this.showError(error.message || 'Login failed. Please check your credentials.');

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
