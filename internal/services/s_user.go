package services

import (
	"fmt"
	"github.com/rs/zerolog"
	"math/rand"
	"time"
	"ui-platform-backend-service/internal/entity"
	"ui-platform-backend-service/internal/storages"
	"ui-platform-backend-service/pkg/rabbit_mq"
)

type User interface {
	EmailVerification(email string, sendCode bool) error
	Register(user entity.User, code string) (string, error)
	Login(user entity.User) (string, error)
}

type UserService struct {
	log      zerolog.Logger
	producer *rabbit_mq.Producer
	storage  *storages.Storage
}

func NewUserService(log zerolog.Logger, producer *rabbit_mq.Producer, storage *storages.Storage) *UserService {
	return &UserService{
		log:      log,
		producer: producer,
		storage:  storage,
	}
}

func (s *UserService) Register(user entity.User, code string) (string, error) {
	// проверяем код
	registerCode, err := s.storage.User.GetRegisterCode(user.Email)
	if err != nil {
		s.log.Error().Err(err).Msg("error getting register code")
		return "", fmt.Errorf("invalid code")
	}
	if registerCode != code {
		s.log.Error().Msg("invalid code")
		return "", fmt.Errorf("invalid code")
	}
	// хешируем пароль
	err = user.HashPassword()
	if err != nil {
		s.log.Error().Err(err).Msg("error hashing password")
		return "", fmt.Errorf("error hashing password")
	}
	// Проверяем уникальность email
	emailRegistered, err := s.storage.User.IsEmailRegistered(user.Email)
	if err != nil {
		s.log.Error().Err(err).Msg("error checking email registration")
		return "", fmt.Errorf("error checking email registration")
	}
	if emailRegistered {
		s.log.Error().Msg("email is already registered")
		return "", fmt.Errorf("email is already registered")
	}
	// сохраняем пользователя
	userId, err := s.storage.User.Create(user)
	if err != nil {
		s.log.Error().Err(err).Msg("error creating user")
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return "", fmt.Errorf("email is already registered")
		}
		return "", fmt.Errorf("error creating user")
	}

	return userId, nil
}

func (s *UserService) Login(user entity.User) (string, error) {
	userDb, err := s.storage.User.GetByEmail(user.Email)
	if err != nil {
		return "", err
	}
	user.PasswordHash = userDb.Password
	ok := user.CheckPassword()
	if !ok {
		return "", fmt.Errorf("invalid password")
	}

	return userDb.ID, nil
}

func (s *UserService) EmailVerification(email string, sendCode bool) error {

	emailRegistered, err := s.storage.User.IsEmailRegistered(email)
	if err != nil {
		return fmt.Errorf("error checking email registration")
	}

	if emailRegistered {
		return fmt.Errorf("email is already registered")
	}

	if sendCode {
		if s.storage.User.GetRegisterCodeEmailLock(email) {
			return fmt.Errorf("code is already sent, please wait 5 minutes")
		}
		err := s.storage.User.SetRegisterCodeEmailLock(email)
		if err != nil {
			return err
		}
		// генерируем код длиной 6 символов
		code := s.generateCode()
		// сохраняем в redis
		err = s.storage.User.SetRegisterCode(email, code)
		if err != nil {
			return err
		}
		// отправляем на почту
		err = s.producer.SendMessage("mail_verification", map[string]string{
			"email": email,
			"code":  code,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserService) generateCode() string {
	rand.Seed(time.Now().UnixNano())   // Инициализация генератора случайных чисел
	code := rand.Intn(900000) + 100000 // Генерация случайного числа от 100000 до 999999
	return fmt.Sprintf("%06d", code)   // Форматирование числа как строки с шестью цифрами
}
