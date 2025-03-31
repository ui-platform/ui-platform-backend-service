package storages

import (
	"fmt"
	"time"
	"ui-platform-backend-service/internal/entity"
	"ui-platform-backend-service/pkg/database"
)

type User interface {
	SetRegisterCode(email string, code string) error
	GetRegisterCode(email string) (string, error)
	SetRegisterCodeEmailLock(email string) error
	GetRegisterCodeEmailLock(email string) bool
	IsEmailRegistered(email string) (bool, error)
	Create(user entity.User) (string, error)
	GetByEmail(email string) (entity.User, error)
	GetById(id string) (entity.User, error)
}

type UserStorage struct {
	postgres *database.PostgresDB
	redis    *database.Redis
}

func NewUserStorage(pg *database.PostgresDB, redis *database.Redis) *UserStorage {
	return &UserStorage{
		postgres: pg,
		redis:    redis,
	}
}

func (s *UserStorage) SetRegisterCode(email string, code string) error {
	key := fmt.Sprintf("register_code:%s", email)

	err := s.redis.Client.Set(key, code, time.Minute*15).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) GetRegisterCode(email string) (string, error) {
	key := fmt.Sprintf("register_code:%s", email)
	code, err := s.redis.Client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *UserStorage) SetRegisterCodeEmailLock(email string) error {
	key := fmt.Sprintf("register_code_email_lock:%s", email)
	err := s.redis.Client.Set(key, true, time.Minute*5).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) GetRegisterCodeEmailLock(email string) bool {
	key := fmt.Sprintf("register_code_email_lock:%s", email)
	val, err := s.redis.Client.Get(key).Result()
	if err != nil {
		return false
	}

	return val == "1"
}

func (s *UserStorage) IsEmailRegistered(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = $1 AND deleted_at IS NULL"
	err := s.postgres.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *UserStorage) Create(user entity.User) (string, error) {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	var id string
	err := s.postgres.DB.QueryRow(query, user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *UserStorage) GetByEmail(email string) (entity.User, error) {
	var user entity.User
	query := "SELECT id, email, password FROM users WHERE email = $1 AND deleted_at IS NULL"
	err := s.postgres.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (s *UserStorage) GetById(id string) (entity.User, error) {
	var user entity.User
	query := "SELECT id, email, password FROM users WHERE id = $1 AND deleted_at IS NULL"
	err := s.postgres.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
