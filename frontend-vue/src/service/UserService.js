import ApiService from './ApiService';

export const UserService = {
    // Get all users
    getUsers() {
        return ApiService.get('/users');
    },

    // Create new user
    createUser(data) {
        return ApiService.post('/users', data);
    },

    // Update existing user
    updateUser(id, data) {
        return ApiService.put(`/users/${id}`, data);
    },

    // Delete user
    deleteUser(id) {
        return ApiService.delete(`/users/${id}`);
    }
};