package main

import (
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/todo-lists-app/user-service/internal/service"
)

var (
	BuildVersion = "dev"
	BuildHash    = "none"
	ServiceName  = "user-service"
)

func main() {
	logs.Local().Info(fmt.Sprintf("Starting %s", ServiceName))
	logs.Local().Info(fmt.Sprintf("Version: %s, Hash: %s", BuildVersion, BuildHash))

	s, err := service.NewService()
	if err != nil {
		_ = logs.Errorf("new service: %v", err)
		return
	}

	if err := s.Start(); err != nil {
		_ = logs.Errorf("start service: %v", err)
		return
	}
}
