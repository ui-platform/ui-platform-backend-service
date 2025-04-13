package services

import (
	"github.com/rs/zerolog"
	"ui-platform-backend-service/internal/entity"
	"ui-platform-backend-service/internal/storages"
	"ui-platform-backend-service/pkg/rabbit_mq"
)

type Project interface {
	Create(project entity.Project, ownerId string) (projectId string, err error)
	GetAllByUserId(userId string) (projects []entity.Project, err error)
	UpdateById(project entity.Project) (err error)
	DeleteById(projectId string) (err error)
}

type ProjectService struct {
	log      zerolog.Logger
	producer *rabbit_mq.Producer
	storage  *storages.Storage
}

func NewProjectService(log zerolog.Logger, producer *rabbit_mq.Producer, storage *storages.Storage) *ProjectService {
	return &ProjectService{
		log:      log,
		producer: producer,
		storage:  storage,
	}
}

func (s *ProjectService) Create(project entity.Project, ownerId string) (projectId string, err error) {
	s.log.Debug().Msgf("creating project: %v", project)
	s.log.Debug().Msgf("ownerId: %v", ownerId)
	return s.storage.Project.Create(project, ownerId)
}

func (s *ProjectService) GetAllByUserId(userId string) (projects []entity.Project, err error) {
	projects, err = s.storage.Project.GetAllByUserId(userId)
	if err != nil {
		return nil, err
	}
	if projects == nil {
		return []entity.Project{}, nil
	}
	return projects, nil
}

func (s *ProjectService) UpdateById(project entity.Project) (err error) {
	return s.storage.Project.UpdateById(project)
}

func (s *ProjectService) DeleteById(projectId string) (err error) {
	return s.storage.Project.DeleteById(projectId)
}
