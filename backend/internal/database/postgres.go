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
	// 1. Ambil data dari file .env yang udah kita set tadi
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")

	// 2. Susun Data Source Name (DSN) sesuai format PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbName, port, sslMode)

	// 3. Buka koneksi pake GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Tambahkan logger kalau mau liat query SQL yang dijalanin Go di terminal
	})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database PostgreSQL: %v", err)
	}

	fmt.Println("✅ Database Connected Successfully!")

	DB = db
	return db
}
