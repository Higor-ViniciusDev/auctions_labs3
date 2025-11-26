package mongodb

import (
	"context"
	"os"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	MONGODB_URL = "MONGODB_URL"
	MONGODB_DB  = "MONGODB_DB"
)

func NewConnectionDataBaseMongoDB(ctx context.Context) (*mongo.Database, error) {
	mongoURL := os.Getenv(MONGODB_URL)
	mongoDB := os.Getenv(MONGODB_DB)

	client, err := mongo.Connect(
		options.Client().ApplyURI(mongoURL),
	)

	if err != nil {
		logger.Error("Error ao tentar conectar no mongoDB database", err)
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		logger.Error("Error ao tentar pingar no mongoDB database", err)
		return nil, err
	}

	return client.Database(mongoDB), nil
}
