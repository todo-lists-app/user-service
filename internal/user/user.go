package user

import (
	"context"
	"github.com/bugfixes/go-bugfixes/logs"
	ConfigBuilder "github.com/keloran/go-config"
	mungo "github.com/keloran/go-config/mongo"
	"go.mongodb.org/mongo-driver/bson"
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

type Service struct {
	ConfigBuilder.Config
	context.Context
	UserID string

	MongoOps MongoOperations
}

func NewUserService(ctx context.Context, config ConfigBuilder.Config, userID string, mongoOps MongoOperations) *Service {
	return &Service{
		Config:   config,
		Context:  ctx,
		UserID:   userID,
		MongoOps: mongoOps,
	}
}

func (s *Service) DeleteUser() error {
	if err := s.MongoOps.GetMongoClient(s.Context, s.Config.Mongo); err != nil {
		return logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := s.MongoOps.Disconnect(s.Context); err != nil {
			_ = logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	filter := bson.D{{"userid", s.UserID}}
	_, err := s.MongoOps.DeleteOne(s.Context, filter)
	if err != nil {
		return logs.Errorf("error deleting user: %v", err)
	}
	return nil
}
