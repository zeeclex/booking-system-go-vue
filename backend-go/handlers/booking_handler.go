package handlers

import (
	"backend/models"
	"backend/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	Repo *repository.BookingRepository
}

// CreateBooking menangani pengajuan jadwal baru dengan validasi bentrok
// @Summary Ajukan pemesanan ruang
// @Description Dosen mengajukan jadwal penggunaan ruang. Sistem akan mengecek bentrok jadwal secara otomatis.
// @Tags Bookings
// @Accept json
// @Produce json
// @Param booking body models.Booking true "Data Pemesanan Ruang"
// @Success 201 {object} map[string]string "message: Pemesanan berhasil diajukan"
// @Failure 400 {object} map[string]string "error: Data input tidak valid"
// @Failure 409 {object} map[string]string "error: Jadwal bentrok"
// @Failure 500 {object} map[string]string "error: Gagal menyimpan data"
// @Router /bookings [post]
// CreateBooking
func (h *BookingHandler) CreateBooking(c *gin.Context) {
    var b models.Booking
    if err := c.ShouldBindJSON(&b); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
        return
    }

    // 1. Validasi Backend
    isAvailable, err := h.Repo.IsRoomAvailable(b.RoomID, b.StartTime, b.EndTime)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error checking availability"})
        return
    }

    if !isAvailable {
        c.JSON(http.StatusConflict, gin.H{
            "error": "Room is already booked for the selected time slot.",
        })
        return
    }

    if err := h.Repo.CreateBooking(b); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Booking submitted successfully"})
}

// GetBookings mengambil daftar pemesanan dengan filter opsional user_id
// @Summary List daftar pemesanan
// @Description Mengambil semua data pemesanan (Admin) atau filter berdasarkan User ID (Dosen).
// @Tags Bookings
// @Produce json
// @Param user_id query int false "Filter berdasarkan User ID (Dosen)"
// @Success 200 {array} models.Booking
// @Failure 400 {object} map[string]string "error: Parameter tidak valid"
// @Failure 500 {object} map[string]string "error: Gagal mengambil data"
// @Router /bookings [get]
func (h *BookingHandler) GetBookings(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter user_id harus berupa angka"})
			return
		}
		bookings, err := h.Repo.GetBookingsByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, bookings)
		return
	}
	bookings, err := h.Repo.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// ApproveBooking mengubah status pemesanan menjadi 'approved'
// @Summary Setujui pemesanan (Admin)
// @Description Mengubah status pemesanan menjadi 'approved' berdasarkan ID.
// @Tags Bookings
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string "message: Pemesanan telah disetujui"
// @Failure 400 {object} map[string]string "error: ID tidak valid"
// @Failure 500 {object} map[string]string "error: Gagal update data"
// @Router /bookings/{id}/approve [post]
func (h *BookingHandler) ApproveBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pemesanan tidak valid"})
		return
	}
	if err := h.Repo.UpdateStatus(id, "approved"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyetujui pemesanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pemesanan telah disetujui (Approved)"})
}

// RejectBooking mengubah status pemesanan menjadi 'rejected'
// @Summary Tolak pemesanan (Admin)
// @Description Mengubah status pemesanan menjadi 'rejected' berdasarkan ID.
// @Tags Bookings
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string "message: Pemesanan telah ditolak"
// @Failure 400 {object} map[string]string "error: ID tidak valid"
// @Failure 500 {object} map[string]string "error: Gagal update data"
// @Router /bookings/{id}/reject [post]
func (h *BookingHandler) RejectBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pemesanan tidak valid"})
		return
	}
	if err := h.Repo.UpdateStatus(id, "rejected"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menolak pemesanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pemesanan telah ditolak (Rejected)"})
}

// @Summary Ambil Statistik Dashboard
// @Tags Stats
// @Router /stats [get]
func (h *BookingHandler) GetStats(c *gin.Context) {
	stats, err := h.Repo.GetDashboardStats()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, stats)
}