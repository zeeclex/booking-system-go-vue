import axios from 'axios';

// 1. BASE URL CONFIGURATION
const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const ApiService = axios.create({
    baseURL: baseURL,
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    },
    timeout: 10000 
});

// --- REQUEST INTERCEPTOR ---
ApiService.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// --- RESPONSE INTERCEPTOR ---
ApiService.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && error.response.status === 401) {
            console.warn('Session expired or invalid token. Logging out...');
            localStorage.clear();
            window.location.href = '/login';
        }
        console.error('API Error:', error.response?.data?.error || error.message);
        return Promise.reject(error);
    }
);

export default ApiService;