// views/NotificationsView.js - Notifications Page

export default {
    async render() {
        return `
            <div class="container">
                <h1>Notifications</h1>
                <p>Your notifications will appear here!</p>
            </div>
        `;
    },

    afterRender() {
        console.log('[NotificationsView] Rendered');
    }
};
