package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Instance database global (optional, tapi ConnectDB lebih clean buat di-passing)
var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var dsn string
	
	// 1. Prioritaskan DATABASE_URL atau DATABASE_PRIVATE_URL (Format universal)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_PRIVATE_URL") // Sering dipakai di Railway internal
	}

	if databaseURL != "" {
		fmt.Println("📍 Menggunakan DATABASE_URL/PRIVATE_URL untuk koneksi")
		dsn = databaseURL
	} else {
		// 2. Fallback: Cek variabel DB_* (Custom kita) atau PG* (Bawaan Railway/Postgres)
		host := getEnvFallback("DB_HOST", "PGHOST", "POSTGRES_HOST")
		user := getEnvFallback("DB_USER", "PGUSER", "POSTGRES_USER")
		password := getEnvFallback("DB_PASSWORD", "PGPASSWORD", "POSTGRES_PASSWORD")
		dbName := getEnvFallback("DB_NAME", "PGDATABASE", "DB_DATABASE", "POSTGRES_DB")
		port := getEnvFallback("DB_PORT", "PGPORT", "POSTGRES_PORT")
		sslMode := getEnvFallback("DB_SSLMODE", "PGSSLMODE", "POSTGRES_SSLMODE")

		// Jika di Railway (biasanya ada variable RAILWAY_ENVIRONMENT), jangan default ke localhost
		isRailway := os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("RAILWAY_STATIC_URL") != ""
		
		if host == "" { 
			if isRailway {
				fmt.Println("⚠️  CRITICAL: Host DB tidak ditemukan! Kamu BELUM menghubungkan service Postgres ke service ini di Railway.")
				fmt.Println("💡 SOLUSI: Di Dashboard Railway, masuk ke menu 'Variables', klik 'New Variable', lalu tambahkan 'DATABASE_URL' dengan value '${{Postgres.DATABASE_URL}}'")
				// Jangan default ke localhost kalo di Railway biar gak bingung
				host = "REQUIRED_VARIABLE_MISSING"
			} else {
				host = "localhost" 
			}
		}
		if port == "" { port = "5432" }
		if user == "" { user = "postgres" }
		
		// Fix BUG-003: Force require SSL if in Railway and no mode set
		if sslMode == "" { 
			if isRailway {
				sslMode = "require" 
			} else {
				sslMode = "disable" 
			}
		}

		// Debug: Log semua key yang tersedia (untuk bantu user cek penamaan)
		fmt.Print("📋 Variabel Environment Tersedia: ")
		for _, env := range os.Environ() {
			key := strings.Split(env, "=")[0]
			if strings.Contains(key, "DB") || strings.Contains(key, "PG") || strings.Contains(key, "PORT") || strings.Contains(key, "DATABASE") {
				fmt.Printf("%s, ", key)
			}
		}
		fmt.Println("")

		// Mask password for logging
		maskedPassword := "********"
		if password == "" {
			maskedPassword = "[kosong]"
		}

		fmt.Printf("🔍 Detail Koneksi: host=%s, user=%s, password=%s, db=%s, port=%s, sslmode=%s\n", 
			host, user, maskedPassword, dbName, port, sslMode)
		
		if dbName == "" {
			fmt.Println("❌ Error: Nama Database (DB_NAME) tidak ditemukan!")
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			host, user, password, dbName, port, sslMode)
	}

	// 3. Buka koneksi
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Mask password in error message if it's part of the DSN
		maskedDSN := maskPasswordInDSN(dsn)
		fmt.Printf("❌ Gagal koneksi database. Pastikan variabel environment sudah di-set di Railway. DSN: %s\n", maskedDSN)
		log.Fatalf("Gagal koneksi ke database PostgreSQL: %v", err)
	}

	fmt.Println("✅ Database Connected Successfully!")
	DB = db

	// Enable uuid-ossp extension (required for UUID types in some PG versions/configurations)
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		fmt.Printf("⚠️  Warning: Could not ensure uuid-ossp extension: %v (Migration might fail if DB-side UUIDs are used)\n", err)
	}

	return db
}

func getEnvFallback(keys ...string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return ""
}

// maskPasswordInDSN replaces the password in a DSN string with asterisks for logging purposes.
func maskPasswordInDSN(dsn string) string {
	parts := strings.Fields(dsn)
	for i, part := range parts {
		if strings.HasPrefix(part, "password=") {
			parts[i] = "password=********"
			break
		}
	}
	return strings.Join(parts, " ")
}
