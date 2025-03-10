package postgresql

import (
	"context"
	"database/sql"
	"fit-tracker/database"
	"os"

	_ "github.com/lib/pq"
)

type (
	postgresql struct {
		DB *sql.DB
	}
)

func New() (*postgresql, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_DSN"))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresql{
		DB: db,
	}, nil
}

func (p *postgresql) SaveTraces(ctx context.Context, input *database.SaveTracesInput) error {
	return nil
}
