package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(connectionUrl string) (*mongo.Database, error) {
	opts := options.Client().ApplyURI(connectionUrl)
	cl, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	return cl.Database(""), nil
}