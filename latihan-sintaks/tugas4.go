package main

import "fmt"

//STRUCT

// Struct Student — menyimpan data mahasiswa
type Student struct {
	ID       int     // nomor identitas
	Name     string  // nama mahasiswa
	Grade    float64 // nilai (0.0 - 100.0)
	IsActive bool    // status aktif atau tidak
}

// GetInfo — mengembalikan info lengkap student
// Pakai VALUE receiver (s Student) karena hanya MEMBACA data, tidak mengubah
func (s Student) GetInfo() string {
	status := "Tidak Aktif"
	if s.IsActive {
		status = "Aktif"
	}
	return fmt.Sprintf("ID: %d | Nama: %s | Nilai: %.1f | Status: %s",
		s.ID, s.Name, s.Grade, status)
}

// UpdateGrade — memperbarui nilai student
// Pakai POINTER receiver (s *Student) karena MENGUBAH data asli
func (s *Student) UpdateGrade(grade float64) {
	s.Grade = grade
}

// Activate — mengubah status menjadi aktif
// Pakai POINTER receiver karena MENGUBAH data asli
func (s *Student) Activate() {
	s.IsActive = true
}

// Deactivate — mengubah status menjadi tidak aktif
// Pakai POINTER receiver karena MENGUBAH data asli
func (s *Student) Deactivate() {
	s.IsActive = false
}

func main() {
	// Buat student baru
	mhs := Student{
		ID:       1,
		Name:     "Budi Santoso",
		Grade:    75.5,
		IsActive: false,
	}

	// Tampilkan info awal
	fmt.Println("=== Info Awal ===")
	fmt.Println(mhs.GetInfo())
	fmt.Println()

	// Aktifkan student
	fmt.Println("=== Setelah Activate() ===")
	mhs.Activate()
	fmt.Println(mhs.GetInfo())
	fmt.Println()

	// Update nilai
	fmt.Println("=== Setelah UpdateGrade(90.5) ===")
	mhs.UpdateGrade(90.5)
	fmt.Println(mhs.GetInfo())
	fmt.Println()

	// Deaktifkan student
	fmt.Println("=== Setelah Deactivate() ===")
	mhs.Deactivate()
	fmt.Println(mhs.GetInfo())
	fmt.Println()

	// --- Penjelasan pilihan receiver ---
	fmt.Println("=== Penjelasan Receiver ===")
	fmt.Println("GetInfo()     → VALUE receiver   → hanya baca, tidak ubah data")
	fmt.Println("UpdateGrade() → POINTER receiver → mengubah Grade, perlu akses data asli")
	fmt.Println("Activate()    → POINTER receiver → mengubah IsActive, perlu akses data asli")
	fmt.Println("Deactivate()  → POINTER receiver → mengubah IsActive, perlu akses data asli")
}
