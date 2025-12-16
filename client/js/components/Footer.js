// Footer.js - Footer Component

import state from '../state.js';

export function renderFooter() {
    const footer = document.getElementById('footer');
    if (!footer) return;

    const user = state.getUser();
    const profileLink = user ? `/profile/${user.id}` : '/profile';

    footer.innerHTML = `
        <div class="footer-content">
            <div class="footer-text">
                © 2025 Real-Time Forum. All rights reserved.
            </div>
            <div class="footer-links">
                <a href="/" class="footer-link">
                    <i class="fas fa-home"></i> Home
                </a>
                <a href="${profileLink}" class="footer-link">
                    <i class="fas fa-user"></i> Profile
                </a>
                <a href="/chat" class="footer-link">
                    <i class="fas fa-comments"></i> Chat
                </a>
            </div>
        </div>
    `;
}

