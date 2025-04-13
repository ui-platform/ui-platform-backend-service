package storages

import "ui-platform-backend-service/pkg/database"

type Storage struct {
	User    *UserStorage
	Project *ProjectStorage
	Screen  *ScreenStorage
}

func NewStorage(pg *database.PostgresDB, redis *database.Redis) *Storage {
	return &Storage{
		User:    NewUserStorage(pg, redis),
		Project: NewProjectStorage(pg, redis),
		Screen:  NewScreenStorage(pg, redis),
	}
}
