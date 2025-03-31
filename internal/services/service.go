package services

import (
	"github.com/rs/zerolog"
	"ui-platform-backend-service/internal/storages"
	"ui-platform-backend-service/pkg/rabbit_mq"
)

type Service struct {
	User    User
	Project Project
}

func NewService(log zerolog.Logger, producer *rabbit_mq.Producer, storage *storages.Storage) *Service {
	return &Service{
		User:    NewUserService(log, producer, storage),
		Project: NewProjectService(log, producer, storage),
	}
}
