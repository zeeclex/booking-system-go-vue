import ApiService from './ApiService';

export const BookingService = {
    getStats() {
        return ApiService.get('/stats');
    },
    getBookings(userId = null) {
        const url = userId ? `/bookings?user_id=${userId}` : '/bookings';
        return ApiService.get(url);
    },
    createBooking(data) {
        return ApiService.post('/bookings', data);
    },
    updateStatus(id, action) {
        return ApiService.post(`/bookings/${id}/${action}`);
    }
};