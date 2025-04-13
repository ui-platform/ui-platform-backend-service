package handlers

import (
	"github.com/gofiber/fiber/v2"
	"ui-platform-backend-service/internal/entity"
)

func (h *Handler) register(c *fiber.Ctx) error {
	var user entity.UserRegister
	// Парсим тело запроса
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request body",
		})
	}
	// Проверяем тело запроса
	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Регистрируем
	userId, err := h.services.User.Register(entity.User{
		Email:    user.Email,
		Password: user.Password,
	}, user.Code)
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
