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
	
	// 1. Prioritaskan DATABASE_URL (Format universal)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		fmt.Println("📍 Menggunakan DATABASE_URL untuk koneksi")
		dsn = databaseURL
	} else {
		// 2. Fallback: Cek variabel DB_* (Custom kita) atau PG* (Bawaan Railway/Postgres)
		host := getEnvFallback("DB_HOST", "PGHOST")
		user := getEnvFallback("DB_USER", "PGUSER")
		password := getEnvFallback("DB_PASSWORD", "PGPASSWORD")
		dbName := getEnvFallback("DB_NAME", "PGDATABASE", "DB_DATABASE")
		port := getEnvFallback("DB_PORT", "PGPORT")
		sslMode := getEnvFallback("DB_SSLMODE", "PGSSLMODE")

		if host == "" { host = "localhost" }
		if port == "" { port = "5432" }
		if user == "" { user = "postgres" }
		if sslMode == "" { sslMode = "disable" }

		// Mask password for logging
		maskedPassword := "********"
		if password == "" {
			maskedPassword = "[empty]"
		}

		fmt.Printf("🔍 Info Koneksi: host=%s, user=%s, password=%s, db=%s, port=%s, sslmode=%s\n", 
			host, user, maskedPassword, dbName, port, sslMode)
		
		if dbName == "" {
			fmt.Println("⚠️  Peringatan: Nama Database (DB_NAME) masih kosong!")
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
