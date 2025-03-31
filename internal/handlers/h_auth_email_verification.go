package handlers

import (
	"github.com/gofiber/fiber/v2"
	"ui-platform-backend-service/internal/entity"
)

func (h *Handler) emailVerification(c *fiber.Ctx) error {
	sendCode := c.Query("send_code")
	var user entity.User
	// Парсим тело запроса
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request body",
		})
	}
	// Проверяем email
	if err := user.ValidateEmail(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	h.log.Debug().Msgf("email: %s", user.Email)
	err := h.services.User.EmailVerification(user.Email, sendCode == "true")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Возвращаем OK
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}
