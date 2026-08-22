package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// --- Fungsi-fungsi untuk kirim respons ---

func ok(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(WebResponse{
		Success: true, Message: message, Data: data,
	})
}

func okList(c *fiber.Ctx, message string, data any, meta *Meta) error {
	return c.Status(fiber.StatusOK).JSON(WebResponse{
		Success: true, Message: message, Data: data, Meta: meta,
	})
}

// created: status 201, sekaligus set header Location supaya client tahu di mana data barunya
func created(c *fiber.Ctx, message string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(WebResponse{
		Success: true, Message: message, Data: data,
	})
}

// noContent: status 204, tidak ada body — dipakai setelah DELETE berhasil
func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(WebResponse{Success: false, Message: message})
}

// failValidation: status 422, sebutkan field mana yang salah
func failValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(WebResponse{
		Success: false, Message: "validasi gagal", Errors: errs,
	})
}

// --- Query string parser ---

// Field yang boleh dipakai untuk sort — kalau client kirim field lain, diabaikan
var allowedSort = map[string]bool{
	"id": true, "name": true, "nim": true, "grade": true,
}

// parseListQuery membaca query string dan memberi nilai default yang aman
func parseListQuery(c *fiber.Ctx) ListQuery {
	q := ListQuery{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: strings.TrimSpace(c.Query("search")),
		Sort:   c.Query("sort", "id"),
		Order:  strings.ToLower(c.Query("order", "asc")),
	}

	// Pastikan nilai tidak aneh
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 { // batas atas wajib ada, cegah ?limit=999999
		q.Limit = 100
	}
	if !allowedSort[q.Sort] { // kalau field tidak ada di daftar putih, pakai default
		q.Sort = "id"
	}
	if q.Order != "desc" {
		q.Order = "asc"
	}

	// Filter is_active: nil = tidak difilter, true/false = difilter
	if raw := c.Query("is_active"); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			q.IsActive = &v
		}
	}

	return q
}
