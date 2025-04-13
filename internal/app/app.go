package app

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"os"
	"time"
	"ui-platform-backend-service/internal/config"
	"ui-platform-backend-service/internal/handlers"
	"ui-platform-backend-service/internal/services"
	"ui-platform-backend-service/internal/storages"
	"ui-platform-backend-service/pkg/database"
	"ui-platform-backend-service/pkg/jwt"
	"ui-platform-backend-service/pkg/rabbit_mq"
)

func Run() {
	// logger
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"}
	logger := zerolog.New(output).With().Caller().Timestamp().Logger()
	dir, err := os.Getwd()
	if err != nil {
		logger.Error().Msg("Cannot get working directory")
	}
	logger.Info().Msg("Current working directory: " + dir)
	// init env
	err = godotenv.Load()
	if err != nil {
		logger.Error().Msgf("Error loading .env file: %v", err)
	}
	// cfg
	cfg := config.GetConfig()
	logger.Info().Msg("Config: OK")
	// rabbitmq
	producer, err := rabbit_mq.NewProducer("amqp://"+cfg.RabbitMQ.User+":"+cfg.RabbitMQ.Password+"@"+cfg.RabbitMQ.Host+":"+cfg.RabbitMQ.Port, logger)
	if err != nil {
		logger.Error().Msgf("Error connecting to RabbitMQ: %v", err)
	}
	logger.Info().Msg("RabbitMQ: OK")
	// postgres
	pg, err := database.NewPostgresDB(cfg.Postgres.DBHost, cfg.Postgres.DBPort, cfg.Postgres.DBUser, cfg.Postgres.DBName, cfg.Postgres.DBPass, cfg.Postgres.DBSSLMode)
	if err != nil {
		logger.Error().Msgf("Error connecting to PostgreSQL: %v", err)
	}
	logger.Info().Msg("Postgres: OK")
	// redis
	redis, err := database.NewRedis(cfg.Redis.Host+":"+cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Error().Msgf("Error connecting to Redis: %v", err)
	}
	logger.Info().Msg("Redis: OK")
	// storage
	storage := storages.NewStorage(storages.StorageDeps{
		PostgresDB: pg,
		Redis:      redis,
		Log:        logger,
	})
	// services
	service := services.NewService(services.ServiceDeps{
		Log:      logger,
		Producer: producer,
		Storage:  storage,
	})
	// jwt service
	jwtService := jwt.New(jwt.Config{
		SecretKey:       cfg.AppSecretKey,
		AccessTokenTTL:  time.Hour * 24,
		RefreshTokenTTL: time.Hour * 24 * 7,
	}, nil)
	// handlers
	handler := handlers.NewHandler(logger, service, jwtService)
	// run
	handler.InitRoutes(cfg.AppPort)
}
