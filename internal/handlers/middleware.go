package handlers

import (
	"github.com/gofiber/fiber/v2"
	"ui-platform-backend-service/pkg/jwt"
)

func (h *Handler) middlewareAuth(c *fiber.Ctx) error {
	// Получаем accessToken из заголовка
	accessToken := c.Get("Authorization")
	// Проверяем accessToken
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "access token is empty",
		})
	}
	// Убираем префикс Bearer
	accessToken = accessToken[7:]
	// Валидируем accessToken
	userId, err := jwt.ValidateJWT(accessToken, h.appSecretKey, "access")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid access token",
		})
	}
	// Сохраняем userId в контексте
	c.Locals("UID", userId)
	// Пропускаем запрос
	return c.Next()
}
