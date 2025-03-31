package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) deleteProject(c *fiber.Ctx) error {
	// Получаем userId из контекста
	userId := c.Locals("UID").(string)
	h.log.Debug().Msgf("userId: %v", userId)
	// Получаем projectId из параметров Path
	projectId := c.Params("project_id")
	h.log.Debug().Msgf("projectId: %v", projectId)
	if projectId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "project id is empty",
		})
	}
	// Удаляем проект
	err := h.services.Project.DeleteById(projectId)
	if err != nil {
		h.log.Error().Msgf("error deleting project: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error deleting project",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}
