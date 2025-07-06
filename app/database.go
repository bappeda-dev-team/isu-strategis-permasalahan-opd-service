package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func GetConnection() *sql.DB {
	godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}

	log.Printf("Mencoba koneksi ke database %s di %s:%s...", dbName, dbHost, dbPort)

	// Format koneksi MySQL: user:password@tcp(host:port)/dbname?parseTime=true
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Error membuka koneksi database: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Gagal terhubung ke database dalam 10 detik: %v", err)
		log.Printf("Mencoba koneksi ulang...")

		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err != nil {
			db.Close()
			log.Fatalf("Koneksi database gagal setelah percobaan ulang: %v", err)
		}
	}

	log.Printf("Berhasil terhubung ke database %s", dbName)
	log.Printf("Max Open Connections: %d", db.Stats().MaxOpenConnections)
	log.Printf("Open Connections: %d", db.Stats().OpenConnections)
	log.Printf("In Use Connections: %d", db.Stats().InUse)
	log.Printf("Idle Connections: %d", db.Stats().Idle)

	return db
}

func RunFlyway() {
	godotenv.Load()

	// Mencari executable Flyway berdasarkan sistem operasi
	var flywayCmd string
	if runtime.GOOS == "windows" {
		// Di Windows, coba cari flyway.cmd di PATH
		flywayPath := os.Getenv("FLYWAY_HOME")
		if flywayPath == "" {
			log.Fatal("Environment variable FLYWAY_HOME harus diatur di Windows (contoh: C:\\Program Files\\flyway)")
		}
		flywayCmd = filepath.Join(flywayPath, "flyway.cmd")
	} else {
		flywayCmd = "flyway" // Untuk Linux/MacOS
	}

	cmd := exec.Command(flywayCmd, "-locations=filesystem:./db/migrations", "migrate")

	// Menggunakan JAVA_HOME dari environment variable
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		log.Fatal("Environment variable JAVA_HOME harus diatur")
	}

	// Membuat path yang sesuai dengan sistem operasi
	javaPath := filepath.Join(javaHome, "bin")
	newPath := fmt.Sprintf("%s%s%s", javaPath, string(os.PathListSeparator), os.Getenv("PATH"))

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbHost == "" || dbPort == "" || dbName == "" || dbUser == "" {
		log.Fatal("Environment variables DB_HOST, DB_PORT, DB_NAME, DB_USER must be set")
	}

	// Format JDBC untuk MySQL
	jdbcURL := fmt.Sprintf("jdbc:mysql://%s:%s/%s?useSSL=false&serverTimezone=UTC", dbHost, dbPort, dbName)

	cmd.Env = append(os.Environ(),
		"PATH="+newPath,
		"FLYWAY_URL="+jdbcURL,
		"FLYWAY_USER="+dbUser,
		"FLYWAY_PASSWORD="+dbPassword,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Migrasi gagal dijalankan: %v", err)
	} else {
		log.Println("Migrasi berhasil dijalankan")
	}
}
