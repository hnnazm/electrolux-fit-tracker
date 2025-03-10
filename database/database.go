package database

import (
	"context"
	"time"
)

type (
	Config func(*Database) *Database

	SaveTracesInput struct {
		Steps     int64     `bson:"steps"`
		HeartBeat int64     `bson:"heartBeat"`
		MET       float64   `bson:"met"`
		UserID    string    `bson:"userId"`
		CreatedAt time.Time `bson:"createdAt"`
	}

	GetTracesInput struct {
		UserID    string    `bson:"userId"`
		CreatedAt time.Time `bson:"createdAt"`
	}

	GetTracesResult struct {
		Steps     int64     `bson:"steps"`
		HeartBeat int64     `bson:"heartBeat"`
		MET       float64   `bson:"met"`
		UserID    string    `bson:"userId"`
		CreatedAt time.Time `bson:"createdAt"`
	}

	IngestorRepository interface {
		SaveTraces(ctx context.Context, input *SaveTracesInput) error
		GetTraces(ctx context.Context, input *GetTracesInput) ([]GetTracesResult, error)
	}

	Database struct {
		IngestorRepository IngestorRepository
	}
)

func New(configs ...Config) *Database {
	db := &Database{}

	for _, config := range configs {
		config(db)
	}

	return db
}

func WithIngestorRepository(ingestorRepository IngestorRepository) Config {
	return func(db *Database) *Database {
		db.IngestorRepository = ingestorRepository
		return db
	}
}
