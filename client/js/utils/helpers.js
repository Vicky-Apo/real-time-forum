// utils/helpers.js - Shared Utility Functions

/**
 * Get user initials from username
 * @param {string} username - The username
 * @returns {string} Two-character uppercase initials or 'U' as fallback
 */
export function getInitials(username) {
    if (!username || typeof username !== 'string') {
        return 'U';
    }
    return username.slice(0, 2).toUpperCase();
}

/**
 * Show image lightbox modal
 * @param {string} imageUrl - The image URL to display
 */
export function showImageLightbox(imageUrl) {
    const modalRoot = document.getElementById('modal-root');
    if (!modalRoot) return;

    modalRoot.innerHTML = `
        <div class="modal-overlay" onclick="this.parentElement.innerHTML = ''">
            <div class="modal-content" onclick="event.stopPropagation()">
                <button class="modal-close" onclick="this.closest('.modal-overlay').parentElement.innerHTML = ''">&times;</button>
                <img src="${imageUrl}" alt="Full size image" style="max-width: 90vw; max-height: 90vh;">
            </div>
        </div>
    `;
}

/**
 * Show toast notification
 * @param {string} message - The message to display
 * @param {string} type - The type of toast ('success', 'error', 'info')
 */
export function showToast(message, type = 'info') {
    const container = document.getElementById('toast-container');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    toast.textContent = message;
    container.appendChild(toast);

    // Auto-remove after 3 seconds
    setTimeout(() => {
        toast.style.opacity = '0';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

/**
 * Get relative time string (e.g., "2 hours ago")
 * @param {string|Date} timestamp - The timestamp to format
 * @returns {string} Relative time string
 */
export function getTimeAgo(timestamp) {
    const now = new Date();
    const past = new Date(timestamp);
    const diffMs = now - past;
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHour = Math.floor(diffMin / 60);
    const diffDay = Math.floor(diffHour / 24);

    if (diffSec < 60) return 'just now';
    if (diffMin < 60) return `${diffMin} minute${diffMin !== 1 ? 's' : ''} ago`;
    if (diffHour < 24) return `${diffHour} hour${diffHour !== 1 ? 's' : ''} ago`;
    if (diffDay < 7) return `${diffDay} day${diffDay !== 1 ? 's' : ''} ago`;

    // For older dates, return formatted date
    return past.toLocaleDateString();
}

/**
 * Format date to locale string
 * @param {string|Date} date - The date to format
 * @returns {string} Formatted date string
 */
export function formatDate(date) {
    return new Date(date).toLocaleDateString();
}

/**
 * Format time to locale string
 * @param {string|Date} date - The date to format
 * @returns {string} Formatted time string (HH:MM)
 */
export function formatTime(date) {
    return new Date(date).toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit'
    });
}
