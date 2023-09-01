package user

import (
	"context"
	"github.com/bugfixes/go-bugfixes/logs"
	ConfigBuilder "github.com/keloran/go-config"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	ConfigBuilder.Config
	context.Context
	UserID      string
	AccessToken string

	MongoOps MongoOperations
}

func NewUserService(ctx context.Context, config ConfigBuilder.Config, userID, accessToken string, mongoOps MongoOperations) *Service {
	return &Service{
		Config:      config,
		Context:     ctx,
		UserID:      userID,
		AccessToken: accessToken,
		MongoOps:    mongoOps,
	}
}

func (s *Service) DeleteUser() error {
	if err := s.deleteFromMongo(); err != nil {
		return logs.Errorf("error deleting from mongo: %v", err)
	}

	if err := s.deleteFromKeycloak(); err != nil {
		return logs.Errorf("error deleting from keycloak: %v", err)
	}

	return nil
}

func (s *Service) deleteFromMongo() error {
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

func (s *Service) deleteFromKeycloak() error {
	if err := DeleteKeyCloakUser(s.Context, s.Config, s.UserID, s.AccessToken); err != nil {
		return logs.Errorf("error deleting user: %v", err)
	}
	return nil
}
