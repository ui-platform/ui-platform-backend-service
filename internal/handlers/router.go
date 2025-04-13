package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rs/zerolog"
	"time"
	"ui-platform-backend-service/internal/services"
	"ui-platform-backend-service/pkg/jwt"
)

type Handler struct {
	log        zerolog.Logger
	services   *services.Service
	jwtService *jwt.Service
}

func NewHandler(log zerolog.Logger, services *services.Service, jwtService *jwt.Service) *Handler {
	return &Handler{
		log:        log,
		services:   services,
		jwtService: jwtService,
	}
}

func (h *Handler) InitRoutes(port string) {
	app := fiber.New(fiber.Config{
		DisableDefaultContentType: true,
		CaseSensitive:             false,
	})

	app.Use(cors.New())

	// 3 requests per 10 seconds max
	app.Use(limiter.New(limiter.Config{
		Expiration: 1 * time.Second,
		Max:        10,
	}))

	api := app.Group("/api/v1")
	{
		api.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
		})

		// auth
		auth := api.Group("/auth")
		{
			auth.Use(limiter.New(limiter.Config{
				Expiration: 1 * time.Second,
				Max:        1,
			}))

			auth.Post("/email_verification", h.emailVerification)
			auth.Post("/register", h.register)
			auth.Post("/login", h.login)
			auth.Post("/refresh", h.refresh)
		}

		// projects
		projects := api.Group("/projects")
		{
			projects.Use(limiter.New(limiter.Config{
				Expiration: 1 * time.Second,
				Max:        10,
			}))

			projects.Use(h.middlewareAuth)

			projects.Post("/", h.createProject)
			projects.Get("/", h.getProjects)
			//projects.Get("/:id", nil)
			//projects.Put("/:id", nil)
			projects.Delete("/:project_id", h.deleteProject)
		}

	}

	h.log.Info().Msg("Starting server on port " + port)
	err := app.Listen(":" + port)
	if err != nil {
		h.log.Error().Msg(err.Error())
	}
}
