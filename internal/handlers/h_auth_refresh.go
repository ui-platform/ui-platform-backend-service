package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) refresh(c *fiber.Ctx) error {
	// Получаем accessToken из заголовка
	accessToken := c.Get("Authorization")
	// Проверяем accessToken
	if accessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "access token is empty",
		})
	}
	// Убираем префикс Bearer
	accessToken = accessToken[7:]
	// Получаем refreshToken из заголовка
	refreshToken := c.Get("Refresh-Token")
	// Проверяем refreshToken
	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "refresh token is empty",
		})
	}
	// Обновляем токены
	newAccessToken, newRefreshToken, err := h.jwtService.RefreshTokens(refreshToken, accessToken)
	if err != nil {
		h.log.Error().Err(err).Msg("error refreshing tokens")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error refreshing tokens",
		})
	}
	// Возвращаем токены
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": fiber.Map{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		},
	})
}
