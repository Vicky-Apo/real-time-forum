import CONSTANTS from './constants.js';
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
    }, CONSTANTS.TOAST_DURATION_MS);
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

/**
 * Throttle function - Limits how often a function can be called
 * Ensures the function is called at most once per specified delay
 * Perfect for scroll events to prevent performance issues
 * @param {Function} func - The function to throttle
 * @param {number} delay - Minimum time between function calls in milliseconds
 * @returns {Function} Throttled function
 */
export function throttle(func, delay = 300) {
    let lastCall = 0;
    let timeoutId = null;

    return function (...args) {
        const now = Date.now();
        const timeSinceLastCall = now - lastCall;

        // Clear any pending timeout
        if (timeoutId) {
            clearTimeout(timeoutId);
        }

        if (timeSinceLastCall >= delay) {
            // Enough time has passed, execute immediately
            lastCall = now;
            func.apply(this, args);
        } else {
            // Schedule execution after remaining delay
            timeoutId = setTimeout(() => {
                lastCall = Date.now();
                func.apply(this, args);
            }, delay - timeSinceLastCall);
        }
    };
}

/**
 * Debounce function - Delays function execution until after specified time has elapsed
 * since the last time it was invoked
 * @param {Function} func - The function to debounce
 * @param {number} delay - Time to wait before calling function in milliseconds
 * @returns {Function} Debounced function
 */
export function debounce(func, delay = 300) {
    let timeoutId = null;

    return function (...args) {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => {
            func.apply(this, args);
        }, delay);
    };
}
