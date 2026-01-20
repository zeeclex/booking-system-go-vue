package handlers

import (
	"backend/models"
	"backend/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	Repo *repository.RoomRepository
}

// GetRooms menangani permintaan untuk melihat semua daftar ruang
// @Summary Daftar semua ruang
// @Description Mengambil semua data ruang kelas dan laboratorium yang terdaftar.
// @Tags Rooms
// @Produce json
// @Success 200 {array} models.Room
// @Failure 500 {object} map[string]string "error: Gagal mengambil data"
// @Router /rooms [get]
func (h *RoomHandler) GetRooms(c *gin.Context) {
	rooms, err := h.Repo.GetAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ruang"})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

// CreateRoom menangani penambahan ruang baru (khusus Admin)
// @Summary Tambah ruang baru
// @Description Menambahkan data ruang baru ke dalam sistem (Khusus Admin).
// @Tags Rooms
// @Accept json
// @Produce json
// @Param room body models.Room true "Data Ruang Baru"
// @Success 201 {object} map[string]string "message: Ruang berhasil ditambahkan"
// @Failure 400 {object} map[string]string "error: Input data tidak valid"
// @Failure 500 {object} map[string]string "error: Gagal menyimpan data"
// @Router /rooms [post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input data tidak valid"})
		return
	}
	if err := h.Repo.CreateRoom(room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data ruang"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Ruang berhasil ditambahkan"})
}

// UpdateRoom menangani perubahan data ruang
// @Summary Perbarui data ruang
// @Description Mengubah informasi ruang (nama, kapasitas, tipe) berdasarkan ID.
// @Tags Rooms
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Param room body models.Room true "Data Ruang yang diupdate"
// @Success 200 {object} map[string]string "message: Data ruang berhasil diperbarui"
// @Failure 400 {object} map[string]string "error: ID atau data tidak valid"
// @Failure 500 {object} map[string]string "error: Gagal memperbarui data"
// @Router /rooms/{id} [put]
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID ruang tidak valid"})
		return
	}
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input data tidak valid"})
		return
	}
	room.ID = id
	if err := h.Repo.UpdateRoom(room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data ruang"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data ruang berhasil diperbarui"})
}

// DeleteRoom menangani penghapusan ruang dengan syarat tidak ada booking aktif
// @Summary Hapus ruang
// @Description Menghapus data ruang jika tidak ada pemesanan aktif (pending/approved) yang terikat.
// @Tags Rooms
// @Param id path int true "Room ID"
// @Success 200 {object} map[string]string "message: Ruang berhasil dihapus"
// @Failure 400 {object} map[string]string "error: ID tidak valid"
// @Failure 409 {object} map[string]string "error: Ruang memiliki pemesanan aktif"
// @Failure 500 {object} map[string]string "error: Gagal menghapus ruang"
// @Router /rooms/{id} [delete]
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID ruang tidak valid"})
		return
	}

	// 1. Validation
	hasActive, err := h.Repo.HasActiveBookings(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi status pemesanan ruang"})
		return
	}
	if hasActive {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Ruang tidak dapat dihapus karena masih memiliki pemesanan aktif",
		})
		return
	}

	// 2. excecute delete
	if err := h.Repo.DeleteRoom(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus ruang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ruang berhasil dihapus"})
}