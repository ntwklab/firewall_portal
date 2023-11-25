package dbrepo

import (
	"context"
	"time"

	"github.com/ntwklab/firewall_portal/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertRule inserts a rule into the database
func (m *postgresDBRepo) InsertRule(rule models.CreateRule) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into rules (source_ip, destination_ip, port, created_at, updated_at)
			values ($1, $2, $3, $4, $5)`

	_, err := m.DB.ExecContext(ctx, stmt,
		rule.SourceIP,
		rule.DestinationIP,
		rule.Port,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
