package repository

import (
    "backend/models"
    "github.com/jmoiron/sqlx"
)

type BookingRepository struct {
    DB *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) *BookingRepository {
    return &BookingRepository{DB: db}
}

// CheckAvailability
// Logic: (StartA < EndB) AND (EndA > StartB)
func (r *BookingRepository) IsRoomAvailable(roomID int, startTime string, endTime string) (bool, error) {
    var count int
    query := `
        SELECT COUNT(*) 
        FROM bookings 
        WHERE room_id = ? 
			AND status IN ('approved', 'pending') 
			AND start_time < ? 
			AND end_time > ?
    `
    err := r.DB.Get(&count, query, roomID, endTime, startTime)
    if err != nil {
        return false, err
    }
    
    return count == 0, nil
}

// CreateBooking
func (r *BookingRepository) CreateBooking(b models.Booking) error {
    query := `INSERT INTO bookings (room_id, user_id, start_time, end_time, purpose, status) 
            VALUES (?, ?, ?, ?, ?, 'pending')`
    _, err := r.DB.Exec(query, b.RoomID, b.UserID, b.StartTime, b.EndTime, b.Purpose)
    return err
}

// GetAllBookings (without date)
func (r *BookingRepository) GetAllBookings() ([]models.Booking, error) {
    var bookings []models.Booking
    query := "SELECT * FROM bookings ORDER BY start_time DESC"
    err := r.DB.Select(&bookings, query)
    return bookings, err
}

// GetBookingsByUserID
func (r *BookingRepository) GetBookingsByUserID(userID int) ([]models.Booking, error) {
    var bookings []models.Booking
    query := "SELECT * FROM bookings WHERE user_id = ? ORDER BY start_time DESC"
    err := r.DB.Select(&bookings, query, userID)
    return bookings, err
}

// UpdateStatus
func (r *BookingRepository) UpdateStatus(id int, status string) error {
    query := "UPDATE bookings SET status = ? WHERE id = ?"
    _, err := r.DB.Exec(query, status, id)
    return err
}

// GetDashboardStats (Update query count without date)
func (r *BookingRepository) GetDashboardStats() (map[string]interface{}, error) {
    var totalRooms, activeRooms, totalBookings, pendingBookings int
    r.DB.Get(&totalRooms, "SELECT COUNT(*) FROM rooms")
    r.DB.Get(&activeRooms, "SELECT COUNT(DISTINCT room_id) FROM bookings WHERE DATE(start_time) = CURDATE() AND status = 'approved'")
    r.DB.Get(&totalBookings, "SELECT COUNT(*) FROM bookings")
    r.DB.Get(&pendingBookings, "SELECT COUNT(*) FROM bookings WHERE status = 'pending'")
    return map[string]interface{}{
        "total_rooms":      totalRooms,
        "active_usage":     activeRooms,
        "total_bookings":   totalBookings,
        "pending_requests": pendingBookings,
    }, nil
}

// GetBookingsByMonth 
func (r *BookingRepository) GetBookingsByMonth(month int, year int) ([]models.Booking, error) {
    var bookings []models.Booking
    
    query := `
        SELECT id, room_id, user_id, start_time, end_time, purpose, status 
        FROM bookings 
        WHERE MONTH(start_time) = ? AND YEAR(start_time) = ? AND status = 'approved'
        ORDER BY start_time ASC
    `
    err := r.DB.Select(&bookings, query, month, year)
    return bookings, err
}