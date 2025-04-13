package services

import (
	"github.com/rs/zerolog"
	"ui-platform-backend-service/internal/storages"
	"ui-platform-backend-service/pkg/rabbit_mq"
)

type Service struct {
	User    User
	Project Project
	Screen  Screen
}

type ServiceDeps struct {
	Log      zerolog.Logger
	Storage  *storages.Storage
	Producer *rabbit_mq.Producer
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		User:    NewUserService(deps.Log, deps.Producer, deps.Storage),
		Project: NewProjectService(deps.Log, deps.Producer, deps.Storage),
		Screen:  NewScreenService(deps.Log, deps.Storage),
	}
}
