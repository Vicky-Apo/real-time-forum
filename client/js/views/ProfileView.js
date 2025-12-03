// views/ProfileView.js - User Profile Page

export default {
    async render(params) {
        return `
            <div class="container">
                <h1>User Profile</h1>
                <p>Viewing profile ID: ${params.id}</p>
            </div>
        `;
    },

    afterRender() {
        console.log('[ProfileView] Rendered');
    }
};
