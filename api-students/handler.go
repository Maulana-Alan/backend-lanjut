package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Penyimpanan data di memori — hilang saat server dimatikan (database masuk pertemuan 3)
var students []Student
var nextID = 1

// findIndex mencari posisi student berdasarkan ID, -1 kalau tidak ketemu
func findIndex(id int) int {
	for i := range students {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

// nimExists cek apakah NIM sudah dipakai student lain (excludeID = abaikan ID tertentu, untuk PUT/PATCH)
func nimExists(nim string, excludeID int) bool {
	for _, s := range students {
		if strings.EqualFold(s.NIM, nim) && s.ID != excludeID {
			return true
		}
	}
	return false
}

// paramID ambil :id dari URL, validasi harus angka positif
func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

// --- Handler: GET /api/v1/students ---
func listStudents(c *fiber.Ctx) error {
	q := parseListQuery(c)

	// 1. Saring berdasarkan is_active dan search
	hasil := []Student{}
	for _, s := range students {
		if q.IsActive != nil && s.IsActive != *q.IsActive {
			continue
		}
		if q.Search != "" && !strings.Contains(strings.ToLower(s.Name), strings.ToLower(q.Search)) {
			continue
		}
		hasil = append(hasil, s)
	}

	// 2. Urutkan
	sort.SliceStable(hasil, func(i, j int) bool {
		var lebihKecil bool
		switch q.Sort {
		case "name":
			lebihKecil = hasil[i].Name < hasil[j].Name
		case "nim":
			lebihKecil = hasil[i].NIM < hasil[j].NIM
		case "grade":
			lebihKecil = hasil[i].Grade < hasil[j].Grade
		default:
			lebihKecil = hasil[i].ID < hasil[j].ID
		}
		if q.Order == "desc" {
			return !lebihKecil
		}
		return lebihKecil
	})

	// 3. Potong sesuai halaman
	total := len(hasil)
	totalPages := (total + q.Limit - 1) / q.Limit
	if totalPages == 0 {
		totalPages = 1
	}
	mulai := (q.Page - 1) * q.Limit
	if mulai > total {
		mulai = total
	}
	akhir := mulai + q.Limit
	if akhir > total {
		akhir = total
	}

	return okList(c, "daftar student berhasil diambil", hasil[mulai:akhir], &Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

// --- Handler: GET /api/v1/students/:id ---
func getStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}

	return ok(c, "student ditemukan", students[i])
}

// --- Handler: POST /api/v1/students ---
func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	// Validasi per field
	errs := map[string]string{}
	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	} else if nimExists(req.NIM, 0) {
		return fail(c, fiber.StatusConflict, "NIM sudah digunakan oleh student lain")
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

	baru := Student{
		ID:        nextID,
		NIM:       req.NIM,
		Name:      req.Name,
		Grade:     req.Grade,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	students = append(students, baru)
	nextID++

	return created(c, "student berhasil dibuat", baru, "/api/v1/students/"+strconv.Itoa(baru.ID))
}

// --- Handler: PUT /api/v1/students/:id ---
// PUT = ganti SELURUH isi, semua field wajib dikirim
func replaceStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}

	var req ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	errs := map[string]string{}
	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
	} else if nimExists(req.NIM, id) {
		return fail(c, fiber.StatusConflict, "NIM sudah digunakan oleh student lain")
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

	// Ganti semua field sekaligus
	students[i].NIM = req.NIM
	students[i].Name = req.Name
	students[i].Grade = req.Grade
	students[i].IsActive = req.IsActive

	return ok(c, "student berhasil diganti seluruhnya", students[i])
}

// --- Handler: PATCH /api/v1/students/:id ---
// PATCH = ubah SEBAGIAN, field yang tidak dikirim tidak disentuh
func patchStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}

	var req PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	// Kalau semua nil, client tidak mengirim apa-apa
	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	// Hanya update field yang non-nil
	errs := map[string]string{}
	if req.NIM != nil {
		v := strings.TrimSpace(*req.NIM)
		if v == "" {
			errs["nim"] = "tidak boleh kosong"
		} else if nimExists(v, id) {
			return fail(c, fiber.StatusConflict, "NIM sudah digunakan oleh student lain")
		} else {
			students[i].NIM = v
		}
	}
	if req.Name != nil {
		v := strings.TrimSpace(*req.Name)
		if v == "" {
			errs["name"] = "tidak boleh kosong"
		} else {
			students[i].Name = v
		}
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			errs["grade"] = "harus antara 0 sampai 100"
		} else {
			students[i].Grade = *req.Grade
		}
	}
	if req.IsActive != nil {
		students[i].IsActive = *req.IsActive
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	return ok(c, "student berhasil diperbarui sebagian", students[i])
}

// --- Handler: DELETE /api/v1/students/:id ---
func deleteStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}

	// Hapus dari slice: sambungkan bagian sebelum dan sesudah index
	students = append(students[:i], students[i+1:]...)

	return noContent(c)
}
