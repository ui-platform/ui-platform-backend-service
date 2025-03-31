package storages

import (
	"errors"
	"time"
	"ui-platform-backend-service/internal/entity"
	"ui-platform-backend-service/pkg/database"
)

type Project interface {
	Create(project entity.Project, ownerId string) (projectId string, err error)
	GetAllByUserId(userId string) (projects []entity.Project, err error)
	UpdateById(projects entity.Project) (err error)
	DeleteById(projectId string) (err error)
}

type ProjectStorage struct {
	postgres *database.PostgresDB
	redis    *database.Redis
}

func NewProjectStorage(pg *database.PostgresDB, redis *database.Redis) *ProjectStorage {
	return &ProjectStorage{
		postgres: pg,
		redis:    redis,
	}
}

func (s *ProjectStorage) Create(project entity.Project, ownerId string) (projectId string, err error) {
	// create project by transaction
	tx, err := s.postgres.DB.Begin()
	if err != nil {
		return "", err
	}

	// create project
	queryCreateProject := `INSERT INTO projects (name, description, status) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(queryCreateProject, project.Name, project.Description, entity.ProjectStatusUnPublished).Scan(&projectId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return "", err
		}
		return "", err
	}

	// add user to project
	queryAddUserToProject := `INSERT INTO projects_membership (project_id, user_id, is_owner) VALUES ($1, $2, $3)`
	_, err = tx.Exec(queryAddUserToProject, projectId, ownerId, true)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return "", err
		}
		return "", err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return projectId, nil
}

func (s *ProjectStorage) GetAllByUserId(userId string) ([]entity.Project, error) {
	query := `
        SELECT p.id, p.name, p.description, p.status, p.created_at
        FROM projects p
        JOIN projects_membership pu ON p.id = pu.project_id
        WHERE pu.user_id = $1 AND p.deleted_at IS NULL
    `

	rows, err := s.postgres.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

	var projects []entity.Project
	for rows.Next() {
		var project entity.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.Status, &project.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *ProjectStorage) UpdateById(project entity.Project) error {
	query := `
        UPDATE projects
        SET name = $2, description = $3, status = $4, updated_at = $5
        WHERE id = $1 AND deleted_at IS NULL
    `

	// Подготовка текущего времени для обновления поля updated_at.
	now := time.Now().UTC()

	// Выполнение запроса с передачей параметров.
	result, err := s.postgres.DB.Exec(query, project.ID, project.Name, project.Description, project.Status, now)
	if err != nil {
		return err
	}

	// Проверка, что была обновлена хотя бы одна строка.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected, project may not exist or may be deleted")
	}

	return nil
}

func (s *ProjectStorage) DeleteById(projectId string) error {
	query := `UPDATE projects SET deleted_at = $1 WHERE id = $2`
	_, err := s.postgres.DB.Exec(query, time.Now().UTC(), projectId)
	if err != nil {
		return err
	}
	return nil
}
