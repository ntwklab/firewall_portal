package models

import "time"

// Create rule holds ip/port data
type CreateRule struct {
	SourceIP      string
	DestinationIP string
	Port          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
