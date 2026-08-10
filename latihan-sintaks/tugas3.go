package main

import "fmt"

// POINTER
// Function 1: swap — menukar nilai dua integer lewat pointer
// Pakai pointer (*int) supaya perubahan TERSIMPAN di variabel asli
func swap(a, b *int) {
	*a, *b = *b, *a // tukar isi di alamat a dan b
}

// Function 2: updateSlice — menambah item baru ke slice lewat pointer
// Pakai pointer (*[]string) karena append() bisa bikin slice baru di memori,
// jadi tanpa pointer, slice asli TIDAK berubah
func updateSlice(s *[]string, newItem string) {
	*s = append(*s, newItem)
}

// --- Perbandingan: pass by value (SALINAN, asli tidak berubah) ---
func ubahNilaiValue(x int) {
	x = 999 // yang berubah cuma salinan, aslinya tetap
}

// --- Perbandingan: pass by pointer (ALAMAT, asli ikut berubah) ---
func ubahNilaiPointer(x *int) {
	*x = 999 // mengubah isi di alamat asli
}

func main() {
	// --- Demo swap ---
	fmt.Println("=== Demo swap ===")
	a, b := 10, 20
	fmt.Println("Sebelum swap: a =", a, ", b =", b)
	swap(&a, &b) // kirim ALAMAT pakai &
	fmt.Println("Sesudah swap: a =", a, ", b =", b)
	fmt.Println()

	// --- Demo updateSlice ---
	fmt.Println("=== Demo updateSlice ===")
	buah := []string{"apel", "jeruk"}
	fmt.Println("Sebelum update:", buah)
	updateSlice(&buah, "mangga")
	updateSlice(&buah, "durian")
	fmt.Println("Sesudah update:", buah)
	fmt.Println()

	// --- Perbandingan pass by value vs pass by pointer ---
	fmt.Println("=== Perbandingan Value vs Pointer ===")
	angka := 42

	fmt.Println("Nilai awal:", angka)

	ubahNilaiValue(angka) // kirim SALINAN
	fmt.Println("Setelah pass by value:", angka, " ← tetap 42, karena yang diubah cuma salinan")

	ubahNilaiPointer(&angka) // kirim ALAMAT
	fmt.Println("Setelah pass by pointer:", angka, " ← jadi 999, karena yang diubah data aslinya")
}
