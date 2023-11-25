package repository

import "github.com/ntwklab/firewall_portal/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertRule(rule models.CreateRule) error
}
