package models

type User struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Role  string `db:"role" json:"role"`
	Email string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type Room struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Type     string `db:"type" json:"type"`
	Capacity int    `db:"capacity" json:"capacity"`
	IsActive bool   `db:"is_active" json:"is_active"`
}

type Booking struct {
    ID        int    `db:"id" json:"id"`
    RoomID    int    `db:"room_id" json:"room_id"`
    UserID    int    `db:"user_id" json:"user_id"`
    StartTime string `db:"start_time" json:"start_time"`
    EndTime   string `db:"end_time" json:"end_time"`
    Purpose   string `db:"purpose" json:"purpose"`
    Status    string `db:"status" json:"status"`
}

type BookingResponse struct {
    Booking
    RoomName string `db:"room_name" json:"room_name"`
    UserName string `db:"user_name" json:"user_name"`
}

type ReportSummary struct {
	RoomName   string  `json:"room_name"`
	TotalHours float64 `json:"total_hours"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}