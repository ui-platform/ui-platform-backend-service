package storages

import (
	"ui-platform-backend-service/internal/entity"
	"ui-platform-backend-service/pkg/database"
)

type Screen interface {
	Create(screen entity.Screen) (screenId string, err error)
	GetAllByProjectId(projectId string) (screens []entity.Screen, err error)
	GetById(screenId string) (screen entity.Screen, err error)
	DeleteById(screenId string) (err error)
}

type ScreenStorage struct {
	postgres *database.PostgresDB
	redis    *database.Redis
}

func NewScreenStorage(pg *database.PostgresDB, redis *database.Redis) *ScreenStorage {
	return &ScreenStorage{
		postgres: pg,
		redis:    redis,
	}
}

func (s *ScreenStorage) Create(screen entity.Screen) (screenId string, err error) {
	query := "INSERT INTO screens (project_id, name, description, status) VALUES ($1, $2, $3, $4) RETURNING id"
	var id string
	err = s.postgres.DB.QueryRow(query, screen.ProjectId, screen.Name, screen.Description, entity.ScreenStatusUnPublished).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *ScreenStorage) GetAllByProjectId(projectId string) (screens []entity.Screen, err error) {
	query := "SELECT id, project_id, name, description, status, content, settings, created_at, updated_at FROM screens WHERE project_id = $1 AND deleted_at IS NULL"
	rows, err := s.postgres.DB.Query(query, projectId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var screen entity.Screen
		err = rows.Scan(&screen.Id, &screen.ProjectId, &screen.Name, &screen.Description, &screen.Status, &screen.Content, &screen.Settings, &screen.CreatedAt, &screen.UpdatedAt)
		if err != nil {
			return nil, err
		}
		screens = append(screens, screen)
	}
	return screens, nil
}

func (s *ScreenStorage) GetById(screenId string) (screen entity.Screen, err error) {
	query := "SELECT id, project_id, name, description, status, content, settings, created_at, updated_at FROM screens WHERE id = $1 AND deleted_at IS NULL"
	var screenDB entity.Screen

	row := s.postgres.DB.QueryRow(query, screenId)
	err = row.Scan(&screenDB.Id, &screenDB.ProjectId, &screenDB.Name, &screenDB.Description, &screenDB.Status, &screenDB.Content, &screenDB.Settings, &screenDB.CreatedAt, &screenDB.UpdatedAt)
	if err != nil {
		return entity.Screen{}, err
	}
	return screenDB, nil
}

func (s *ScreenStorage) DeleteById(screenId string) (err error) {
	query := "UPDATE screens SET deleted_at = NOW() WHERE id = $1"
	_, err = s.postgres.DB.Exec(query, screenId)
	if err != nil {
		return err
	}
	return nil
}
