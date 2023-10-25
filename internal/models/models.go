package models

// Createrule holds ip/port data
type CreateRule struct {
	SourceIP      string
	DestinationIP string
	Port          string
}
