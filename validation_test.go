package tcmrsv

import (
	"testing"
	"time"
)

func TestIsIDValid(t *testing.T) {
	tests := []struct {
		id    string
		valid bool
	}{
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"123E4567-E89B-12D3-A456-426614174000", false},
		{"in-va-li-d-id", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsIDValid(tt.id); got != tt.valid {
			t.Errorf("IsIDValid(%q) = %v; want %v", tt.id, got, tt.valid)
		}
	}
}

func TestIsDateWithin2Days(t *testing.T) {
	now := time.Now().In(jst())
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, jst())
	tomorrow := today.AddDate(0, 0, 1)
	dayAfterTomorrow := today.AddDate(0, 0, 2)
	threeDaysLater := today.AddDate(0, 0, 3)
	yesterday := today.AddDate(0, 0, -1)

	tests := []struct {
		date  time.Time
		valid bool
	}{
		{today, true},
		{tomorrow, true},
		{dayAfterTomorrow, true},
		{threeDaysLater, false},
		{yesterday, false},
	}

	for _, tt := range tests {
		if got := IsDateWithin2Days(tt.date); got != tt.valid {
			t.Errorf("IsDateWithin2Days(%v) = %v; want %v", tt.date, got, tt.valid)
		}
	}
}

func TestIsTimeRangeValid(t *testing.T) {
	tests := []struct {
		fromH, fromM, toH, toM int
		valid                  bool
	}{
		{7, 0, 8, 0, true},
		{7, 30, 8, 0, true},
		{22, 30, 23, 0, true},
		{22, 30, 22, 0, false},
		{6, 0, 8, 0, false},
		{7, 15, 8, 0, false},
		{7, 0, 8, 45, false},
		{22, 0, 23, 30, false},
		{20, 0, 19, 30, false},
	}

	for _, tt := range tests {
		if got := IsTimeRangeValid(tt.fromH, tt.fromM, tt.toH, tt.toM); got != tt.valid {
			t.Errorf("IsTimeRangeValid(%d:%02d to %d:%02d) = %v; want %v",
				tt.fromH, tt.fromM, tt.toH, tt.toM, got, tt.valid)
		}
	}
}

func TestIsTimeInFuture(t *testing.T) {
	now := time.Now().In(jst())
	h := now.Hour()
	m := now.Minute()
	currentTotal := h*60 + m

	tests := []struct {
		fromH, fromM int
		want         bool
	}{
		{h, m + 1, true},
		{h + 1, 0, true},
		{h, m, true},
		{h, m - 1, false},
		{h - 1, 59, false},
	}

	for _, tt := range tests {
		got := IsTimeInFuture(tt.fromH, tt.fromM)
		if got != tt.want {
			t.Errorf("IsTimeInFuture(%02d:%02d) = %v; want %v (current: %02d:%02d, total=%d)",
				tt.fromH, tt.fromM, got, tt.want, h, m, currentTotal)
		}
	}
}

func TestIsCommentValid(t *testing.T) {
	tests := []struct {
		comment string
		valid   bool
	}{
		{"hello", true},
		{"  world  ", true},
		{"", false},
		{"     ", false},
	}

	for _, tt := range tests {
		if got := IsCommentValid(tt.comment); got != tt.valid {
			t.Errorf("IsCommentValid(%q) = %v; want %v", tt.comment, got, tt.valid)
		}
	}
}
