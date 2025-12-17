const hostname = window.location.hostname;
const isDevelopment = hostname === 'localhost' || hostname === '127.0.0.1';

export const config = {
    // Environment
    environment: isDevelopment ? 'development' : 'production',
    isDevelopment,
    
    // API
    apiBaseURL: '/api',
    
    // WebSocket
    wsProtocol: window.location.protocol === 'https:' ? 'wss:' : 'ws:',
    wsHost: window.location.host,
    
    get wsUrl() {
        return `${this.wsProtocol}//${this.wsHost}/ws`;
    },
    
    // Features
    features: {
        browserNotifications: 'Notification' in window,
        webSocket: typeof WebSocket !== 'undefined',
        localStorage: typeof localStorage !== 'undefined'
    }
};

export default config;