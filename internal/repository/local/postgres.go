package local

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dB *gorm.DB

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GORM: %v\n", err)
	}

	sqlDB.SetMaxOpenConns(5)          
	sqlDB.SetMaxIdleConns(3)           
	sqlDB.SetConnMaxLifetime(30 * time.Minute) 

	// Run migrations
	if err := db.AutoMigrate(&domain.ConnectionDetails{}, &domain.Table{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v\n", err)
	}
	if err := db.AutoMigrate( &domain.Schema{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v\n", err)
	}
	if err := db.AutoMigrate(&domain.Summary{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v\n", err)
	}

	dB = db
	log.Println("Connected to PostgreSQL using GORM with connection pooling")
}
