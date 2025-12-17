export function escapeHtml(text) {
    if (!text) return '';
    
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    
    return String(text).replace(/[&<>"']/g, m => map[m]);
}

export function sanitizeUsername(username) {
    return escapeHtml(username);
}

export function sanitizeContent(content) {
    return escapeHtml(content);
}