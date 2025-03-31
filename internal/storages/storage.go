package storages

import "ui-platform-backend-service/pkg/database"

type Storage struct {
	User    *UserStorage
	Project *ProjectStorage
}

func NewStorage(pg *database.PostgresDB, redis *database.Redis) *Storage {
	return &Storage{
		User:    NewUserStorage(pg, redis),
		Project: NewProjectStorage(pg, redis),
	}
}
