package helper

import (
	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
)

// Success kirim response 200 atau status lain yang berhasil
func Success(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(model.WebResponse{
		Success: true, Message: message, Data: data,
	})
}

// SuccessList kirim response 200 dengan meta paginasi
func SuccessList(c *fiber.Ctx, message string, data any, meta *model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Success: true, Message: message, Data: data, Meta: meta,
	})
}

// Created kirim 201 + header Location
func Created(c *fiber.Ctx, message string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Success: true, Message: message, Data: data,
	})
}

// NoContent kirim 204 tanpa body
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Fail kirim response gagal dengan status tertentu
func Fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(model.WebResponse{Success: false, Message: message})
}

// FailValidation kirim 422 dengan detail error per field
func FailValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(model.WebResponse{
		Success: false, Message: "validasi gagal", Errors: errs,
	})
}
