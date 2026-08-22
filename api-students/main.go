package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// metodeBerbody = metode yang punya body (harus cek Content-Type)
var metodeBerbody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

// requireJSON menolak request yang body-nya bukan JSON — status 415
func requireJSON(c *fiber.Ctx) error {
	if metodeBerbody[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return fail(c, fiber.StatusUnsupportedMediaType, "Content-Type harus application/json")
		}
	}
	return c.Next()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "API Students — Pertemuan 2",
		// Error handler global: tangkap panic dan error tak terduga
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			pesan := "terjadi kesalahan pada server"
			if e, ok := err.(*fiber.Error); ok {
				status = e.Code
				pesan = e.Message
			}
			return fail(c, status, pesan)
		},
	})

	// Middleware global — urutan penting: requestid dulu, baru logger pakai ID-nya
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${locals:requestid} ${method} ${path} → ${status} (${latency})\n",
	}))
	app.Use(cors.New())

	// Route group
	api := app.Group("/api/v1")

	// Health check — bisa dipakai untuk cek server masih jalan
	api.Get("/health", func(c *fiber.Ctx) error {
		return ok(c, "server berjalan", nil)
	})

	// requireJSON dipasang di grup /students saja, bukan global
	s := api.Group("/students", requireJSON)
	s.Get("/", listStudents)
	s.Get("/:id", getStudent)
	s.Post("/", createStudent)
	s.Put("/:id", replaceStudent)
	s.Patch("/:id", patchStudent)
	s.Delete("/:id", deleteStudent)

	// Catch-all: endpoint tidak dikenal
	app.Use(func(c *fiber.Ctx) error {
		return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	log.Fatal(app.Listen(":3000"))
}
