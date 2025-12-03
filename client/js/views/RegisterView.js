// views/RegisterView.js - Registration Page

import apiClient from '../api/client.js';
import state from '../state.js';
import { navigate } from '../router.js';
import wsManager from '../websocket/WebSocketManager.js';

export default {
    async render() {
        return `
            <div class="auth-container">
                <div class="auth-card">
                    <h1>Create Account</h1>
                    <p class="auth-subtitle">Join Real-Time Forum today</p>

                    <form id="register-form" class="auth-form">
                        <div class="form-row">
                            <div class="form-group">
                                <label for="first-name">First Name *</label>
                                <input
                                    type="text"
                                    id="first-name"
                                    name="first_name"
                                    placeholder="John"
                                    required
                                    autofocus
                                    minlength="1"
                                    maxlength="50"
                                >
                            </div>
                            <div class="form-group">
                                <label for="last-name">Last Name *</label>
                                <input
                                    type="text"
                                    id="last-name"
                                    name="last_name"
                                    placeholder="Doe"
                                    required
                                    minlength="1"
                                    maxlength="50"
                                >
                            </div>
                        </div>

                        <div class="form-group">
                            <label for="username">Username *</label>
                            <input
                                type="text"
                                id="username"
                                name="username"
                                placeholder="johndoe"
                                required
                                minlength="3"
                                maxlength="20"
                                pattern="[a-zA-Z0-9_]+"
                                title="Username can only contain letters, numbers, and underscores"
                            >
                            <small class="form-hint">3-20 characters, letters, numbers, and underscores only</small>
                        </div>

                        <div class="form-group">
                            <label for="email">Email *</label>
                            <input
                                type="email"
                                id="email"
                                name="email"
                                placeholder="john@example.com"
                                required
                            >
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label for="age">Age *</label>
                                <input
                                    type="number"
                                    id="age"
                                    name="age"
                                    placeholder="18"
                                    required
                                    min="13"
                                    max="120"
                                >
                                <small class="form-hint">Must be 13 or older</small>
                            </div>
                            <div class="form-group">
                                <label for="gender">Gender *</label>
                                <select id="gender" name="gender" required>
                                    <option value="">Select...</option>
                                    <option value="Male">Male</option>
                                    <option value="Female">Female</option>
                                    <option value="Other">Other</option>
                                </select>
                            </div>
                        </div>

                        <div class="form-group">
                            <label for="password">Password *</label>
                            <input
                                type="password"
                                id="password"
                                name="password"
                                placeholder="Enter a strong password"
                                required
                                minlength="8"
                            >
                            <small class="form-hint">At least 8 characters</small>
                        </div>

                        <div class="form-group">
                            <label for="confirm-password">Confirm Password *</label>
                            <input
                                type="password"
                                id="confirm-password"
                                name="confirm_password"
                                placeholder="Re-enter your password"
                                required
                            >
                        </div>

                        <div id="error-message" class="error-message"></div>

                        <button type="submit" class="btn btn-primary btn-block">
                            Create Account
                        </button>
                    </form>

                    <p class="auth-footer">
                        Already have an account?
                        <a href="/login" data-link class="auth-link">Sign in</a>
                    </p>
                </div>
            </div>
        `;
    },

    afterRender() {
        console.log('[RegisterView] Rendered');

        const form = document.getElementById('register-form');
        form.addEventListener('submit', this.handleSubmit.bind(this));

        // Real-time password match validation
        const password = document.getElementById('password');
        const confirmPassword = document.getElementById('confirm-password');

        confirmPassword.addEventListener('input', () => {
            if (confirmPassword.value && password.value !== confirmPassword.value) {
                confirmPassword.setCustomValidity('Passwords do not match');
            } else {
                confirmPassword.setCustomValidity('');
            }
        });

        password.addEventListener('input', () => {
            if (confirmPassword.value) {
                if (password.value !== confirmPassword.value) {
                    confirmPassword.setCustomValidity('Passwords do not match');
                } else {
                    confirmPassword.setCustomValidity('');
                }
            }
        });
    },

    async handleSubmit(e) {
        e.preventDefault();

        const form = e.target;
        const formData = new FormData(form);

        // Collect form data
        const data = {
            first_name: formData.get('first_name').trim(),
            last_name: formData.get('last_name').trim(),
            username: formData.get('username').trim(),
            email: formData.get('email').trim(),
            age: parseInt(formData.get('age')),
            gender: formData.get('gender'),
            password: formData.get('password')
        };

        const confirmPassword = formData.get('confirm_password');

        // Validate passwords match
        if (data.password !== confirmPassword) {
            this.showError('Passwords do not match');
            return;
        }

        // Validate age
        if (data.age < 13 || data.age > 120) {
            this.showError('Age must be between 13 and 120');
            return;
        }

        // Disable submit button
        const submitBtn = form.querySelector('button[type="submit"]');
        const originalText = submitBtn.textContent;
        submitBtn.disabled = true;
        submitBtn.textContent = 'Creating account...';

        try {
            // Call register API
            const response = await apiClient.post('/auth/register', data);

            console.log('[RegisterView] Registration successful:', response);

            // Auto-login: call login API
            const loginData = {
                identifier: data.username,
                password: data.password
            };

            const loginResponse = await apiClient.post('/auth/login', loginData);

            // Store user in state
            if (loginResponse.user) {
                state.setUser(loginResponse.user);
            }

            // Connect WebSocket
            wsManager.connect();

            // Redirect to home
            navigate('/');

        } catch (error) {
            console.error('[RegisterView] Registration failed:', error);
            this.showError(error.message || 'Registration failed. Please try again.');

            // Re-enable button
            submitBtn.disabled = false;
            submitBtn.textContent = originalText;
        }
    },

    showError(message) {
        const errorDiv = document.getElementById('error-message');
        errorDiv.textContent = message;
        errorDiv.classList.add('show');

        // Scroll to error
        errorDiv.scrollIntoView({ behavior: 'smooth', block: 'nearest' });

        // Hide after 5 seconds
        setTimeout(() => {
            errorDiv.classList.remove('show');
        }, 5000);
    }
};
