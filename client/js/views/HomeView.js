// views/HomeView.js - Home/Feed Page

export default {
    async render() {
        return `
            <div class="container">
                <h1>Home Feed</h1>
                <p>Posts will appear here!</p>
            </div>
        `;
    },

    afterRender() {
        console.log('[HomeView] Rendered');
    }
};
