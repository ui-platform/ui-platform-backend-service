package handlers

import (
	"github.com/gofiber/fiber/v2"
	"ui-platform-backend-service/internal/entity"
)

func (h *Handler) login(c *fiber.Ctx) error {
	// Получаем заголовки в запросе
	h.log.Debug().Msgf("headers: %v", c.Get("User-Agent"))
	var user entity.User
	// Парсим тело запроса
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request body",
		})
	}
	// Проверяем тело запроса
	if err := user.ValidateLogin(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	//
	userId, err := h.services.User.Login(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Создаем токены
	accessToken, refreshToken, err := h.jwtService.GenerateTokenPair(userId)
	if err != nil {
		h.log.Error().Err(err).Msg("error generating tokens")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error generating tokens",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}
