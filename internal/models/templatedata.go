package models

import "github.com/ntwklab/firewall_portal/internal/forms"

// TemplateData holds data sent from handlers to temnplate
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
