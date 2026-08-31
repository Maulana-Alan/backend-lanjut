package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"
)

// Register daftar semua route
// Isi file ini cuma alamat + siapa yang melayani, gak ada logika bisnis
func Register(app *fiber.App, pool *pgxpool.Pool, studentService *service.StudentService) {
	api := app.Group("/api/v1")

	api.Get("/health", healthCheck(pool))

	s := api.Group("/students", middleware.RequireJSON)
	s.Get("/", studentService.List)
	s.Get("/:id", studentService.Get)
	s.Post("/", studentService.Create)
	s.Put("/:id", studentService.Replace)
	s.Patch("/:id", studentService.Patch)
	s.Delete("/:id", studentService.Delete)
}

// healthCheck cek server + database
func healthCheck(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable,
				"database tidak dapat dihubungi")
		}

		return helper.Success(c, fiber.StatusOK, "server dan database berjalan", nil)
	}
}
