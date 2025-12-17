import config from '../config.js';

const isDev = config.isDevelopment;

export default {
    log: (...args) => {
        if (isDev) {
            console.log(...args);
        }
    },
    
    warn: (...args) => {
        if (isDev) {
            console.warn(...args);
        }
    },
    
    error: (...args) => {
        // Always log errors, even in production
        console.error(...args);
    },
    
    info: (...args) => {
        if (isDev) {
            console.info(...args);
        }
    }
};