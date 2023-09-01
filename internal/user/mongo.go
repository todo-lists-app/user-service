package user

import (
	"context"
	"github.com/bugfixes/go-bugfixes/logs"
	mungo "github.com/keloran/go-config/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOperations interface {
	GetMongoClient(ctx context.Context, config mungo.Mongo) error
	Disconnect(ctx context.Context) error
	DeleteOne(ctx context.Context, filter interface{}) (interface{}, error)
}

type RealMongoOperations struct {
	Client     *mongo.Client
	Collection string
	Database   string
}

func (r *RealMongoOperations) GetMongoClient(ctx context.Context, config mungo.Mongo) error {
	client, err := mungo.GetMongoClient(ctx, config)
	if err != nil {
		return logs.Errorf("error getting mongo client: %v", err)
	}
	r.Client = client
	return nil
}

func (r *RealMongoOperations) Disconnect(ctx context.Context) error {
	return r.Client.Disconnect(ctx)
}

func (r *RealMongoOperations) DeleteOne(ctx context.Context, filter interface{}) (interface{}, error) {
	return r.Client.Database(r.Database).Collection(r.Collection).DeleteOne(ctx, filter)
}
