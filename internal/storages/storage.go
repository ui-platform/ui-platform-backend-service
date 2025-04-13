package storages

import (
	"github.com/rs/zerolog"
	"ui-platform-backend-service/pkg/database"
)

type Storage struct {
	User    User
	Project Project
	Screen  Screen
}

type StorageDeps struct {
	PostgresDB *database.PostgresDB
	Redis      *database.Redis
	Log        zerolog.Logger
}

func NewStorage(deps StorageDeps) *Storage {
	return &Storage{
		User:    NewUserStorage(deps.PostgresDB, deps.Redis),
		Project: NewProjectStorage(deps.PostgresDB, deps.Redis, deps.Log),
		Screen:  NewScreenStorage(deps.PostgresDB, deps.Redis),
	}
}
