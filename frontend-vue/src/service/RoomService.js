import ApiService from './ApiService';

export const RoomService = {
    getRooms() {
        return ApiService.get('/rooms');
    },
    createRoom(data) {
        const payload = { 
            ...data, 
            capacity: parseInt(data.capacity) 
        };
        return ApiService.post('/rooms', payload);
    },
    updateRoom(id, data) {
        const { id: _, ...updateData } = data; 
        const payload = { 
            ...updateData, 
            capacity: parseInt(data.capacity) 
        };
        return ApiService.put(`/rooms/${id}`, payload);
    },
    deleteRoom(id) {
        return ApiService.delete(`/rooms/${id}`);
    }
};