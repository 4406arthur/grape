package postgres

import (
	"context"
	"database/sql"
)

type SummaryRepository struct {
	DB *sql.DB
}

func NewSummaryRepository(db *sql.DB) *SummaryRepository {
	return &SummaryRepository{
		DB: db,
	}
}

func (db *SummaryRepository) GetMaxBlockNum(ctx context.Context) int64 {
	var max int64
	db.DB.QueryRowContext(ctx, "SELECT max FROM max_block").Scan(&max)
	return max
}
