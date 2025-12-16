// router.js - Client-Side Router using History API

import state from './state.js';

class Router {
    constructor(routes) {
        this.routes = routes;
        this.currentRoute = null;

        // Handle browser back/forward buttons
        window.addEventListener('popstate', () => {
            this.handleRoute();
        });

        // Intercept all link clicks with data-link attribute
        document.addEventListener('click', (e) => {
            const link = e.target.closest('[data-link]');
            if (link) {
                e.preventDefault();
                const href = link.getAttribute('href');
                this.navigate(href);
            }
        });
    }

    // Navigate to a new path
    navigate(path) {
        window.history.pushState(null, null, path);
        this.handleRoute();
    }

    // Handle the current route
    async handleRoute() {
        const path = window.location.pathname;

        const route = this.matchRoute(path);

        if (!route) {
            console.warn('[Router] No route found for:', path);
            // Redirect to home
            this.navigate('/');
            return;
        }

        // Check authentication
        if (route.requiresAuth && !state.getUser()) {
            this.navigate('/login');
            return;
        }

        // If logged in and trying to access login/register, redirect to home
        if (!route.requiresAuth && state.getUser() && (path === '/login' || path === '/register')) {
            this.navigate('/');
            return;
        }

        // Render the view
        this.currentRoute = route;
        await this.renderView(route);
    }

    // Match the current path to a route
    matchRoute(path) {
        for (const route of this.routes) {
            const regex = this.pathToRegex(route.path);
            const match = regex.exec(path);

            if (match) {
                // Extract route parameters (e.g., /post/:id)
                const params = this.extractParams(route.path, match);
                return { ...route, params };
            }
        }
        return null;
    }

    // Convert path pattern to regex
    pathToRegex(path) {
        // Convert /post/:id to /post/([^/]+)
        const pattern = path.replace(/:\w+/g, '([^/]+)');
        return new RegExp(`^${pattern}$`);
    }

    // Extract parameters from matched route
    extractParams(path, match) {
        const keys = [...path.matchAll(/:(\w+)/g)].map(m => m[1]);
        const values = match.slice(1); // Remove first element (full match)

        return keys.reduce((params, key, index) => {
            params[key] = values[index];
            return params;
        }, {});
    }

    // Render the view component
    async renderView(route) {
        const app = document.getElementById('app');

        // Hide hero section by default (individual views can show it)
        const heroSection = document.getElementById('hero-section');
        if (heroSection) {
            heroSection.style.display = 'none';
        }

        // Show loading state
        app.innerHTML = `
            <div class="loading-container">
                <div class="loading-spinner"></div>
                <p>Loading...</p>
            </div>
        `;

        try {
            // Import the view component
            const view = await route.component();

            // Render the view
            const html = await view.render(route.params);
            app.innerHTML = html;

            // Call afterRender if it exists (for event listeners, etc.)
            if (view.afterRender) {
                view.afterRender(route.params);
            }

            // Scroll to top
            window.scrollTo(0, 0);

        } catch (error) {
            console.error('[Router] Error rendering view:', error);
            app.innerHTML = `
                <div class="error-container">
                    <h2>Error Loading Page</h2>
                    <p>${error.message}</p>
                    <button onclick="window.location.reload()" class="btn btn-primary">
                        Reload
                    </button>
                </div>
            `;
        }
    }

    // Get current route
    getCurrentRoute() {
        return this.currentRoute;
    }
}

export default Router;

// Helper function for navigation (can be used in views)
export function navigate(path) {
    if (window.router) {
        window.router.navigate(path);
    } else {
        console.error('[Router] Router not initialized');
    }
}
