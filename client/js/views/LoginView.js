// views/LoginView.js - Login Page

export default {
    async render() {
        return `
            <div class="auth-container">
                <div class="auth-card">
                    <h1>Welcome Back</h1>
                    <p>Login view - Coming soon!</p>
                    <a href="/register" data-link>Don't have an account? Register</a>
                </div>
            </div>
        `;
    },

    afterRender() {
        console.log('[LoginView] Rendered');
    }
};
