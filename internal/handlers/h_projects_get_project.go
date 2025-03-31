package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) getProjects(c *fiber.Ctx) error {
	// Получаем userId из контекста
	userId := c.Locals("UID").(string)
	h.log.Debug().Msgf("userId: %v", userId)
	// Создаем проект
	projects, err := h.services.Project.GetAllByUserId(userId)
	if err != nil {
		h.log.Error().Msgf("error creating project: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error creating project",
		})
	}
	// Возвращаем projectId
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": fiber.Map{
			"projects": projects,
		},
	})
}
