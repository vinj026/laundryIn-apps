package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Instance database global (optional, tapi ConnectDB lebih clean buat di-passing)
var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var dsn string
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL != "" {
		// Jika ada DATABASE_URL (stau format standar di Railway/Heroku/Vercel)
		fmt.Println("📍 Using DATABASE_URL for connection")
		dsn = databaseURL
	} else {
		// Fallback ke variabel individu (biasanya buat local)
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslMode := os.Getenv("DB_SSLMODE")

		if sslMode == "" {
			sslMode = "disable"
		}

		// Logging untuk debug (jangan log password beneran)
		fmt.Printf("🔍 Connection Info: host=%s, user=%s, db=%s, port=%s, sslmode=%s\n",
			host, user, dbName, port, sslMode)

		if host == "" || user == "" || dbName == "" || port == "" {
			fmt.Println("⚠️  Warning: Missing one or more individual DB environment variables!")
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			host, user, password, dbName, port, sslMode)
	}

	// 3. Buka koneksi pake GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("❌ Failed to connect to database using DSN: %s\n", maskDSN(dsn))
		log.Fatalf("Gagal koneksi ke database PostgreSQL: %v", err)
	}

	fmt.Println("✅ Database Connected Successfully!")

	DB = db
	return db
}

func maskDSN(dsn string) string {
	// Sederhana aja buat masking password di log error
	return "host=... user=... password=***** dbname=... port=..."
}
