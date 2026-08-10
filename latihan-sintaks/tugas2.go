package main

import "fmt"

func main() {

	// TUGAS 2: Variabel dan Struktur Data ---

	// --- Bagian 1: Deklarasi 5 variabel tipe berbeda ---
	var nama string = "NALA SUTRISNO"               // string  = teks
	var umur int = 21                               // int     = bilangan bulat
	var ipk float64 = 3.75                          // float64 = bilangan desimal
	var aktif bool = true                           // bool    = true / false
	hobi := []string{"coding", "gaming", "membaca"} // slice = daftar dinamis

	fmt.Println("=== 5 Variabel Tipe Berbeda ===")
	fmt.Println("Nama  (string) :", nama)
	fmt.Println("Umur  (int)    :", umur)
	fmt.Println("IPK   (float64):", ipk)
	fmt.Println("Aktif (bool)   :", aktif)
	fmt.Println("Hobi  (slice)  :", hobi)
	fmt.Println() // spasi  -

	// --- Bagian 2: Map data mahasiswa (nama = kunci, nilai = isi) ---
	nilaiMahasiswa := make(map[string]int)

	// Operasi 1: MENAMBAH data ke map
	nilaiMahasiswa["Budi"] = 85
	nilaiMahasiswa["Sari"] = 92
	nilaiMahasiswa["Andi"] = 78
	nilaiMahasiswa["Dewi"] = 88

	fmt.Println("=== Map Data Mahasiswa ===")
	fmt.Println("Setelah ditambah:", nilaiMahasiswa)
	fmt.Println()

	// Operasi 2: MEMBACA dengan pengecekan keberadaan
	fmt.Println("=== Membaca dengan Pengecekan ===")

	// Cek "Sari" — ada di map
	nilai, ada := nilaiMahasiswa["Sari"]
	if ada {
		fmt.Println("Sari ditemukan, nilainya:", nilai)
	} else {
		fmt.Println("Sari tidak ditemukan")
	}

	// Cek "Rudi" — tidak ada di map
	nilai2, ada2 := nilaiMahasiswa["Rudi"]
	if ada2 {
		fmt.Println("Rudi ditemukan, nilainya:", nilai2)
	} else {
		fmt.Println("Rudi tidak ditemukan di data")
	}
	fmt.Println()

	// Operasi 3: MENGHAPUS data dari map
	fmt.Println("=== Menghapus Data ===")
	fmt.Println("Sebelum hapus:", nilaiMahasiswa)
	delete(nilaiMahasiswa, "Andi")
	fmt.Println("Setelah hapus Andi:", nilaiMahasiswa)
	fmt.Println()

	// Operasi 4: MENELUSURI (loop) seluruh isi map
	fmt.Println("=== Loop Seluruh Isi Map ===")
	for nama, nilai := range nilaiMahasiswa {
		fmt.Printf("Mahasiswa: %s, Nilai: %d\n", nama, nilai)
	}
}
