package db_connection

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context, mongoURI, dbName string) (*mongo.Client, *mongo.Collection, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, nil, fmt.Errorf("db_connection: failed to connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("db_connection: failed to ping: %w", err)
	}

	col := client.Database(dbName).Collection("spans")

	ttlSeconds := int32(3600)
	_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "startTime", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(ttlSeconds),
		},
		{
			Keys: bson.D{{Key: "traceId", Value: 1}},
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("db_connection: failed to create indexes: %w", err)
	}

	return client, col, nil
}
