package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"api-students/helper"
)

// Register pasang semua middleware global (urutan penting!)
func Register(app *fiber.App, logger *slog.Logger) {
	app.Use(requestid.New())       // 1. ID unik tiap request
	app.Use(recover.New())         // 2. tangkap panic supaya server gak crash
	app.Use(helmet.New())          // 3. header keamanan dasar
	app.Use(cors.New())            // 4. izinkan request dari browser
	app.Use(RequestLogger(logger)) // 5. catat setiap request
}

// RequestLogger catat setiap request ke log terstruktur
func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next() // lanjut ke handler

		requestID, _ := c.Locals("requestid").(string)

		logger.Info("http_request",
			slog.String("request_id", requestID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", c.IP()),
		)

		return err
	}
}

var methodsWithBody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

// RequireJSON tolak request yang Content-Type-nya bukan JSON
func RequireJSON(c *fiber.Ctx) error {
	if methodsWithBody[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return helper.Fail(c, fiber.StatusUnsupportedMediaType,
				"Content-Type harus application/json")
		}
	}
	return c.Next()
}
