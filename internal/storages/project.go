package storages

import (
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"ui-platform-backend-service/internal/entity"
	"ui-platform-backend-service/pkg/database"
)

type Project interface {
	Create(project entity.Project, ownerId string) (projectId string, err error)
	GetAllByUserId(userId string) ([]entity.Project, error)
	UpdateById(project entity.Project) error
	DeleteById(projectId string) error
}

type ProjectStorage struct {
	postgres *database.PostgresDB
	redis    *database.Redis
	log      zerolog.Logger
}

func NewProjectStorage(pg *database.PostgresDB, redis *database.Redis, log zerolog.Logger) *ProjectStorage {
	return &ProjectStorage{
		postgres: pg,
		redis:    redis,
		log:      log,
	}
}

func (s *ProjectStorage) Create(project entity.Project, ownerId string) (string, error) {
	s.log.Debug().Str("ownerId", ownerId).Interface("project", project).Msg("creating new project")

	tx, err := s.postgres.DB.Begin()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to begin transaction")
		return "", err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	queryCreateProject := `INSERT INTO projects (name, description, status) VALUES ($1, $2, $3) RETURNING id`
	var projectId string
	err = tx.QueryRow(queryCreateProject, project.Name, project.Description, entity.ProjectStatusUnPublished).Scan(&projectId)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to insert project")
		tx.Rollback()
		return "", err
	}

	queryAddUserToProject := `INSERT INTO projects_membership (project_id, user_id, is_owner) VALUES ($1, $2, $3)`
	_, err = tx.Exec(queryAddUserToProject, projectId, ownerId, true)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to add user to project")
		tx.Rollback()
		return "", err
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("failed to commit transaction")
		return "", err
	}

	s.log.Debug().Str("projectId", projectId).Msg("project created successfully")
	return projectId, nil
}

func (s *ProjectStorage) GetAllByUserId(userId string) ([]entity.Project, error) {
	s.log.Debug().Str("userId", userId).Msg("fetching all projects for user")

	query := `
		SELECT p.id, p.name, p.description, p.status, p.created_at
		FROM projects p
		JOIN projects_membership pu ON p.id = pu.project_id
		WHERE pu.user_id = $1 AND p.deleted_at IS NULL
	`

	rows, err := s.postgres.DB.Query(query, userId)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to query projects")
		return nil, err
	}
	defer rows.Close()

	var projects []entity.Project
	for rows.Next() {
		var project entity.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.Status, &project.CreatedAt); err != nil {
			s.log.Error().Err(err).Msg("failed to scan project row")
			return nil, err
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		s.log.Error().Err(err).Msg("rows iteration error")
		return nil, err
	}

	s.log.Debug().Int("count", len(projects)).Msg("projects fetched successfully")
	return projects, nil
}

func (s *ProjectStorage) UpdateById(project entity.Project) error {
	s.log.Debug().Str("projectId", project.ID).Msg("updating project")

	query := `
		UPDATE projects
		SET name = $2,
			description = $3,
			status = $4,
			updated_at = $5
		WHERE id = $1 AND deleted_at IS NULL
	`

	now := time.Now().UTC()
	result, err := s.postgres.DB.Exec(query, project.ID, project.Name, project.Description, project.Status, now)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to update project")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to fetch affected rows")
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected, project may not exist or may be deleted")
	}

	s.log.Debug().Msg("project updated successfully")
	return nil
}

func (s *ProjectStorage) DeleteById(projectId string) error {
	s.log.Debug().Str("projectId", projectId).Msg("soft-deleting project")

	query := `UPDATE projects SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	res, err := s.postgres.DB.Exec(query, time.Now().UTC(), projectId)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to delete project")
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to fetch affected rows on delete")
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	s.log.Debug().Msg("project deleted successfully")
	return nil
}
