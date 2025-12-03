// views/PostView.js - Single Post View

export default {
    async render(params) {
        return `
            <div class="container">
                <h1>Post View</h1>
                <p>Viewing post ID: ${params.id}</p>
            </div>
        `;
    },

    afterRender() {
        console.log('[PostView] Rendered');
    }
};
