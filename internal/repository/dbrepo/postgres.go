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

// CheckDuplicateRule checks if a duplicate rule already exists in the database
func (p *postgresDBRepo) CheckDuplicateRule(createrule models.CreateRule) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1 FROM rules
            WHERE source_ip = $1 AND destination_ip = $2 AND port = $3
        )
    `

	var exists bool
	err := p.DB.QueryRow(query, createrule.SourceIP, createrule.DestinationIP, createrule.Port).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
