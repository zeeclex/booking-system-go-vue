package handlers

import (
	"backend/repository"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	BookingRepo *repository.BookingRepository
	RoomRepo    *repository.RoomRepository
}

// Helper: Smart Time Parser
func parseFlexibleTime(timeStr string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return t, nil
	}
	
	if t, err := time.Parse("2006-01-02 15:04:05", timeStr); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("failed to parse time: %s", timeStr)
}

// GenerateMonthlyReport generates monthly usage statistics
// @Summary      Generate Monthly Report
// @Description  Create CSV & JSON report containing usage stats by Month and Year
// @Tags         Reports
// @Produce      json
// @Param        month query int true "Month (1-12)"
// @Param        year  query int true "Year (e.g. 2026)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /reports/generate [get]
func (h *ReportHandler) GenerateMonthlyReport(c *gin.Context) {
	// 1. Validate Input
	monthStr := c.Query("month")
	yearStr := c.Query("year")
	month, errMonth := strconv.Atoi(monthStr)
	year, errYear := strconv.Atoi(yearStr)
	if errMonth != nil || errYear != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Month and Year parameters must be numbers"})
		return
	}

	// 2. Fetch Booking Data
	bookings, err := h.BookingRepo.GetBookingsByMonth(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings: " + err.Error()})
		return
	}

	// 3. Prepare Storage Folder
	storagePath := filepath.Join("storage", "reports")
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage directory"})
		return
	}
	fileNameBase := fmt.Sprintf("report_%d-%02d", year, month)
	jsonFilePath := filepath.Join(storagePath, fileNameBase+".json")
	csvFilePath := filepath.Join(storagePath, fileNameBase+".csv")

	// 4. Write JSON File (Raw Data Backup)
	jsonData, _ := json.MarshalIndent(bookings, "", "  ")
	if err := os.WriteFile(jsonFilePath, jsonData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save JSON file"})
		return
	}

	// 5. Write CSV File (Summary Report)
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CSV file"})
		return
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write CSV Header
	header := []string{"Room ID", "Room Name", "Type", "Total Hours Used", "Total Bookings"}
	if err := writer.Write(header); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header"})
		return
	}

	// --- STATS CALCULATION LOGIC ---
	summaryHours := make(map[int]float64)
	summaryCount := make(map[int]int)
	for _, b := range bookings {
		start, errStart := parseFlexibleTime(b.StartTime)
		end, errEnd := parseFlexibleTime(b.EndTime)
		if errStart != nil || errEnd != nil {
			log.Printf("[Warning] Skipping Booking ID %d due to time parse error: %v", b.ID, errStart)
			continue 
		}
		duration := end.Sub(start).Hours()
		if duration > 0 {
			summaryHours[b.RoomID] += duration
			summaryCount[b.RoomID]++
		}
	}

	// Fetch All Rooms
	rooms, err := h.RoomRepo.GetAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch room data"})
		return
	}

	// Loop Rooms and Write to CSV
	for _, r := range rooms {
		totalHours := 0.0
		totalCount := 0
		if val, ok := summaryHours[r.ID]; ok {
			totalHours = val
		}
		if val, ok := summaryCount[r.ID]; ok {
			totalCount = val
		}
		row := []string{
			strconv.Itoa(r.ID),
			r.Name,
			r.Type,
			fmt.Sprintf("%.2f", totalHours),
			strconv.Itoa(totalCount),
		}
		
		writer.Write(row)
	}

	// 6. Response with Download URLs
	// [Note] Replace localhost with your actual domain when deploying
	baseURL := "http://localhost:8080" 
	
	c.JSON(http.StatusOK, gin.H{
		"message":           "Report generated successfully",
		"csv_file_name":     fileNameBase + ".csv",
		"json_file_name":    fileNameBase + ".json",
		"download_url_csv":  fmt.Sprintf("%s/storage/reports/%s.csv", baseURL, fileNameBase),
		"download_url_json": fmt.Sprintf("%s/storage/reports/%s.json", baseURL, fileNameBase),
	})
}