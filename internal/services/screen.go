package services

import (
	"github.com/rs/zerolog"
	"ui-platform-backend-service/internal/entity"
	"ui-platform-backend-service/internal/storages"
)

type Screen interface {
	Create(screen *entity.Screen) (screenId string, err error)
}

type screenService struct {
	log     zerolog.Logger
	storage *storages.Storage
}

func NewScreenService(log zerolog.Logger, storage *storages.Storage) Screen {
	return &screenService{
		log:     log,
		storage: storage,
	}
}

func (s *screenService) Create(screen *entity.Screen) (screenId string, err error) {
	s.log.Info().Str("screen_id", screen.Id).Msg("Creating screen")
	//TODO: implement
	return "", err
}
