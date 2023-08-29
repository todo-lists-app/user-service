package user

import (
	"context"
	ConfigBuilder "github.com/keloran/go-config"
	pb "github.com/todo-lists-app/protobufs/generated/user/v1"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	*ConfigBuilder.Config
}

func (s *Server) Delete(ctx context.Context, r *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	u := NewUserService(ctx, *s.Config, r.UserId, &RealMongoOperations{Collection: s.Mongo.Collections["user"], Database: s.Mongo.Database})
	err := u.DeleteUser()
	if err != nil {
		return &pb.UserDeleteResponse{
			UserId: r.UserId,
			Status: err.Error(),
		}, nil
	}

	return &pb.UserDeleteResponse{
		UserId: r.UserId,
		Status: "ok",
	}, nil
}
