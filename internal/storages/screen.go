package storages

import (
	"time"
	"ui-platform-backend-service/internal/entity"
	"ui-platform-backend-service/pkg/database"
)

type Screen interface {
	Create(screen *entity.Screen) (screenId string, err error)
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

func (s *ScreenStorage) Create(screen *entity.Screen) (screenId string, err error) {
	tx, err := s.postgres.DB.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// 1. Создаём экран
	query := "INSERT INTO screens (project_id, name, description, status, widgets, settings) VALUES ($1, $2, $3, $4) RETURNING id"
	err = tx.QueryRow(query, screen.ProjectId, screen.Name, screen.Description, entity.ScreenStatusUnPublished, map[string]interface{}{}, map[string]interface{}{}).Scan(&screenId)
	if err != nil {
		return "", err
	}

	// 2. Создаём ветки
	branchQuery := "INSERT INTO screens_branches (screen_id, name, created_at) VALUES ($1, $2, $3)"
	_, err = tx.Exec(branchQuery, screenId, "Main", time.Now().UTC())
	if err != nil {
		return "", err
	}

	return screenId, nil
}
