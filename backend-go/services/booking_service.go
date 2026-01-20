package services

import (
	"backend/models"
	"backend/repository"
	"errors"
	"time"
)

type BookingService struct {
	Repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{Repo: repo}
}

// 1. BUSINESS LOGIC FUNCTIONS (PURE FUNCTIONS)

// IsTimeConflict 
func IsTimeConflict(newStart, newEnd, existingStart, existingEnd time.Time) bool {
	return newStart.Before(existingEnd) && newEnd.After(existingStart)
}

// CalculateDurationHours 
func CalculateDurationHours(startStr, endStr string) (float64, error) {
	layout := "2006-01-02 15:04:05" // Format MySQL DateTime
	startTime, err := time.Parse(layout, startStr)
	if err != nil {
		return 0, err
	}
	endTime, err := time.Parse(layout, endStr)
	if err != nil {
		return 0, err
	}
	if endTime.Before(startTime) {
		return 0, errors.New("waktu selesai tidak boleh lebih awal dari waktu mulai")
	}
	duration := endTime.Sub(startTime)
	return duration.Hours(), nil
}

// 2. SERVICE METHODS (INTEGRATION WITH REPO)
func (s *BookingService) CreateBooking(req models.Booking) error {
	// 1. Validation
	duration, err := CalculateDurationHours(req.StartTime, req.EndTime)
	if err != nil {
		return err
	}
	if duration <= 0 {
		return errors.New("durasi peminjaman tidak valid")
	}

	// 2. check Availability
	available, err := s.Repo.IsRoomAvailable(req.RoomID, req.StartTime, req.EndTime)
	if err != nil {
		return err
	}
	if !available {
		return errors.New("conflict: room is already booked for this time")
	}

	// 3. save to DB
	return s.Repo.CreateBooking(req)
}

func (s *BookingService) GetBookings(userIDStr string) ([]models.Booking, error) {
	if userIDStr == "" {
		return s.Repo.GetAllBookings()
	}
	return s.Repo.GetAllBookings() 
}

// Wrapper for UpdateStatus
func (s *BookingService) UpdateStatus(id int, status string) error {
	return s.Repo.UpdateStatus(id, status)
}