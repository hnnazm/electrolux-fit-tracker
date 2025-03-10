package mongo

import (
	"context"
	"os"
	"time"

	"fit-tracker/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const DATABASE_NAME = "FitTracker"

type (
	Mongo struct {
		Client *mongo.Client
	}
)

func New(ctx context.Context) (*Mongo, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("DATABASE_DSN")))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &Mongo{
		Client: client,
	}, nil
}

func (m Mongo) SaveTraces(ctx context.Context, input *database.SaveTracesInput) error {
	var (
		db  = m.Client.Database(DATABASE_NAME)
		col = db.Collection("ingestions")
	)

	input.CreatedAt = time.Now()

	_, err := col.InsertOne(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (m Mongo) GetTraces(ctx context.Context, input *database.GetTracesInput) ([]database.GetTracesResult, error) {
	var (
		db   = m.Client.Database(DATABASE_NAME)
		coll = db.Collection("ingestions")

		endOfDay = input.CreatedAt.Add(24 * time.Hour)

		filter = bson.M{
			"userId": input.UserID,
			"createdAt": bson.M{
				"$gte": input.CreatedAt,
				"$lt":  endOfDay,
			},
		}
	)

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []database.GetTracesResult
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
