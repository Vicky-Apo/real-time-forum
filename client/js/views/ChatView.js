// views/ChatView.js - Chat/Messaging Page

export default {
    async render() {
        return `
            <div class="container">
                <h1>Chat</h1>
                <p>Real-time messaging - Coming soon!</p>
            </div>
        `;
    },

    afterRender() {
        console.log('[ChatView] Rendered');
    }
};
