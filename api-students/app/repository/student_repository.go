package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-students/app/model"
)

// Error milik repository — handler cuma kenal dua ini, gak perlu tahu soal pgx
var (
	ErrNotFound  = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("data sudah ada")
)

// StudentRepository = kontrak akses data, gak ada kata SQL di sini
type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, s model.Student) (model.Student, error)
	Update(ctx context.Context, s model.Student) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

// Daftar putih kolom yang boleh di-sort — cegah SQL injection
var kolomUrut = map[string]string{
	"id":    "id",
	"name":  "name",
	"nim":   "nim",
	"grade": "grade",
}

// struct implementasi — SQL hidup di sini
type studentPostgresRepo struct {
	pool *pgxpool.Pool
}

// NewStudentRepository balikin interface, bukan struct
func NewStudentRepository(pool *pgxpool.Pool) StudentRepository {
	return &studentPostgresRepo{pool: pool}
}

// buildFilter susun WHERE dari query string, nilai selalu jadi parameter ($1, $2)
func buildFilter(q model.ListQuery) (string, []any) {
	where := " WHERE 1=1"
	args := []any{}

	if q.Search != "" {
		where += fmt.Sprintf(" AND name ILIKE $%d", len(args)+1)
		args = append(args, "%"+q.Search+"%")
	}
	if q.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", len(args)+1)
		args = append(args, *q.IsActive)
	}

	return where, args
}

func (r *studentPostgresRepo) FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error) {
	where, args := buildFilter(q)

	// Hitung total dulu (sebelum LIMIT), buat meta paginasi
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM students"+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("menghitung student: %w", err)
	}

	// Ambil 1 halaman saja — filter, sort, paginasi semua dikerjakan SQL
	arah := "ASC"
	if q.Order == "desc" {
		arah = "DESC"
	}

	sql := fmt.Sprintf(
		`SELECT id, nim, name, grade, is_active, created_at
		 FROM students%s
		 ORDER BY %s %s
		 LIMIT $%d OFFSET $%d`,
		where, kolomUrut[q.Sort], arah, len(args)+1, len(args)+2,
	)
	args = append(args, q.Limit, q.Offset())

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("mengambil daftar student: %w", err)
	}
	defer rows.Close() // WAJIB — kalau lupa, koneksi pool habis

	hasil := []model.Student{}
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("membaca baris student: %w", err)
		}
		hasil = append(hasil, s)
	}

	return hasil, total, nil
}

func (r *studentPostgresRepo) FindByID(ctx context.Context, id int) (model.Student, error) {
	var s model.Student
	err := r.pool.QueryRow(ctx,
		`SELECT id, nim, name, grade, is_active, created_at
		 FROM students WHERE id = $1`, id,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("mengambil student: %w", err)
	}
	return s, nil
}

func (r *studentPostgresRepo) Create(ctx context.Context, s model.Student) (model.Student, error) {
	// RETURNING = ambil id dan created_at yang dibikin database, tanpa query kedua
	err := r.pool.QueryRow(ctx,
		`INSERT INTO students (nim, name, grade, is_active)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		s.NIM, s.Name, s.Grade, s.IsActive,
	).Scan(&s.ID, &s.CreatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("menyimpan student: %w", err)
	}
	return s, nil
}

func (r *studentPostgresRepo) Update(ctx context.Context, s model.Student) (model.Student, error) {
	err := r.pool.QueryRow(ctx,
		`UPDATE students SET nim=$1, name=$2, grade=$3, is_active=$4
		 WHERE id=$5
		 RETURNING id, nim, name, grade, is_active, created_at`,
		s.NIM, s.Name, s.Grade, s.IsActive, s.ID,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("memperbarui student: %w", err)
	}
	return s, nil
}

func (r *studentPostgresRepo) Delete(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM students WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("menghapus student: %w", err)
	}
	// Perintah jalan tapi gak ada baris yang kena = id gak ada
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// isUniqueViolation cek apakah error dari pelanggaran UNIQUE (kode 23505)
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
