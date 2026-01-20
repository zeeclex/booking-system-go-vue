package repository

import (
	"backend/models"

	"github.com/jmoiron/sqlx"
)

type RoomRepository struct {
	DB *sqlx.DB
}

func NewRoomRepository(db *sqlx.DB) *RoomRepository {
	return &RoomRepository{DB: db}
}

// GetAllRooms
func (r *RoomRepository) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	query := "SELECT id, name, type, capacity, is_active FROM rooms"
	err := r.DB.Select(&rooms, query)
	return rooms, err
}

// CreateRoom 
func (r *RoomRepository) CreateRoom(room models.Room) error {
	query := "INSERT INTO rooms (name, type, capacity, is_active) VALUES (?, ?, ?, ?)"
	_, err := r.DB.Exec(query, room.Name, room.Type, room.Capacity, room.IsActive)
	return err
}

// UpdateRoom 
func (r *RoomRepository) UpdateRoom(room models.Room) error {
    query := "UPDATE rooms SET name = ?, type = ?, capacity = ?, is_active = ? WHERE id = ?"
    _, err := r.DB.Exec(query, room.Name, room.Type, room.Capacity, room.IsActive, room.ID)
    return err
}

// HasActiveBookings
func (r *RoomRepository) HasActiveBookings(roomID int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM bookings WHERE room_id = ? AND status IN ('pending', 'approved')"
	err := r.DB.Get(&count, query, roomID)
	return count > 0, err
}

// DeleteRoom 
func (r *RoomRepository) DeleteRoom(id int) error {
	query := "DELETE FROM rooms WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	return err
}
