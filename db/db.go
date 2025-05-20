package db

import (
	"log"
	"os"

	"token-transfer-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DBconnect() {
	dsn := "host=postgres user=btp password=btp dbname=btp_db port=5432 sslmode=disable TimeZone=Europe/Warsaw"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n[GORM] ", log.LstdFlags),
			logger.Config{
				IgnoreRecordNotFoundError: true,
			},
		),
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	DB = gormDB
	DB.AutoMigrate(&models.Wallet{})

	var count int64
	if err := DB.Model(&models.Wallet{}).Where("address = ?", "0x0000000000000000000000000000000000000000").Count(&count).Error; err != nil {
		log.Fatalf("%v", err)
	}
	if count == 0 {
		DB.Create(&models.Wallet{Address: "0x0000000000000000000000000000000000000000", Balance: 1000000})
	}
}
