package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
)

// Handler struct — repo disuntikkan dari luar, handler gak tahu datanya di mana
type StudentHandler struct {
	repo repository.StudentRepository
}

func NewStudentHandler(repo repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

// terjemahkanError = ubah error repository jadi status HTTP yang tepat
func terjemahkanError(c *fiber.Ctx, err error, pesanUmum string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return fail(c, fiber.StatusConflict, "NIM sudah digunakan oleh student lain")
	default:
		return fail(c, fiber.StatusInternalServerError, pesanUmum)
	}
}

// --- GET /api/v1/students ---
func (h *StudentHandler) List(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	q := parseListQuery(c)

	students, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "gagal mengambil data student")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}
	if totalPages == 0 {
		totalPages = 1
	}

	return okList(c, "daftar student berhasil diambil", students, &model.Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

// --- GET /api/v1/students/:id ---
func (h *StudentHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	student, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "gagal mengambil data student")
	}

	return ok(c, "student ditemukan", student)
}

// --- POST /api/v1/students ---
func (h *StudentHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	errs := map[string]string{}
	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "harus antara 0 sampai 100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	// Keunikan NIM dijaga database (UNIQUE INDEX), bukan loop manual
	baru, err := h.repo.Create(ctx, model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: true,
	})
	if err != nil {
		return terjemahkanError(c, err, "gagal menyimpan student")
	}

	return created(c, "student berhasil dibuat", baru,
		"/api/v1/students/"+strconv.Itoa(baru.ID))
}

// --- PUT /api/v1/students/:id ---
func (h *StudentHandler) Replace(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	errs := map[string]string{}
	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "harus antara 0 sampai 100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	hasil, err := h.repo.Update(ctx, model.Student{
		ID: id, NIM: req.NIM, Name: req.Name, Grade: req.Grade, IsActive: req.IsActive,
	})
	if err != nil {
		return terjemahkanError(c, err, "gagal memperbarui student")
	}

	return ok(c, "student berhasil diganti seluruhnya", hasil)
}

// --- PATCH /api/v1/students/:id ---
func (h *StudentHandler) Patch(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	// PATCH = baca data dulu, ubah seperlunya, lalu simpan balik
	saatIni, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "gagal mengambil data student")
	}

	errs := map[string]string{}
	if req.NIM != nil {
		v := strings.TrimSpace(*req.NIM)
		if v == "" {
			errs["nim"] = "tidak boleh kosong"
		} else {
			saatIni.NIM = v
		}
	}
	if req.Name != nil {
		v := strings.TrimSpace(*req.Name)
		if v == "" {
			errs["name"] = "tidak boleh kosong"
		} else {
			saatIni.Name = v
		}
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			errs["grade"] = "harus antara 0 sampai 100"
		} else {
			saatIni.Grade = *req.Grade
		}
	}
	if req.IsActive != nil {
		saatIni.IsActive = *req.IsActive
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	hasil, err := h.repo.Update(ctx, saatIni)
	if err != nil {
		return terjemahkanError(c, err, "gagal memperbarui student")
	}

	return ok(c, "student berhasil diperbarui sebagian", hasil)
}

// --- DELETE /api/v1/students/:id ---
func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		return terjemahkanError(c, err, "gagal menghapus student")
	}

	return noContent(c)
}
