package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv memuat variabel dari file .env
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("peringatan: file .env tidak ditemukan, memakai environment sistem")
	}
}

// GetEnv ambil nilai environment, kalau kosong pakai fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

// GetEnvInt sama kayak GetEnv tapi hasilnya angka
func GetEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("peringatan: %s bukan angka (%q), pakai bawaan %d", key, value, fallback)
		return fallback
	}
	return parsed
}
