// views/CreatePostView.js - Create Post Page

export default {
    async render() {
        return `
            <div class="container">
                <h1>Create New Post</h1>
                <p>Post creation form - Coming soon!</p>
            </div>
        `;
    },

    afterRender() {
        console.log('[CreatePostView] Rendered');
    }
};
