// views/RegisterView.js - Registration Page

export default {
    async render() {
        return `
            <div class="auth-container">
                <div class="auth-card">
                    <h1>Create Account</h1>
                    <p>Register view - Coming soon!</p>
                    <a href="/login" data-link>Already have an account? Login</a>
                </div>
            </div>
        `;
    },

    afterRender() {
        console.log('[RegisterView] Rendered');
    }
};
