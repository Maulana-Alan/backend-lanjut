package service

import (
	"testing"

	"api-students/app/model"
)

// Test ini jalan TANPA server, TANPA database, TANPA fiber.Ctx
// Bukti bahwa pemisahan business rules berguna

func TestCountTotalPages(t *testing.T) {
	cases := []struct{ total, limit, want int }{
		{0, 10, 1},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{137, 20, 7},
	}

	for _, tc := range cases {
		if got := CountTotalPages(tc.total, tc.limit); got != tc.want {
			t.Errorf("total=%d limit=%d: harap %d, dapat %d",
				tc.total, tc.limit, tc.want, got)
		}
	}
}

func TestValidateCreate_Kosong(t *testing.T) {
	errs := ValidateCreate(model.CreateStudentRequest{})
	if len(errs) == 0 {
		t.Error("seharusnya ada error validasi untuk request kosong")
	}
	if _, ok := errs["nim"]; !ok {
		t.Error("nim seharusnya error")
	}
	if _, ok := errs["name"]; !ok {
		t.Error("name seharusnya error")
	}
}

func TestValidateCreate_Valid(t *testing.T) {
	errs := ValidateCreate(model.CreateStudentRequest{
		NIM: "2024001", Name: "Budi", Grade: 85,
	})
	if len(errs) != 0 {
		t.Errorf("seharusnya tidak ada error: %v", errs)
	}
}

func TestApplyPatch_SebagianField(t *testing.T) {
	awal := model.Student{ID: 1, NIM: "2024001", Name: "Budi", Grade: 85, IsActive: true}
	namaBaru := "Budi Baru"

	hasil, errs := ApplyPatch(awal, model.PatchStudentRequest{Name: &namaBaru})

	if len(errs) != 0 {
		t.Fatalf("gak seharusnya ada error: %v", errs)
	}
	if hasil.Name != "Budi Baru" {
		t.Error("name seharusnya berubah jadi Budi Baru")
	}
	if hasil.NIM != "2024001" {
		t.Error("nim seharusnya gak berubah")
	}
	if hasil.Grade != 85 {
		t.Error("grade seharusnya gak berubah")
	}
}

func TestApplyPatch_FieldKosong(t *testing.T) {
	awal := model.Student{ID: 1, NIM: "2024001", Name: "Budi"}
	nimKosong := ""

	_, errs := ApplyPatch(awal, model.PatchStudentRequest{NIM: &nimKosong})
	if _, ok := errs["nim"]; !ok {
		t.Error("nim kosong seharusnya error")
	}
}

func TestIsEmptyPatch(t *testing.T) {
	if !IsEmptyPatch(model.PatchStudentRequest{}) {
		t.Error("seharusnya empty patch")
	}
	nama := "test"
	if IsEmptyPatch(model.PatchStudentRequest{Name: &nama}) {
		t.Error("seharusnya bukan empty patch")
	}
}
