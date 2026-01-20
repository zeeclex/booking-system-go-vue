import ApiService from './ApiService';

export const AuthService = {
    login(credentials) {
        return ApiService.post('/login', credentials);
    },
    logout() {
        localStorage.clear();
        window.location.href = '/login';
    }
};