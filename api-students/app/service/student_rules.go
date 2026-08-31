package service

import (
	"strings"

	"api-students/app/model"
)

// File ini berisi business rules MURNI:
// gak sentuh fiber.Ctx, gak sentuh database, gak tahu soal HTTP.
// Bisa dites tanpa nyalain server.

// ValidateCreate cek isi request POST
func ValidateCreate(req model.CreateStudentRequest) map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "wajib diisi"
	}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "harus antara 0 sampai 100"
	}
	return errs
}

// ValidateReplace cek isi request PUT (semua field wajib)
func ValidateReplace(req model.ReplaceStudentRequest) map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "harus antara 0 sampai 100"
	}
	return errs
}

// ApplyPatch timpa field yang dikirim ke data yang sudah ada
// Field nil = gak disentuh
func ApplyPatch(current model.Student, req model.PatchStudentRequest) (model.Student, map[string]string) {
	errs := map[string]string{}

	if req.NIM != nil {
		v := strings.TrimSpace(*req.NIM)
		if v == "" {
			errs["nim"] = "tidak boleh kosong"
		} else {
			current.NIM = v
		}
	}
	if req.Name != nil {
		v := strings.TrimSpace(*req.Name)
		if v == "" {
			errs["name"] = "tidak boleh kosong"
		} else {
			current.Name = v
		}
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			errs["grade"] = "harus antara 0 sampai 100"
		} else {
			current.Grade = *req.Grade
		}
	}
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	return current, errs
}

// IsEmptyPatch cek apakah PATCH gak kirim field sama sekali
func IsEmptyPatch(req model.PatchStudentRequest) bool {
	return req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil
}

// CountTotalPages hitung jumlah halaman (bulatkan ke atas)
func CountTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	pages := (total + limit - 1) / limit
	if pages == 0 {
		pages = 1
	}
	return pages
}
