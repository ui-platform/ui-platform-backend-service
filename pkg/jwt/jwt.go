package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	accessTokenDuration  = 24 * time.Hour
	refreshTokenDuration = 7 * 24 * time.Hour
	issuer               = "ui_platform_auth_service"
)

type CustomClaims struct {
	jwt.StandardClaims
	UserId    string `json:"user_id,omitempty"`
	TokenId   string `json:"token_id,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
}

func generateNonce() (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	return hex.EncodeToString(nonce), nil
}

func generateJWT(secretKey string, userId string, duration time.Duration, accessTokenHash string, tokenType string, nonce string) (string, error) {
	claims := CustomClaims{
		UserId:    userId,
		TokenId:   accessTokenHash,
		TokenType: tokenType,
		Nonce:     nonce,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateTokenPair(secretKey string, userID string) (accessToken string, refreshToken string, err error) {
	nonce, err := generateNonce()
	if err != nil {
		return "", "", err
	}

	accessToken, err = generateJWT(secretKey, userID, accessTokenDuration, "", "access", "")
	if err != nil {
		return "", "", err
	}

	hash := sha256.Sum256([]byte(accessToken))
	accessTokenHash := hex.EncodeToString(hash[:])

	refreshToken, err = generateJWT(secretKey, userID, refreshTokenDuration, accessTokenHash, "refresh", nonce)
	if err != nil {
		return "", "", err
	}

	// Сохраните nonce в базе данных
	if err := saveNonce(nonce, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateJWT(tokenString string, secretKey string, expectedTokenType string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.TokenType != expectedTokenType {
			return "", errors.New("invalid token type")
		}
		return claims.UserId, nil
	}

	return "", jwt.ErrInvalidKey
}

func RefreshTokens(refreshToken string, accessToken string, secretKey string) (newAccessToken string, newRefreshToken string, err error) {
	claims, err := ValidateRefreshToken(refreshToken, accessToken, secretKey)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err = generateJWT(secretKey, claims.UserId, accessTokenDuration, "", "access", "")
	if err != nil {
		return "", "", err
	}

	hash := sha256.Sum256([]byte(accessToken))
	accessTokenHash := hex.EncodeToString(hash[:])

	newRefreshToken, err = generateJWT(secretKey, claims.UserId, refreshTokenDuration, accessTokenHash, "refresh", claims.Nonce)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func ValidateRefreshToken(refreshToken string, accessToken string, secretKey string) (*CustomClaims, error) {
	hash := sha256.Sum256([]byte(accessToken))
	accessTokenHash := hex.EncodeToString(hash[:])

	token, err := jwt.ParseWithClaims(refreshToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.TokenType != "refresh" || claims.TokenId != accessTokenHash {
			return nil, errors.New("invalid refresh token or access token hash does not match")
		}
		valid, err := validateNonce(claims.Nonce, refreshToken)
		if err != nil || !valid {
			return nil, errors.New("invalid nonce or nonce check failed")
		}
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// Заглушки для функций saveNonce и validateNonce, которые вам нужно реализовать
func saveNonce(nonce string, refreshToken string) error {
	// Реализуйте сохранение nonce в вашей базе данных
	return nil
}

func validateNonce(nonce string, refreshToken string) (bool, error) {
	// Реализуйте проверку nonce в вашей базе данных
	return true, nil
}
