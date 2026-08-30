# API Students — Pertemuan 3

REST API data mahasiswa menggunakan Go, Fiber v2, dan **PostgreSQL**. Data tersimpan permanen di database (tidak hilang saat server dimatikan).

## Cara Menyiapkan Database

1. Pastikan PostgreSQL sudah berjalan
2. Buat database:
   ```sql
   CREATE DATABASE praktikum_backend;
   ```
3. Jalankan migrasi:
   ```bash
   psql -U <username> -d praktikum_backend -f migrations/001_create_students.sql
   ```

## Skema Tabel

```sql
CREATE TABLE students (
    id         SERIAL       PRIMARY KEY,
    nim        VARCHAR(20)  NOT NULL,
    name       VARCHAR(100) NOT NULL,
    grade      NUMERIC(5,2) NOT NULL DEFAULT 0,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- NIM unik (case-insensitive)
CREATE UNIQUE INDEX students_nim_lower_key ON students (LOWER(nim));
-- Index pencarian nama
CREATE INDEX students_name_lower_idx ON students (LOWER(name));
```

## Environment Variables

Salin `.env.example` menjadi `.env`, lalu isi nilainya:

| Variabel | Keterangan | Contoh |
|----------|-----------|--------|
| `APP_PORT` | Port server | `3000` |
| `DB_HOST` | Host database | `localhost` |
| `DB_PORT` | Port database | `5432` |
| `DB_USER` | User PostgreSQL | `postgres` |
| `DB_PASSWORD` | Password | *(isi sendiri)* |
| `DB_NAME` | Nama database | `praktikum_backend` |
| `DB_SSLMODE` | Mode SSL | `disable` |
| `DB_MAX_CONNS` | Maks koneksi pool | `10` |

## Cara Menjalankan

```bash
cd api-students
go run .
# Server jalan di http://localhost:3000
```

## Kontrak API

| Metode | Endpoint | Deskripsi | Status Mungkin |
|--------|----------|-----------|----------------|
| GET | `/api/v1/health` | Cek server + database | 200, 503 |
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
| `search` | — | Cari berdasarkan nama (ILIKE) |
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

### Response Gagal Validasi (422)
```json
{
  "success": false,
  "message": "validasi gagal",
  "errors": { "nim": "wajib diisi", "grade": "harus antara 0 sampai 100" }
}
```

## Struktur Folder

```
api-students/
├── app/
│   ├── model/
│   │   └── student.go            — struct dan tipe data
│   └── repository/
│       └── student_repository.go — interface + implementasi SQL
├── config/
│   └── env.go                    — baca file .env
├── database/
│   └── postgres.go               — connection pool
├── migrations/
│   └── 001_create_students.sql   — skema tabel
├── .env.example                  — contoh konfigurasi
├── main.go                       — perakitan app
├── handler.go                    — handler CRUD
└── helper.go                     — response helpers
```
