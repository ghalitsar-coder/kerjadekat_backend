package main

import (
	"log"
	"kerjadekat/backend/config"
	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/pkg/database"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Config load failed: %v", err)
	}

	db, err := database.OpenPostgres(cfg)
	if err != nil {
		log.Fatalf("DB Open Error: %v", err)
	}

	log.Println("Running AutoMigrate for User schema...")
	// Hapus kode AutoMigrate yang lama dan ganti dengan ini:
    if err := db.AutoMigrate(
        &domain.Kelurahan{},
        &domain.SkillCategory{},
        &domain.User{},
        &domain.WorkerProfile{},
        &domain.WorkerSkill{},
        &domain.AgentTerritory{},
        &domain.Order{},
        &domain.OrderStatusLog{},
        &domain.OrderRating{},
        &domain.IncomeRecord{},
        &domain.Wallet{},
        &domain.WalletTransaction{},
    ); err != nil {
        log.Printf("Gagal melakukan migrasi database: %v", err)
    }
	
	log.Println("Database migration completed successfully!")
}
