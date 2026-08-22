# API Students — Pertemuan 2

REST API data mahasiswa menggunakan Go dan Fiber v2. Data disimpan di memori (belum database).

## Cara Menjalankan

```bash
cd api-students
go run .
# Server jalan di http://localhost:3000
```

## Kontrak API

| Metode | Endpoint | Deskripsi | Status Mungkin |
|--------|----------|-----------|----------------|
| GET | `/api/v1/students` | Daftar mahasiswa | 200 |
| GET | `/api/v1/students/:id` | Detail satu mahasiswa | 200, 400, 404 |
| POST | `/api/v1/students` | Tambah mahasiswa baru | 201, 400, 409, 415, 422 |
| PUT | `/api/v1/students/:id` | Ganti seluruh data | 200, 400, 404, 409, 415, 422 |
| PATCH | `/api/v1/students/:id` | Ubah sebagian data | 200, 400, 404, 409, 415, 422 |
| DELETE | `/api/v1/students/:id` | Hapus mahasiswa | 204, 400, 404 |

## Query String (GET /api/v1/students)

| Parameter | Default | Keterangan |
|-----------|---------|------------|
| `page` | 1 | Halaman ke berapa |
| `limit` | 10 | Jumlah per halaman (maks 100) |
| `search` | — | Cari berdasarkan nama (case-insensitive) |
| `sort` | id | Urut berdasarkan: id, name, nim, grade |
| `order` | asc | Arah urut: asc / desc |
| `is_active` | — | Filter: true / false |

## Contoh Request & Response

### POST /api/v1/students
```json
// Request body
{ "nim": "2024001", "name": "Budi Santoso", "grade": 85 }

// Response 201
{
  "success": true,
  "message": "student berhasil dibuat",
  "data": { "id": 1, "nim": "2024001", "name": "Budi Santoso", "grade": 85, "is_active": true, "created_at": "..." }
}
```

### PUT /api/v1/students/1 — semua field wajib
```json
// Request body
{ "nim": "2024001", "name": "Budi Baru", "grade": 90, "is_active": false }
```

### PATCH /api/v1/students/1 — hanya field yang diubah
```json
// Request body — hanya grade yang berubah, field lain tetap
{ "grade": 95 }
```

### GET /api/v1/students?page=1&limit=2&sort=name
```json
{
  "success": true,
  "message": "daftar student berhasil diambil",
  "data": [...],
  "meta": { "page": 1, "limit": 2, "total": 5, "total_pages": 3 }
}
```

### Response Gagal Validasi (422)
```json
{
  "success": false,
  "message": "validasi gagal",
  "errors": { "nim": "wajib diisi", "grade": "harus antara 0 sampai 100" }
}
```

## Status HTTP yang Diterapkan

| Status | Nama | Situasi |
|--------|------|---------|
| 200 | OK | GET / PUT / PATCH berhasil |
| 201 | Created | POST berhasil, ada header `Location` |
| 204 | No Content | DELETE berhasil, tanpa body |
| 400 | Bad Request | Body bukan JSON / id bukan angka / tidak ada field yang diubah |
| 404 | Not Found | Student tidak ditemukan |
| 409 | Conflict | NIM sudah digunakan student lain |
| 415 | Unsupported Media Type | Content-Type bukan application/json |
| 422 | Unprocessable Entity | Validasi gagal (rincian per field) |
