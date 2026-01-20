package services

import (
	"testing"
	"time"
)

// Helper for parsing time strings
func parseTime(tStr string) time.Time {
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, tStr)
	return t
}

// UNIT TEST 1: function IsTimeConflict
func TestIsTimeConflict(t *testing.T) {
	existStart := parseTime("2026-02-01 10:00:00")
	existEnd := parseTime("2026-02-01 12:00:00")
	tests := []struct {
		name     string
		newStart string
		newEnd   string
		want     bool
	}{
		{
			name:     "Aman - Jadwal Pagi (08:00 - 10:00)",
			newStart: "2026-02-01 08:00:00",
			newEnd:   "2026-02-01 10:00:00",
			want:     false,
		},
		{
			name:     "Bentrok - Overlap Awal (09:00 - 11:00)",
			newStart: "2026-02-01 09:00:00",
			newEnd:   "2026-02-01 11:00:00",
			want:     true,
		},
		{
			name:     "Bentrok - Overlap Total (10:30 - 11:30)",
			newStart: "2026-02-01 10:30:00",
			newEnd:   "2026-02-01 11:30:00",
			want:     true,
		},
		{
			name:     "Bentrok - Overlap Akhir (11:00 - 13:00)",
			newStart: "2026-02-01 11:00:00",
			newEnd:   "2026-02-01 13:00:00",
			want:     true,
		},
		{
			name:     "Aman - Jadwal Siang (12:00 - 14:00)",
			newStart: "2026-02-01 12:00:00",
			newEnd:   "2026-02-01 14:00:00",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newS := parseTime(tt.newStart)
			newE := parseTime(tt.newEnd)
			
			got := IsTimeConflict(newS, newE, existStart, existEnd)
			
			if got != tt.want {
				t.Errorf("IsTimeConflict() = %v, want %v", got, tt.want)
			}
		})
	}
}

// UNIT TEST 2: function CalculateDurationHours
func TestCalculateDurationHours(t *testing.T) {
	tests := []struct {
		name      string
		startTime string
		endTime   string
		want      float64
		wantErr   bool
	}{
		{
			name:      "Hitung 2 Jam",
			startTime: "2026-02-01 10:00:00",
			endTime:   "2026-02-01 12:00:00",
			want:      2.0,
			wantErr:   false,
		},
		{
			name:      "Hitung 1.5 Jam (90 Menit)",
			startTime: "2026-02-01 10:00:00",
			endTime:   "2026-02-01 11:30:00",
			want:      1.5,
			wantErr:   false,
		},
		{
			name:      "Error - Waktu Terbalik",
			startTime: "2026-02-01 12:00:00",
			endTime:   "2026-02-01 10:00:00",
			want:      0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateDurationHours(tt.startTime, tt.endTime)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateDurationHours() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateDurationHours() = %v, want %v", got, tt.want)
			}
		})
	}
}