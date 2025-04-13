// Package jwt предоставляет сервис для генерации, валидации и обновления
// JWT access/refresh токенов с привязкой к nonce-хранилищу.
//
// Он предназначен для использования в микросервисной архитектуре, где
// аутентификация пользователя разделена по сервисам.
//
// Формат токенов: HMAC-SHA256, со встроенными кастомными клеймами.

package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer = "ui-platform-auth-service"
)

// Config содержит параметры конфигурации для JWT-сервиса.
type Config struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Service реализует логику генерации и валидации JWT-токенов.
type Service struct {
	cfg          Config
	nonceStorage NonceStorage
}

// NonceStorage описывает интерфейс для хранения и проверки nonce.
// Используется для обеспечения one-time использования refresh-токенов.
type NonceStorage interface {
	Save(nonce, token string) error
	Validate(nonce, token string) (bool, error)
}

// CustomClaims расширяет стандартные JWT claims специфичными полями
// для пользовательского идентификатора, типа токена, nonce и хеша access-токена.
type CustomClaims struct {
	UserId    string `json:"user_id,omitempty"`
	TokenId   string `json:"token_id,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	jwt.RegisteredClaims
}

// New создает экземпляр JWT-сервиса с заданной конфигурацией и хранилищем nonce.
func New(cfg Config, ns NonceStorage) *Service {
	return &Service{cfg: cfg, nonceStorage: ns}
}

// generateNonce создает криптографически безопасный уникальный идентификатор nonce.
func generateNonce() (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	return hex.EncodeToString(nonce), nil
}

// generateJWT создает JWT с заданными параметрами.
// Используется как для access, так и для refresh токенов.
func (s *Service) generateJWT(userId string, duration time.Duration, accessTokenHash, tokenType, nonce string) (string, error) {
	claims := CustomClaims{
		UserId:    userId,
		TokenId:   accessTokenHash,
		TokenType: tokenType,
		Nonce:     nonce,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.SecretKey))
}

// GenerateTokenPair создает пару access/refresh токенов.
// В refresh-токен вшивается хеш access-токена и nonce.
func (s *Service) GenerateTokenPair(userID string) (string, string, error) {
	nonce, err := generateNonce()
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.generateJWT(userID, s.cfg.AccessTokenTTL, "", "access", "")
	if err != nil {
		return "", "", err
	}

	hash := sha256.Sum256([]byte(accessToken))
	accessTokenHash := hex.EncodeToString(hash[:])

	refreshToken, err := s.generateJWT(userID, s.cfg.RefreshTokenTTL, accessTokenHash, "refresh", nonce)
	if err != nil {
		return "", "", err
	}

	if s.nonceStorage == nil {
		return accessToken, refreshToken, nil
	}

	if err := s.nonceStorage.Save(nonce, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateJWT проверяет валидность JWT и соответствие ожидаемому типу ("access"/"refresh").
// Возвращает userID, если токен валиден.
func (s *Service) ValidateJWT(tokenStr, expectedType string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.SecretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims.TokenType != expectedType {
		return "", errors.New("unexpected token type")
	}

	return claims.UserId, nil
}

// RefreshTokens валидирует refresh-токен и access-токен, генерирует новую пару токенов.
// Повторное использование одного и того же refresh-токена не допускается.
func (s *Service) RefreshTokens(refreshToken, accessToken string) (string, string, error) {
	claims, err := s.validateRefreshToken(refreshToken, accessToken)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := s.generateJWT(claims.UserId, s.cfg.AccessTokenTTL, "", "access", "")
	if err != nil {
		return "", "", err
	}

	hash := sha256.Sum256([]byte(newAccessToken))
	accessTokenHash := hex.EncodeToString(hash[:])

	newRefreshToken, err := s.generateJWT(claims.UserId, s.cfg.RefreshTokenTTL, accessTokenHash, "refresh", claims.Nonce)
	if err != nil {
		return "", "", err
	}

	if s.nonceStorage == nil {
		return newAccessToken, newRefreshToken, nil
	}

	if err := s.nonceStorage.Save(claims.Nonce, newRefreshToken); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// validateRefreshToken выполняет полную проверку refresh-токена:
// проверку подписи, типа, связи с access-токеном и валидацию nonce.
func (s *Service) validateRefreshToken(refreshToken, accessToken string) (*CustomClaims, error) {
	hash := sha256.Sum256([]byte(accessToken))
	accessTokenHash := hex.EncodeToString(hash[:])

	token, err := jwt.ParseWithClaims(refreshToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" || claims.TokenId != accessTokenHash {
		return nil, errors.New("refresh token hash mismatch")
	}

	if s.nonceStorage == nil {
		return claims, nil
	}

	valid, err := s.nonceStorage.Validate(claims.Nonce, refreshToken)
	if err != nil || !valid {
		return nil, errors.New("nonce validation failed")
	}

	return claims, nil
}
