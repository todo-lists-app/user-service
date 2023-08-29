package service

import (
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	ConfigBuilder "github.com/keloran/go-config"
	"github.com/keloran/go-healthcheck"
	pb "github.com/todo-lists-app/protobufs/generated/user/v1"
	"github.com/todo-lists-app/user-service/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"time"
)

type Service struct {
	Config *ConfigBuilder.Config
}

func NewService() (*Service, error) {
	c, err := ConfigBuilder.Build(ConfigBuilder.Local, ConfigBuilder.Mongo, ConfigBuilder.Keycloak)
	if err != nil {
		return nil, logs.Errorf("build config: %v", err)
	}

	return &Service{
		Config: c,
	}, nil
}

func (s *Service) Start() error {
	errChan := make(chan error)
	go startHTTP(s.Config, errChan)
	go startGRPC(s.Config, errChan)

	return <-errChan
}

func startGRPC(config *ConfigBuilder.Config, errChan chan error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Local.GRPCPort))
	if err != nil {
		errChan <- err
		return
	}

	gs := grpc.NewServer()
	reflection.Register(gs)
	pb.RegisterUserServiceServer(gs, &user.Server{
		Config: config,
	})
	logs.Local().Infof("starting grpc on port %d", config.Local.GRPCPort)
	if err := gs.Serve(lis); err != nil {
		errChan <- err
	}
}

func startHTTP(config *ConfigBuilder.Config, errChan chan error) {
	allowedOrigins := []string{
		"http://localhost:3000",
		"https://api.todo-list.app",
		"https://todo-list.app",
		"https://beta.todo-list.app",
	}
	if config.Local.Development {
		allowedOrigins = append(allowedOrigins, "http://*")
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			"GET",
		},
	}))
	r.Get("/health", healthcheck.HTTP)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Local.HTTPPort),
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	logs.Local().Infof("starting http on port %d", config.Local.HTTPPort)
	if err := srv.ListenAndServe(); err != nil {
		errChan <- err
	}
}
