package handlers

import (
	"github.com/gofiber/fiber/v2"
	"ui-platform-backend-service/internal/entity"
)

func (h *Handler) createProject(c *fiber.Ctx) error {
	// Получаем userId из контекста
	userId := c.Locals("UID").(string)
	h.log.Debug().Msgf("userId: %v", userId)
	// Парсим тело запроса
	var project entity.Project
	if err := c.BodyParser(&project); err != nil {
		h.log.Error().Msgf("invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}
	// Проверяем валидность данных
	if err := project.Validate(); err != nil {
		h.log.Error().Msgf("invalid project data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid project data",
		})
	}
	// Создаем проект
	projectId, err := h.services.Project.Create(project, userId)
	if err != nil {
		h.log.Error().Msgf("error creating project: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error creating project",
		})
	}
	// Возвращаем projectId
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "ok",
		"details": fiber.Map{
			"project_id": projectId,
		},
	})
}
