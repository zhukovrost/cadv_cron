package service

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type Service struct {
	db  *sql.DB
	log *zap.Logger
}

func New(db *sql.DB, log *zap.Logger) *Service {
	return &Service{
		db:  db,
		log: log,
	}
}

func (s *Service) ClearTokens(ctx context.Context) (int64, error) {
	q := `
	DELETE FROM sessions
	WHERE expires_at<=NOW()
	OR revoked=true
	`

	result, err := s.db.Exec(q)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
